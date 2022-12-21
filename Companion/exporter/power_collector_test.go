package exporter_test

import (
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/m/v2/exporter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PowerCollector", func() {
	var collector *exporter.PowerCollector

	BeforeEach(func() {
		FRMServer.Reset()
		collector = exporter.NewPowerCollector("http://localhost:9080/getPower")

		FRMServer.ReturnsPowerData([]exporter.PowerDetails{
			{
				CircuitId:           1,
				PowerConsumed:       30,
				PowerCapacity:       44,
				PowerMaxConsumed:    50,
				BatteryDifferential: 12,
				BatteryPercent:      60,
				BatteryCapacity:     100,
				BatteryTimeEmpty:    "00:00:00",
				BatteryTimeFull:     "33:22:11",
				FuseTriggered:       false,
			},
			{
				CircuitId:           2,
				PowerConsumed:       55,
				PowerCapacity:       44,
				PowerMaxConsumed:    60,
				BatteryDifferential: -12,
				BatteryPercent:      60,
				BatteryCapacity:     100,
				BatteryTimeEmpty:    "00:3:00",
				BatteryTimeFull:     "00:00:00",
				FuseTriggered:       true,
			},
		})
	})

	AfterEach(func() {
		collector = nil
	})

	Describe("Power metrics collection", func() {
		It("sets the 'power_consumed' metric with the right labels", func() {
			collector.Collect()

			val, err := gaugeValue(exporter.PowerConsumed, "1")

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(30)))
		})
		It("sets the 'battery_seconds_full' metric with the right labels", func() {
			collector.Collect()

			val, err := gaugeValue(exporter.BatterySecondsFull, "1")

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(120131)))
		})
		It("sets the 'fuse_triggered' metric with the right labels", func() {
			collector.Collect()

			val, err := gaugeValue(exporter.FuseTriggered, "1")

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(0)))

			val2, err := gaugeValue(exporter.FuseTriggered, "2")
			Expect(val2).To(Equal(float64(1)))
		})
		It("sets the 'battery_differential' metric with the right labels", func() {
			collector.Collect()

			val, err := gaugeValue(exporter.BatteryDifferential, "1")

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(12)))

			val2, err := gaugeValue(exporter.BatteryDifferential, "2")
			Expect(val2).To(Equal(float64(-12)))
		})
	})
})
