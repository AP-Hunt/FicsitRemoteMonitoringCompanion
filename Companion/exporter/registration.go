package exporter

import (
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
	RegisteredMetricVectors = append(RegisteredMetricVectors, MetricVectorDetails{
		Name:   opts.Name,
		Help:   opts.Help,
		Labels: labelNames,
	})

	metric := promauto.NewGaugeVec(opts, labelNames)
	RegisteredMetrics = append(RegisteredMetrics, metric)
	return metric
}
