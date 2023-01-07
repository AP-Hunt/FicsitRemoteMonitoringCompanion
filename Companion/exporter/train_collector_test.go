package exporter_test

import (
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/m/v2/exporter"
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
			Status:        "TS_SelfDriving",
			TimeTable: []exporter.TimeTable{
				{StationName: "First"},
				{StationName: "Second"},
				{StationName: "Third"},
			},
		},
		{
			TrainName:     "Not In Use",
			PowerConsumed: 0,
			TrainStation:  "Offsite",
			Derailed:      false,
			Status:        "TS_Parked",
			TimeTable: []exporter.TimeTable{
				{StationName: "Offsite"},
			},
		},
	})
}

var _ = Describe("TrainCollector", func() {
	var collector *exporter.TrainCollector

	BeforeEach(func() {
		FRMServer.Reset()
		collector = exporter.NewTrainCollector("http://localhost:9080/getTrains")

		FRMServer.ReturnsTrainData([]exporter.TrainDetails{
			{
				TrainName:     "Train1",
				PowerConsumed: 67,
				TrainStation:  "NextStation",
				Derailed:      false,
				Status:        "TS_SelfDriving",
				TimeTable: []exporter.TimeTable{
					{StationName: "First"},
					{StationName: "Second"},
				},
			},
			{
				TrainName:     "DerailedTrain",
				PowerConsumed: 0,
				TrainStation:  "NextStation",
				Derailed:      true,
				Status:        "Derailed",
			},
		})
	})

	AfterEach(func() {
		collector = nil
	})

	Describe("Train metrics collection", func() {
		It("sets the 'train_derailed' metric with the right labels", func() {
			collector.Collect()

			val, err := gaugeValue(exporter.TrainDerailed, "DerailedTrain")

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(1)))
		})

		It("sets the 'train_power_consumed' metric with the right labels", func() {
			collector.Collect()

			val, err := gaugeValue(exporter.TrainPower, "Train1")

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(67)))
		})

		It("sets 'train_segment_trip_seconds' metric with the right labels", func() {

			testTime := clock.NewMock()
			exporter.Clock = testTime
			updateTrain("First")

			collector.Collect()
			val, err := gaugeValue(exporter.TrainSegmentTrip, "Train1", "First", "Second")
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(0)))
			testTime.Add(5 * time.Second)
			collector.Collect()
			val, err = gaugeValue(exporter.TrainSegmentTrip, "Train1", "First", "Second")
			Expect(val).To(Equal(float64(0)))
			testTime.Add(25 * time.Second)

			// Start timing the trains here - No metrics yet because we just got our first "start" marker from the station change.
			updateTrain("Second")
			collector.Collect()
			val, err = gaugeValue(exporter.TrainSegmentTrip, "Train1", "First", "Second")
			Expect(val).To(Equal(float64(0)))

			testTime.Add(15 * time.Second)
			collector.Collect()
			testTime.Add(10 * time.Second)
			collector.Collect()
			// No stats again since train is still "en route"
			val, err = gaugeValue(exporter.TrainSegmentTrip, "Train1", "First", "Second")
			Expect(val).To(Equal(float64(0)))

			testTime.Add(5 * time.Second)

			// Can record elapsed time between Second and Third stations
			updateTrain("Third")
			collector.Collect()
			val, err = gaugeValue(exporter.TrainSegmentTrip, "Train1", "Second", "Third")
			Expect(val).To(Equal(float64(30)))

			testTime.Add(30 * time.Second)
			updateTrain("First")
			collector.Collect()
			val, err = gaugeValue(exporter.TrainSegmentTrip, "Train1", "Third", "First")
			Expect(val).To(Equal(float64(30)))

			testTime.Add(30 * time.Second)
			updateTrain("Second")
			collector.Collect()

			val, err = gaugeValue(exporter.TrainSegmentTrip, "Train1", "First", "Second")
			Expect(val).To(Equal(float64(30)))

		})

		It("sets 'train_round_trip_seconds' metric with the right labels", func() {
			testTime := clock.NewMock()
			exporter.Clock = testTime
			updateTrain("Third")

			collector.Collect()
			val, err := gaugeValue(exporter.TrainRoundTrip, "Train1")
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(0)))
			testTime.Add(30 * time.Second)

			// Started recording round trip on first station arrival
			updateTrain("First")
			collector.Collect()
			val, err = gaugeValue(exporter.TrainRoundTrip, "Train1")
			Expect(val).To(Equal(float64(0)))

			testTime.Add(30 * time.Second)
			updateTrain("Second")
			collector.Collect()
			testTime.Add(30 * time.Second)
			updateTrain("Third")
			collector.Collect()
			testTime.Add(30 * time.Second)
			updateTrain("First")
			collector.Collect()

			val, err = gaugeValue(exporter.TrainRoundTrip, "Train1")
			Expect(val).To(Equal(float64(90)))

			//second round trip should also record properly
			testTime.Add(10 * time.Second)
			updateTrain("Second")
			collector.Collect()
			testTime.Add(10 * time.Second)
			updateTrain("Third")
			collector.Collect()
			testTime.Add(10 * time.Second)
			updateTrain("First")
			collector.Collect()

			val, err = gaugeValue(exporter.TrainRoundTrip, "Train1")
			Expect(val).To(Equal(float64(30)))

		})
	})
})
