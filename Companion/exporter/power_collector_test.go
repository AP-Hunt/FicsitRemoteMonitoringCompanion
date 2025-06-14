package exporter_test

import (
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/exporter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PowerCollector", func() {
	var collector *exporter.PowerCollector
	var url string
	var sessionName = "default"

	BeforeEach(func() {
		FRMServer.Reset()
		url = FRMServer.server.URL
		collector = exporter.NewPowerCollector("/getPower")

		FRMServer.ReturnsPowerData([]exporter.PowerDetails{
			{
				CircuitGroupId:      1,
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
				CircuitGroupId:      2,
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
			collector.Collect(url, sessionName)

			val, err := gaugeValue(exporter.PowerConsumed, "1", url, sessionName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(30)))
		})
		It("sets the 'battery_seconds_full' metric with the right labels", func() {
			collector.Collect(url, sessionName)

			val, err := gaugeValue(exporter.BatterySecondsFull, "1", url, sessionName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(120131)))
		})
		It("sets the 'fuse_triggered' metric with the right labels", func() {
			collector.Collect(url, sessionName)

			val, err := gaugeValue(exporter.FuseTriggered, "1", url, sessionName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(0)))

			val2, _ := gaugeValue(exporter.FuseTriggered, "2", url, sessionName)
			Expect(val2).To(Equal(float64(1)))
		})
		It("sets the 'battery_differential' metric with the right labels", func() {
			collector.Collect(url, sessionName)

			val, err := gaugeValue(exporter.BatteryDifferential, "1", url, sessionName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(12)))

			val2, _ := gaugeValue(exporter.BatteryDifferential, "2", url, sessionName)
			Expect(val2).To(Equal(float64(-12)))
		})
		It("corrects max power when calculated is higher", func() {
			exporter.FactoryPowerMax.WithLabelValues("1", url, sessionName).Set(100.0)
			exporter.ExtractorPowerMax.WithLabelValues("1", url, sessionName).Set(100.0)
			exporter.DronePortPowerMax.WithLabelValues("1", url, sessionName).Set(100.0)
			exporter.FrackingPowerMax.WithLabelValues("1", url, sessionName).Set(100.0)
			exporter.HypertubePowerMax.WithLabelValues("1", url, sessionName).Set(100.0)
			exporter.PortalPowerMax.WithLabelValues("1", url, sessionName).Set(100.0)
			exporter.PumpPowerMax.WithLabelValues("1", url, sessionName).Set(100.0)
			exporter.ResourceSinkPowerMax.WithLabelValues("1", url, sessionName).Set(100.0)
			exporter.TrainCircuitPowerMax.WithLabelValues("1", url, sessionName).Set(100.0)
			exporter.TrainStationPowerMax.WithLabelValues("1", url, sessionName).Set(100.0)
			exporter.VehicleStationPowerMax.WithLabelValues("1", url, sessionName).Set(100.0)
			collector.Collect(url, sessionName)

			val, err := gaugeValue(exporter.PowerMaxConsumed, "1", url, sessionName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(1100.0))
		})
	})
})
