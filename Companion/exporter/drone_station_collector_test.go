package exporter_test

import (
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/exporter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DroneStationCollector", func() {
	var collector *exporter.DroneStationCollector
	var url = "http://localhost:9080"
	var sessionName = "default"

	BeforeEach(func() {
		FRMServer.Reset()
		collector = exporter.NewDroneStationCollector("/getDroneStation")

		FRMServer.ReturnsDroneStationData([]exporter.DroneStationDetails{
			{
				Id:               "1",
				HomeStation:      "home",
				PairedStation:    "remote station",
				DroneStatus:      "EDS_EN_ROUTE",
				AvgIncRate:       1,
				AvgOutRate:       1,
				LatestIncStack:   0.2,
				LatestOutStack:   0.3,
				LatestRndTrip:    264,
				LatestTripIncAmt: 82,
				LatestTripOutAmt: 50,
				EstBatteryRate:   30,
				PowerInfo: exporter.PowerInfo{
					CircuitGroupId: 1.0,
					PowerConsumed:  100,
				},
			},
			{
				Id:               "2",
				HomeStation:      "home2",
				PairedStation:    "remote station2",
				DroneStatus:      "EDS_EN_ROUTE",
				AvgIncRate:       1,
				AvgOutRate:       1,
				LatestIncStack:   0.2,
				LatestOutStack:   0.3,
				LatestRndTrip:    264,
				LatestTripIncAmt: 82,
				LatestTripOutAmt: 50,
				EstBatteryRate:   30,
				PowerInfo: exporter.PowerInfo{
					CircuitGroupId: 1.0,
					PowerConsumed:  100,
				},
			},
		})
	})

	AfterEach(func() {
		collector = nil
	})

	Describe("Drone metrics collection", func() {
		It("sets the 'drone_port_battery_rate' metric with the right labels", func() {
			collector.Collect(url, sessionName)

			val, err := gaugeValue(exporter.DronePortBatteryRate, "1", "home", "remote station", url, sessionName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(30)))
		})
		It("sets the 'drone_port_round_trip_seconds' metric with the right labels", func() {
			collector.Collect(url, sessionName)

			val, err := gaugeValue(exporter.DronePortRndTrip, "1", "home", "remote station", url, sessionName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(264)))
		})
		It("sets the 'drone_port_power' metric with the right labels", func() {
			collector.Collect(url, sessionName)

			val, _ := gaugeValue(exporter.DronePortPower, "1", url, sessionName)
			Expect(val).To(Equal(200.0))
		})
	})
})
