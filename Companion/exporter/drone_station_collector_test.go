package exporter_test

import (
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/exporter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DroneStationCollector", func() {
	var collector *exporter.DroneStationCollector
	var url string
	var sessionName = "default"

	BeforeEach(func() {
		FRMServer.Reset()
		url = FRMServer.server.URL
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
				PowerInfo: exporter.PowerInfo{
					CircuitGroupId: 1.0,
					PowerConsumed:  100,
				},
				Fuel: []exporter.DroneFuelInventory{
					exporter.DroneFuelInventory{
						Name:   "Battery",
						Amount: 200,
					},
				},
				ActiveFuel: exporter.DroneActiveFuel{
					Name: "Battery",
					Rate: 20,
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
				PowerInfo: exporter.PowerInfo{
					CircuitGroupId: 1.0,
					PowerConsumed:  100,
				},
				Fuel: []exporter.DroneFuelInventory{
					exporter.DroneFuelInventory{
						Name:   "Battery",
						Amount: 200,
					},
				},
				ActiveFuel: exporter.DroneActiveFuel{
					Name: "Battery",
					Rate: 30,
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

			val, err := gaugeValue(exporter.DronePortFuelRate, "1", "home", "Battery", url, sessionName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(20)))
			val, err = gaugeValue(exporter.DronePortFuelAmount, "1", "home", "Battery", url, sessionName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(200)))
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
