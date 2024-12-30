package exporter_test

import (
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/exporter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResourceSinkCollector", func() {
	var url string
	var sessionName = "default"
	var collector *exporter.ResourceSinkCollector

	BeforeEach(func() {
		FRMServer.Reset()
		url = FRMServer.server.URL
		collector = exporter.NewResourceSinkCollector("/getResourceSinkBuilding", "/getResourceSink")

		FRMServer.ReturnsResourceSinkData([]exporter.ResourceSinkDetails{
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

		FRMServer.ReturnsGlobalResourceSinkData([]exporter.ResourceSinkGlobalDetails{
			{
				TotalPoints:    100,
				PointsToCoupon: 200,
				NumCoupon:      1,
			},
		})
	})

	AfterEach(func() {
		collector = nil
	})

	Describe("Resource sink metrics collection", func() {
		It("sets the 'resource_sink_power' metric with the right labels", func() {
			collector.Collect(url, sessionName)

			val, err := gaugeValue(exporter.ResourceSinkPower, "1", url, sessionName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(60.1))
		})
		It("sets the 'resource_sink_power_max' metric with the right labels", func() {
			collector.Collect(url, sessionName)

			val, err := gaugeValue(exporter.ResourceSinkPowerMax, "1", url, sessionName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(90.0))
		})
	})

	Describe("Resource sink global metrics collection", func() {
		It("sets the 'resource_sink_total_points' metric with the right labels", func() {
			collector.Collect(url, sessionName)

			val, err := gaugeValue(exporter.ResourceSinkTotalPoints, url, sessionName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(100.0))
		})
		It("sets the 'resource_sink_points_to_coupon' metric with the right labels", func() {
			collector.Collect(url, sessionName)

			val, err := gaugeValue(exporter.ResourceSinkPointsToCoupon, url, sessionName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(200.0))
		})
		It("sets the 'resource_sink_collected_coupons' metric with the right labels", func() {
			collector.Collect(url, sessionName)

			val, err := gaugeValue(exporter.ResourceSinkCollectedCoupons, url, sessionName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(1.0))
		})
	})
})
