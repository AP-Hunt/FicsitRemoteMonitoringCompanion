package exporter_test

import (
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/exporter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
)

func expectGauge(metric *prometheus.GaugeVec, labels ...string) Assertion {
	val, err := gaugeValue(metric, labels...)
	Expect(err).ToNot(HaveOccurred())
	return Expect(val)
}

func eventuallyExpectGauge(metric *prometheus.GaugeVec, labels ...string) AsyncAssertion {
	val, err := gaugeValue(metric, labels...)
	Expect(err).ToNot(HaveOccurred())
	return Eventually(val)
}

var _ = Describe("metrics dropper", func() {
	var dropper *exporter.MetricsDropper
	var metric *prometheus.GaugeVec
	var globalReg *prometheus.Registry
	BeforeEach(func() {
		globalReg = prometheus.NewRegistry()
		metric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "test_metric",
			Help: "Test Metric",
		}, []string{"label1", "label2"})
		globalReg.MustRegister(metric)
		dropper = exporter.NewMetricsDropper(metric)
	})
	AfterEach(func() {
		globalReg.Unregister(metric)
	})
	It("Keeps fresh metrics", func() {
		dropper.CacheFreshMetricLabel(prometheus.Labels{"label1":"val1"})
		metric.WithLabelValues("val1","val2").Set(1)
		dropper.DropStaleMetricLabels()
		expectGauge(metric, "val1", "val2").To(Equal(1.0))
	})
	It("drops old metrics", func() {
		dropper.CacheFreshMetricLabel(prometheus.Labels{"label1":"val1"})
		metric.WithLabelValues("val1","val2").Set(1)
		metric.WithLabelValues("val3","val4").Set(1)
		dropper.DropStaleMetricLabels()
		metric.WithLabelValues("val3","val4").Set(2)
		dropper.DropStaleMetricLabels()
		expectGauge(metric, "val1", "val2").To(Equal(0.0))
		expectGauge(metric, "val3", "val4").To(Equal(2.0))
	})
})
