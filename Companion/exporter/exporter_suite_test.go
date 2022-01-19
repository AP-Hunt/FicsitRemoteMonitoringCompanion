package exporter_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

var FRMServer *FRMServerFake

func TestExporter(t *testing.T) {
	RegisterFailHandler(Fail)

	BeforeSuite(func() {
		FRMServer = NewFRMServerFake()
		FRMServer.Start()
	})

	AfterSuite(func() {
		FRMServer.Stop()
	})
	RunSpecs(t, "Exporter Suite")
}

func f(x float64) *float64 { return &x }

func gaugeValue(gauge *prometheus.GaugeVec, labelValues ...string) (float64, error) {
	var m = &dto.Metric{}

	err := gauge.WithLabelValues(labelValues...).Write(m)
	if err != nil {
		return -1, err
	}

	return m.Gauge.GetValue(), nil
}

func getMetric(gauge *prometheus.GaugeVec, labelValues ...string) (*dto.Metric, error) {
	var m = &dto.Metric{}

	metric, err := gauge.GetMetricWithLabelValues(labelValues...)
	if err != nil {
		return nil, err
	}

	if metric == nil {
		return nil, nil
	}

	metric.Write(m)

	return m, nil
}

func hasLabelPair(metric *dto.Metric, key string, value string) bool {
	for _, lbl := range metric.Label {
		if lbl != nil {
			if lbl.GetName() == key && lbl.GetValue() == value {
				return true
			}
		}
	}

	return false
}
