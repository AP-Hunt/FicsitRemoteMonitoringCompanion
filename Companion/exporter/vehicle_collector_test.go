package exporter_test

import (
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/exporter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/coder/quartz"
	"time"
)

func updateLocation(x float64, y float64, rotation float64) {
	FRMServer.ReturnsVehicleData([]exporter.VehicleDetails{
		{
			Id:           "1",
			VehicleType:  "Truck",
			ForwardSpeed: 0,
			Location: exporter.Location{
				X:        x,
				Y:        y,
				Z:        0,
				Rotation: rotation,
			},
			AutoPilot: true,
			Fuel: []exporter.Fuel{{Name: "Coal",
				Amount: 23,
			}},
			PathName: "Path",
		},
		{
			Id:           "2",
			VehicleType:  "Truck",
			ForwardSpeed: 0,
			Location: exporter.Location{
				X:        0,
				Y:        0,
				Z:        0,
				Rotation: rotation,
			},
			AutoPilot: false,
			Fuel: []exporter.Fuel{{Name: "Coal",
				Amount: 23,
			}},
			PathName: "no path",
		},
	})
}

var _ = Describe("VehicleCollector", func() {
	var collector *exporter.VehicleCollector
	var url string
	var sessionName = "default"

	BeforeEach(func() {
		FRMServer.Reset()
		url = FRMServer.server.URL
		collector = exporter.NewVehicleCollector("/getVehicles")

		FRMServer.ReturnsVehicleData([]exporter.VehicleDetails{
			{
				Id:           "1",
				VehicleType:  "Truck",
				ForwardSpeed: 80,
				Location: exporter.Location{
					X:        1000,
					Y:        2000,
					Z:        1000,
					Rotation: 60,
				},
				AutoPilot: true,
				Fuel: []exporter.Fuel{{Name: "Coal",
					Amount: 23,
				}},
				PathName: "Path",
			},
		})
	})

	AfterEach(func() {
		collector = nil
	})

	Describe("Vehicle metrics collection", func() {
		It("sets the 'vehicle_fuel' metric with the right labels", func() {
			collector.Collect(url, sessionName)

			val, err := gaugeValue(exporter.VehicleFuel, "1", "Truck", "Coal", url, sessionName)

			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(23)))
		})

		It("sets the 'vehicle_round_trip_seconds' metric with the right labels", func() {

			testTime := quartz.NewMock(GinkgoTB())
			exporter.Clock = testTime
			// truck will be too fast here, nothing recorded
			collector.Collect(url, sessionName)
			val, err := gaugeValue(exporter.VehicleRoundTrip, "1", "Truck", "Path", url, sessionName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(0)))

			testTime.Advance(30 * time.Second)
			updateLocation(0, 0, 0)
			// first time collecting stats, nothing yet but it does set start location to 0,0,0
			collector.Collect(url, sessionName)
			val, err = gaugeValue(exporter.VehicleRoundTrip, "1", "Truck", "Path", url, sessionName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(0)))

			testTime.Advance(30 * time.Second)
			updateLocation(8000, 2000, 0)
			//go to a far away location now, star the timer
			collector.Collect(url, sessionName)
			val, err = gaugeValue(exporter.VehicleRoundTrip, "1", "Truck", "Path", url, sessionName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(0)))

			testTime.Advance(10 * time.Second)
			updateLocation(1000, 2000, 180)
			//We are near but not facing the right way. Do not record this until we face near the right direction
			collector.Collect(url, sessionName)
			val, err = gaugeValue(exporter.VehicleRoundTrip, "1", "Truck", "Path", url, sessionName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(0)))

			testTime.Advance(20 * time.Second)
			updateLocation(1000, 2000, 0)
			//Now we are back near enough to where we began recording, and facing near the same way end recording
			collector.Collect(url, sessionName)
			val, err = gaugeValue(exporter.VehicleRoundTrip, "1", "Truck", "Path", url, sessionName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(30)))
		})

		It("Does not track if there's been a long time since the last pull", func() {

			testTime := quartz.NewMock(GinkgoTB())
			exporter.Clock = testTime
			updateLocation(0, 0, 0)
			// first time collecting stats, nothing yet but it does set start location to 0,0,0
			collector.Collect(url, sessionName)
			val, err := gaugeValue(exporter.VehicleRoundTrip, "1", "Truck", "Path", url, sessionName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(0)))

			testTime.Advance(30 * time.Second)
			updateLocation(8000, 2000, 0)
			//go to a far away location now, star the timer
			collector.Collect(url, sessionName)
			val, err = gaugeValue(exporter.VehicleRoundTrip, "1", "Truck", "Path", url, sessionName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(0)))

			// A long time passes before we collect again... too long! so nothing is tracked
			testTime.Advance(120 * time.Second)
			updateLocation(1000, 2000, 0)
			collector.Collect(url, sessionName)
			val, err = gaugeValue(exporter.VehicleRoundTrip, "1", "Truck", "Path", url, sessionName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(0)))
		})
	})
})
