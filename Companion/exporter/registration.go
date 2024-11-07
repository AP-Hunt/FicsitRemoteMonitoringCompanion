package exporter

import (
	"maps"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricVectorDetails struct {
	Name   string
	Help   string
	Labels []string
}

var RegisteredMetricVectors = []MetricVectorDetails{}
var RegisteredMetrics = []*prometheus.GaugeVec{}

func RegisterNewGaugeVec(opts prometheus.GaugeOpts, labelNames []string) *prometheus.GaugeVec {
	// All metrics include url and session_name labels
	labelNames = append(labelNames, "url", "session_name")
	RegisteredMetricVectors = append(RegisteredMetricVectors, MetricVectorDetails{
		Name:   opts.Name,
		Help:   opts.Help,
		Labels: labelNames,
	})

	metric := promauto.NewGaugeVec(opts, labelNames)
	RegisteredMetrics = append(RegisteredMetrics, metric)
	return metric
}

type MetricsDropper struct {
	GaugeVecs       []*prometheus.GaugeVec
	OldMetricLabels []prometheus.Labels
	NewMetricLabels []prometheus.Labels
}

func NewMetricsDropper(gauges ...*prometheus.GaugeVec) *MetricsDropper {
	return &MetricsDropper{
		GaugeVecs: gauges,
	}
}

func (m *MetricsDropper) CacheFreshMetricLabel(newLabels prometheus.Labels) {
	m.NewMetricLabels = append(m.NewMetricLabels, newLabels)
}

func (m *MetricsDropper) DropStaleMetricLabels() {
	for _, labels := range m.OldMetricLabels {
		dropLabels := true
		for _, newLabels := range m.NewMetricLabels {
			if maps.Equal(labels, newLabels) {
				dropLabels = false
			}
		}
		if dropLabels {
			for _, gaugeVec := range m.GaugeVecs {
				gaugeVec.DeletePartialMatch(labels)
			}
		}
	}
	//rotate metric labels, consider the previously evaluated labels as old
	m.OldMetricLabels = m.NewMetricLabels
	m.NewMetricLabels = []prometheus.Labels{}
}
