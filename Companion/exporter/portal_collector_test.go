package exporter_test

import (
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/exporter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PortalCollector", func() {
	var url string
	var sessionName = "default"
	var collector *exporter.PortalCollector

	BeforeEach(func() {
		FRMServer.Reset()
		url = FRMServer.server.URL
		collector = exporter.NewPortalCollector("/getPortal")

		FRMServer.ReturnsPortalData([]exporter.PortalDetails{
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

	Describe("portal metrics collection", func() {
		It("sets the 'portal_power' metric with the right labels", func() {
			collector.Collect(url, sessionName)

			val, err := gaugeValue(exporter.PortalPower, "1", url, sessionName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(60.1))
		})
		It("sets the 'portal_power_max' metric with the right labels", func() {
			collector.Collect(url, sessionName)

			val, err := gaugeValue(exporter.PortalPowerMax, "1", url, sessionName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(90.0))
		})
	})
})
