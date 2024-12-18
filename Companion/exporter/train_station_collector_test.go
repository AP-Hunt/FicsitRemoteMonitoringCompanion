package exporter_test

import (
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/exporter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TrainStationCollector", func() {
	var url string
	var sessionName = "default"
	var collector *exporter.TrainStationCollector

	BeforeEach(func() {
		FRMServer.Reset()
		url = FRMServer.server.URL
		collector = exporter.NewTrainStationCollector("/getTrainStation")

		FRMServer.ReturnsTrainStationData([]exporter.TrainStationDetails{
			{
				Name: "test1",
				CargoPlatforms: []exporter.CargoPlatform{
					{
						LoadingDock:   "Freight Platform",
						TransferRate:  66,
						LoadingStatus: "Idle",
						LoadingMode:   "Loading",
						PowerInfo: exporter.PowerInfo{
							CircuitGroupId:   1,
							PowerConsumed:    0.1,
							MaxPowerConsumed: 50,
						},
					},
					{
						LoadingDock:   "Freight Platform",
						TransferRate:  66,
						LoadingStatus: "Idle",
						LoadingMode:   "Loading",
						PowerInfo: exporter.PowerInfo{
							CircuitGroupId:   1,
							PowerConsumed:    0.1,
							MaxPowerConsumed: 50,
						},
					},
				},
				PowerInfo: exporter.PowerInfo{
					CircuitGroupId:   1,
					PowerConsumed:    0.1,
					MaxPowerConsumed: 50,
				},
			},
			{
				Name: "test2",
				CargoPlatforms: []exporter.CargoPlatform{
					{
						LoadingDock:   "Freight Platform",
						TransferRate:  66,
						LoadingStatus: "Unloading",
						LoadingMode:   "Unloading",
						PowerInfo: exporter.PowerInfo{
							CircuitGroupId:   1,
							PowerConsumed:    50,
							MaxPowerConsumed: 50,
						},
					},
					{
						LoadingDock:   "Freight Platform",
						TransferRate:  66,
						LoadingStatus: "Idle",
						LoadingMode:   "Unloading",
						PowerInfo: exporter.PowerInfo{
							CircuitGroupId:   1,
							PowerConsumed:    0.1,
							MaxPowerConsumed: 50,
						},
					},
				},
				PowerInfo: exporter.PowerInfo{
					CircuitGroupId:   1,
					PowerConsumed:    0.1,
					MaxPowerConsumed: 50,
				},
			},
		})
	})

	AfterEach(func() {
		collector = nil
	})

	Describe("Train station metrics collection", func() {
		It("sets the 'train_station_power' metric with the right labels", func() {
			collector.Collect(url, sessionName)

			val, err := gaugeValue(exporter.TrainStationPower, "1", url, sessionName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(50.5)))
		})

		It("sets the 'train_station_power_max' metric with the right labels", func() {
			collector.Collect(url, sessionName)

			val, err := gaugeValue(exporter.TrainStationPowerMax, "1", url, sessionName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(300.0)))
		})
	})
})
