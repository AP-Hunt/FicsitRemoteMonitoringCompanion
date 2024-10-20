package exporter_test

import (
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/exporter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TrainStationCollector", func() {
	var url = "http://localhost:9080"
	var saveName = "default"
	var collector *exporter.TrainStationCollector

	BeforeEach(func() {
		FRMServer.Reset()
		trackedStations := &(map[string]exporter.TrainStationDetails{})
		collector = exporter.NewTrainStationCollector("/getTrainStation", trackedStations)

		FRMServer.ReturnsTrainStationData([]exporter.TrainStationDetails{
			{
				Name: "test1",
				CargoPlatforms: []exporter.CargoPlatform{
					{
						LoadingDock: "Freight Platform",
						TransferRate: 66,
						LoadingStatus: "Idle",
						LoadingMode: "Loading",
					},
					{
						LoadingDock: "Freight Platform",
						TransferRate: 66,
						LoadingStatus: "Idle",
						LoadingMode: "Loading",
					},
				},
				PowerInfo: exporter.PowerInfo{
					CircuitId: 1,
					PowerConsumed: 50,
				},
			},
			{
				Name: "test2",
				CargoPlatforms: []exporter.CargoPlatform{
					{
						LoadingDock: "Freight Platform",
						TransferRate: 66,
						LoadingStatus: "Unloading",
						LoadingMode: "Unloading",
					},
					{
						LoadingDock: "Freight Platform",
						TransferRate: 66,
						LoadingStatus: "Idle",
						LoadingMode: "Unloading",
					},
				},
				PowerInfo: exporter.PowerInfo{
					CircuitId: 1,
					PowerConsumed: 50,
				},
			},
		})
	})

	AfterEach(func() {
		collector = nil
	})

	Describe("Train station metrics collection", func() {
		It("sets the 'train_station_power' metric with the right labels", func() {
			collector.Collect(url, saveName)

			val, err := gaugeValue(exporter.TrainStationPower, "1", url, saveName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(150.3)))
		})

		It("sets the 'train_station_power_max' metric with the right labels", func() {
			collector.Collect(url, saveName)

			val, err := gaugeValue(exporter.TrainStationPowerMax, "1", url, saveName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(200.2)))
		})
	})
})
