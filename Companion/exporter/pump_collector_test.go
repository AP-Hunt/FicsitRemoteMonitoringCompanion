package exporter_test

import (
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/exporter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PumpCollector", func() {
	var url string
	var sessionName = "default"
	var collector *exporter.PumpCollector

	BeforeEach(func() {
		FRMServer.Reset()
		url = FRMServer.server.URL
		collector = exporter.NewPumpCollector("/getPump")

		FRMServer.ReturnsPumpData([]exporter.PumpDetails{
			{
				PowerInfo: exporter.PowerInfo{
					CircuitGroupId:   1,
					PowerConsumed:    30,
					MaxPowerConsumed: 30,
				},
			},
			{
				PowerInfo: exporter.PowerInfo{
					CircuitGroupId:   1,
					PowerConsumed:    30,
					MaxPowerConsumed: 30,
				},
			},
			{
				PowerInfo: exporter.PowerInfo{
					CircuitGroupId:   1,
					PowerConsumed:    0.1,
					MaxPowerConsumed: 30,
				},
			},
		})
	})

	AfterEach(func() {
		collector = nil
	})

	Describe("Resource sink metrics collection", func() {
		It("sets the 'pump_power' metric with the right labels", func() {
			collector.Collect(url, sessionName)

			val, err := gaugeValue(exporter.PumpPower, "1", url, sessionName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(60.1))
		})
		It("sets the 'pump_power_max' metric with the right labels", func() {
			collector.Collect(url, sessionName)

			val, err := gaugeValue(exporter.PumpPowerMax, "1", url, sessionName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(90.0))
		})
	})
})
