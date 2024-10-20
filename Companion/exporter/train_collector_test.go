package exporter_test

import (
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/exporter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/benbjohnson/clock"
	"time"
)

func updateTrain(station string) {
	FRMServer.ReturnsTrainData([]exporter.TrainDetails{
		{
			TrainName:     "Train1",
			PowerConsumed: 0,
			TrainStation:  station,
			Derailed:      false,
			Status:        "Self-Driving",
			TimeTable: []exporter.TimeTable{
				{StationName: "First"},
				{StationName: "Second"},
				{StationName: "Third"},
			},
			TrainConsist: []exporter.TrainCar{
				{Name: "Electric Locomotive", TotalMass: 3000, PayloadMass: 0, MaxPayloadMass: 0},
				{Name: "Freight Car", TotalMass: 47584, PayloadMass: 17584, MaxPayloadMass: 70000},
			},
		},
		{
			TrainName:     "Not In Use",
			PowerConsumed: 0,
			TrainStation:  "Offsite",
			Derailed:      false,
			Status:        "Parked",
			TimeTable: []exporter.TimeTable{
				{StationName: "Offsite"},
			},
			TrainConsist: []exporter.TrainCar{
				{Name: "Electric Locomotive", TotalMass: 3000, PayloadMass: 0, MaxPayloadMass: 0},
				{Name: "Freight Car", TotalMass: 47584, PayloadMass: 17584, MaxPayloadMass: 70000},
			},
		},
	})
}

var _ = Describe("TrainCollector", func() {
	var collector *exporter.TrainCollector
	var url = "http://localhost:9080"
	var saveName = "default"

	BeforeEach(func() {
		FRMServer.Reset()
		trackedStations := &(map[string]exporter.TrainStationDetails{
			"First": {
				Name: "First",
				PowerInfo: exporter.PowerInfo{
					CircuitId: 1,
				},
			},
			"Second": {
				Name: "Second",
				PowerInfo: exporter.PowerInfo{
					CircuitId: 1,
				},
			},
			"Third": {
				Name: "Third",
				PowerInfo: exporter.PowerInfo{
					CircuitId: 1,
				},
			},
		})
		collector = exporter.NewTrainCollector("/getTrains", trackedStations)

		FRMServer.ReturnsTrainData([]exporter.TrainDetails{
			{
				TrainName:     "Train1",
				PowerConsumed: 67,
				TrainStation:  "NextStation",
				Derailed:      false,
				Status:        "Self-Driving",
				TimeTable: []exporter.TimeTable{
					{StationName: "First"},
					{StationName: "Second"},
				},
				TrainConsist: []exporter.TrainCar{
					{Name: "Electric Locomotive", TotalMass: 3000, PayloadMass: 0, MaxPayloadMass: 0},
					{Name: "Electric Locomotive", TotalMass: 3000, PayloadMass: 0, MaxPayloadMass: 0},
					{Name: "Freight Car", TotalMass: 47584, PayloadMass: 17584, MaxPayloadMass: 70000},
					{Name: "Freight Car", TotalMass: 47584, PayloadMass: 17584, MaxPayloadMass: 70000},
				},
			},
			{
				TrainName:     "Train2",
				PowerConsumed: 22,
				TrainStation:  "NextStation",
				Derailed:      false,
				Status:        "Self-Driving",
				TimeTable: []exporter.TimeTable{
					{StationName: "Second"},
					{StationName: "Third"},
				},
				TrainConsist: []exporter.TrainCar{
					{Name: "Electric Locomotive", TotalMass: 3000, PayloadMass: 0, MaxPayloadMass: 0},
					{Name: "Freight Car", TotalMass: 47584, PayloadMass: 17584, MaxPayloadMass: 70000},
					{Name: "Freight Car", TotalMass: 47584, PayloadMass: 17584, MaxPayloadMass: 70000},
				},
			},
			{
				TrainName:     "DerailedTrain",
				PowerConsumed: 0,
				TrainStation:  "NextStation",
				Derailed:      true,
				Status:        "Derailed",
				TrainConsist:  []exporter.TrainCar{},
			},
		})
	})

	AfterEach(func() {
		collector = nil
	})

	Describe("Train metrics collection", func() {
		It("sets the 'train_derailed' metric with the right labels", func() {
			collector.Collect(url, saveName)

			val, err := gaugeValue(exporter.TrainDerailed, "DerailedTrain", url, saveName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(1)))
		})

		It("sets the 'train_power_consumed' metric with the right labels", func() {
			collector.Collect(url, saveName)

			val, err := gaugeValue(exporter.TrainPower, "Train1", url, saveName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(67 * 2))) //expects reported power to be per train, so multiply by # of trains
		})

		It("sets the 'train_power_circuit_consumed' metric with the right labels", func() {
			collector.Collect(url, saveName)

			val, err := gaugeValue(exporter.TrainCircuitPower, "1", url, saveName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal((67.0 * 2) + 22.0))
		})

		It("sets the 'train_power_circuit_consumed_max' metric with the right labels", func() {
			collector.Collect(url, saveName)

			val, err := gaugeValue(exporter.TrainCircuitPowerMax, "1", url, saveName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(exporter.MaxTrainPowerConsumption * 3))
		})

		It("sets the mass metrics with the right labels", func() {
			collector.Collect(url, saveName)

			val, _ := gaugeValue(exporter.TrainTotalMass, "Train1", url, saveName)
			Expect(val).To(Equal(3000.0 + 3000.0 + 47584.0 + 47584.0))
			val, _ = gaugeValue(exporter.TrainPayloadMass, "Train1", url, saveName)
			Expect(val).To(Equal(17584.0 + 17584.0))
			val, _ = gaugeValue(exporter.TrainMaxPayloadMass, "Train1", url, saveName)
			Expect(val).To(Equal(70000.0 * 2))
		})

		It("sets 'train_segment_trip_seconds' metric with the right labels", func() {

			testTime := clock.NewMock()
			exporter.Clock = testTime
			updateTrain("First")

			collector.Collect(url, saveName)
			val, err := gaugeValue(exporter.TrainSegmentTrip, "Train1", "First", "Second", url, saveName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(0)))
			testTime.Add(5 * time.Second)
			collector.Collect(url, saveName)
			val, err = gaugeValue(exporter.TrainSegmentTrip, "Train1", "First", "Second", url, saveName)
			Expect(val).To(Equal(float64(0)))
			testTime.Add(25 * time.Second)

			// Start timing the trains here - No metrics yet because we just got our first "start" marker from the station change.
			updateTrain("Second")
			collector.Collect(url, saveName)
			val, err = gaugeValue(exporter.TrainSegmentTrip, "Train1", "First", "Second", url, saveName)
			Expect(val).To(Equal(float64(0)))

			testTime.Add(15 * time.Second)
			collector.Collect(url, saveName)
			testTime.Add(10 * time.Second)
			collector.Collect(url, saveName)
			// No stats again since train is still "en route"
			val, err = gaugeValue(exporter.TrainSegmentTrip, "Train1", "First", "Second", url, saveName)
			Expect(val).To(Equal(float64(0)))

			testTime.Add(5 * time.Second)

			// Can record elapsed time between Second and Third stations
			updateTrain("Third")
			collector.Collect(url, saveName)
			val, err = gaugeValue(exporter.TrainSegmentTrip, "Train1", "Second", "Third", url, saveName)
			Expect(val).To(Equal(float64(30)))

			testTime.Add(30 * time.Second)
			updateTrain("First")
			collector.Collect(url, saveName)
			val, err = gaugeValue(exporter.TrainSegmentTrip, "Train1", "Third", "First", url, saveName)
			Expect(val).To(Equal(float64(30)))

			testTime.Add(30 * time.Second)
			updateTrain("Second")
			collector.Collect(url, saveName)

			val, err = gaugeValue(exporter.TrainSegmentTrip, "Train1", "First", "Second", url, saveName)
			Expect(val).To(Equal(float64(30)))

		})

		It("sets 'train_round_trip_seconds' metric with the right labels", func() {
			testTime := clock.NewMock()
			exporter.Clock = testTime
			updateTrain("Third")

			collector.Collect(url, saveName)
			val, err := gaugeValue(exporter.TrainRoundTrip, "Train1", url, saveName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(0)))
			testTime.Add(30 * time.Second)

			// Started recording round trip on first station arrival
			updateTrain("First")
			collector.Collect(url, saveName)
			val, err = gaugeValue(exporter.TrainRoundTrip, "Train1", url, saveName)
			Expect(val).To(Equal(float64(0)))

			testTime.Add(30 * time.Second)
			updateTrain("Second")
			collector.Collect(url, saveName)
			testTime.Add(30 * time.Second)
			updateTrain("Third")
			collector.Collect(url, saveName)
			testTime.Add(30 * time.Second)
			updateTrain("First")
			collector.Collect(url, saveName)

			val, err = gaugeValue(exporter.TrainRoundTrip, "Train1", url, saveName)
			Expect(val).To(Equal(float64(90)))

			//second round trip should also record properly
			testTime.Add(10 * time.Second)
			updateTrain("Second")
			collector.Collect(url, saveName)
			testTime.Add(10 * time.Second)
			updateTrain("Third")
			collector.Collect(url, saveName)
			testTime.Add(10 * time.Second)
			updateTrain("First")
			collector.Collect(url, saveName)

			val, err = gaugeValue(exporter.TrainRoundTrip, "Train1", url, saveName)
			Expect(val).To(Equal(float64(30)))

		})
	})
})
