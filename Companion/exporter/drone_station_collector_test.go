package exporter_test

import (
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/m/v2/exporter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DroneStationCollector", func() {
	var collector *exporter.DroneStationCollector

	BeforeEach(func() {
		FRMServer.Reset()
		collector = exporter.NewDroneStationCollector("http://localhost:9080/getDroneStation")

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
				LatestRndTrip:    "00:04:24",
				LatestTripIncAmt: 82,
				LatestTripOutAmt: 50,
				EstBatteryRate:   30,
			},
		})
	})

	AfterEach(func() {
		collector = nil
	})

	Describe("Drone metrics collection", func() {
		It("sets the 'drone_port_battery_rate' metric with the right labels", func() {
			collector.Collect()

			val, err := gaugeValue(exporter.DronePortBatteryRate, "1", "home", "remote station")

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(30)))
		})
		It("sets the 'drone_port_round_trip_seconds' metric with the right labels", func() {
			collector.Collect()

			val, err := gaugeValue(exporter.DronePortRndTrip, "1", "home", "remote station")

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(264)))
		})
	})
})
