package exporter_test

import (
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/exporter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("VehicleStationCollector", func() {
	var collector *exporter.VehicleStationCollector

	BeforeEach(func() {
		FRMServer.Reset()
		collector = exporter.NewVehicleStationCollector("http://localhost:9080/getTruckStation")

		FRMServer.ReturnsVehicleStationData([]exporter.VehicleStationDetails{
			{
				Name: "Truck Station",
				PowerInfo: exporter.PowerInfo{
					CircuitId:     1,
					PowerConsumed: 20,
				},
			},
			{
				Name: "Truck Station",
				PowerInfo: exporter.PowerInfo{
					CircuitId:     1,
					PowerConsumed: 20,
				},
			},
			{
				Name: "Truck Station",
				PowerInfo: exporter.PowerInfo{
					CircuitId:     1,
					PowerConsumed: 0.1,
				},
			},
		})
	})

	AfterEach(func() {
		collector = nil
	})

	Describe("Truck station metrics collection", func() {
		It("sets the 'vehicle_station_power' metric with the right labels", func() {
			collector.Collect()

			val, err := gaugeValue(exporter.VehicleStationPower, "1")

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(40.1))
		})
		It("sets the 'vehicle_station_power_max' metric with the right labels", func() {
			collector.Collect()

			val, err := gaugeValue(exporter.VehicleStationPowerMax, "1")

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(60.0))
		})
	})
})
