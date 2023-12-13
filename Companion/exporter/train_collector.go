package exporter

import (
	"log"
	"time"
)

type TrainCollector struct {
	FRMAddress    string
	TrackedTrains map[string]*TrainDetails
}

type TimeTable struct {
	StationName string `json:"StationName"`
}

type TrainCar struct {
	Name           string  `json:"Name"`
	TotalMass      float64 `json:"TotalMass"`
	PayloadMass    float64 `json:"PayloadMass"`
	MaxPayloadMass float64 `json:"MaxPayloadMass"`
}

type TrainDetails struct {
	TrainName        string      `json:"Name"`
	PowerConsumed    float64     `json:"PowerConsumed"`
	TrainStation     string      `json:"TrainStation"`
	Derailed         bool        `json:"Derailed"`
	Status           string      `json:"Status"` //"Self-Driving",
	TimeTable        []TimeTable `json:"TimeTable"`
	TrainConsist     []TrainCar  `json:"TrainConsist"`
	ArrivalTime      time.Time
	StationCounter   int
	FirstArrivalTime time.Time
}

func NewTrainCollector(frmAddress string) *TrainCollector {
	return &TrainCollector{
		FRMAddress:    frmAddress,
		TrackedTrains: make(map[string]*TrainDetails),
	}
}

func (t *TrainDetails) recordRoundTripTime(now time.Time) {
	if len(t.TimeTable) <= t.StationCounter {
		roundTripSeconds := now.Sub(t.FirstArrivalTime).Seconds()
		TrainRoundTrip.WithLabelValues(t.TrainName).Set(roundTripSeconds)
		t.StationCounter = 0
		t.FirstArrivalTime = now
	}
}

func (t *TrainDetails) recordSegmentTripTime(destination string, now time.Time) {
	tripSeconds := now.Sub(t.ArrivalTime).Seconds()
	TrainSegmentTrip.WithLabelValues(t.TrainName, t.TrainStation, destination).Set(tripSeconds)
}

func (t *TrainDetails) recordNextStation(d *TrainDetails) {
	if t.TrainStation != d.TrainStation {
		t.StationCounter = t.StationCounter + 1
		now := Clock.Now()
		t.recordSegmentTripTime(d.TrainStation, now)
		t.recordRoundTripTime(now)
		t.ArrivalTime = now
		t.TrainStation = d.TrainStation
	}
}

func (t *TrainDetails) markFirstStation(d *TrainDetails) {
	if t.TrainStation != d.TrainStation {
		t.StationCounter = 0
		t.FirstArrivalTime = Clock.Now()
		t.ArrivalTime = Clock.Now()
		t.TrainStation = d.TrainStation
	}
}

func (t *TrainDetails) startTracking(trackedTrains map[string]*TrainDetails) {
	trackedTrain := TrainDetails{
		TrainName:      t.TrainName,
		TrainStation:   t.TrainStation,
		StationCounter: 0,
		TimeTable:      t.TimeTable,
	}
	trackedTrains[t.TrainName] = &trackedTrain
}

func (d *TrainDetails) handleTimingUpdates(trackedTrains map[string]*TrainDetails) {
	// track self driving train timing
	if d.Status == "Self-Driving" {
		train, exists := trackedTrains[d.TrainName]
		if exists && !train.FirstArrivalTime.IsZero() {
			train.recordNextStation(d)
		} else if exists {
			train.markFirstStation(d)
		} else {
			d.startTracking(trackedTrains)
		}
	} else {
		//remove manual trains, nothing to mark
		_, exists := trackedTrains[d.TrainName]
		if exists {
			delete(trackedTrains, d.TrainName)
		}
	}
}

func (c *TrainCollector) Collect() {
	details := []TrainDetails{}
	err := retrieveData(c.FRMAddress, &details)
	if err != nil {
		log.Printf("error reading train statistics from FRM: %s\n", err)
		return
	}

	for _, d := range details {
		totalMass := 0.0
		payloadMass := 0.0
		maxPayloadMass := 0.0
		locomotives := 0.0

		for _, car := range d.TrainConsist {
			if car.Name == "Electric Locomotive" {
				locomotives = locomotives + 1
			}
			totalMass = totalMass + car.TotalMass
			payloadMass = payloadMass + car.PayloadMass
			maxPayloadMass = maxPayloadMass + car.MaxPayloadMass
		}

		TrainPower.WithLabelValues(d.TrainName).Set(d.PowerConsumed * locomotives) // for now, the total power consumed is a multiple of the reported power consumed by the number of locomotives
		TrainTotalMass.WithLabelValues(d.TrainName).Set(totalMass)
		TrainPayloadMass.WithLabelValues(d.TrainName).Set(payloadMass)
		TrainMaxPayloadMass.WithLabelValues(d.TrainName).Set(maxPayloadMass)

		isDerailed := parseBool(d.Derailed)
		TrainDerailed.WithLabelValues(d.TrainName).Set(isDerailed)

		d.handleTimingUpdates(c.TrackedTrains)
	}
}
