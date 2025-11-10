package exporter

import (
	"log"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type TrainCollector struct {
	endpoint       string
	TrackedTrains  map[string]*TrainDetails
	metricsDropper *MetricsDropper
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
	TrainStation     string      `json:"TrainStation"`
	Derailed         bool        `json:"Derailed"`
	Status           string      `json:"Status"` //"Self-Driving",
	TimeTable        []TimeTable `json:"TimeTable"`
	TrainCars        []TrainCar  `json:"Vehicles"`
	PowerInfo        PowerInfo   `json:"PowerInfo"`
	ArrivalTime      time.Time
	StationCounter   int
	FirstArrivalTime time.Time
	LastTracked      time.Time
}

func NewTrainCollector(endpoint string) *TrainCollector {
	return &TrainCollector{
		endpoint:      endpoint,
		TrackedTrains: make(map[string]*TrainDetails),
		metricsDropper: NewMetricsDropper(
			TrainRoundTrip,
			TrainSegmentTrip,
			TrainPower,
			TrainTotalMass,
			TrainPayloadMass,
			TrainMaxPayloadMass,
			TrainDerailed,
		),
	}
}

func (t *TrainDetails) recordRoundTripTime(now time.Time, frmAddress string, sessionName string) {
	if len(t.TimeTable) <= t.StationCounter {
		roundTripSeconds := now.Sub(t.FirstArrivalTime).Seconds()
		TrainRoundTrip.WithLabelValues(t.TrainName, frmAddress, sessionName).Set(roundTripSeconds)
		t.StationCounter = 0
		t.FirstArrivalTime = now
	}
}

func (t *TrainDetails) recordSegmentTripTime(destination string, now time.Time, frmAddress string, sessionName string) {
	tripSeconds := now.Sub(t.ArrivalTime).Seconds()
	TrainSegmentTrip.WithLabelValues(t.TrainName, t.TrainStation, destination, frmAddress, sessionName).Set(tripSeconds)
}

func (t *TrainDetails) recordNextStation(d *TrainDetails, frmAddress string, sessionName string) {
	if t.TrainStation != d.TrainStation {
		t.StationCounter = t.StationCounter + 1
		now := Clock.Now()
		t.recordSegmentTripTime(d.TrainStation, now, frmAddress, sessionName)
		t.recordRoundTripTime(now, frmAddress, sessionName)
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
		LastTracked:    Clock.Now(),
	}
	trackedTrains[t.TrainName] = &trackedTrain
}

func (d *TrainDetails) handleTimingUpdates(trackedTrains map[string]*TrainDetails, frmAddress string, sessionName string) {
	// track self driving train timing
	if d.Status == "Self-Driving" {
		train, exists := trackedTrains[d.TrainName]
		if !exists || train.LastTracked.Before(Clock.Now().Add(-time.Minute)) {
			d.startTracking(trackedTrains)
		} else if exists && !train.FirstArrivalTime.IsZero() {
			train.recordNextStation(d, frmAddress, sessionName)
		} else if exists {
			train.markFirstStation(d)
		}

		// mark that we saw this vehicle
		if exists {
			train.LastTracked = Clock.Now()
		}
	} else {
		//remove manual trains, nothing to mark
		_, exists := trackedTrains[d.TrainName]
		if exists {
			delete(trackedTrains, d.TrainName)
		}
	}
}

func (c *TrainCollector) Collect(frmAddress string, sessionName string) {
	details := []TrainDetails{}
	err := retrieveData(frmAddress, c.endpoint, &details)
	if err != nil {
		c.metricsDropper.DropStaleMetricLabels()
		log.Printf("error reading train statistics from FRM: %s\n", err)
		return
	}

	powerInfo := map[float64]float64{}
	maxPowerInfo := map[float64]float64{}
	for _, d := range details {
		c.metricsDropper.CacheFreshMetricLabel(prometheus.Labels{"url": frmAddress, "session_name": sessionName, "name": d.TrainName})
		totalMass := 0.0
		payloadMass := 0.0
		maxPayloadMass := 0.0
		locomotives := 0.0

		for _, car := range d.TrainCars {
			if car.Name == "Electric Locomotive" {
				locomotives = locomotives + 1
			}
			totalMass = totalMass + car.TotalMass
			payloadMass = payloadMass + car.PayloadMass
			maxPayloadMass = maxPayloadMass + car.MaxPayloadMass
		}

		// for now, the total power consumed is a multiple of the reported power consumed by the number of locomotives
		trainPowerConsumed := d.PowerInfo.PowerConsumed * locomotives
		maxTrainPowerConsumed := d.PowerInfo.MaxPowerConsumed * locomotives

		TrainPower.WithLabelValues(d.TrainName, frmAddress, sessionName).Set(trainPowerConsumed)
		TrainTotalMass.WithLabelValues(d.TrainName, frmAddress, sessionName).Set(totalMass)
		TrainPayloadMass.WithLabelValues(d.TrainName, frmAddress, sessionName).Set(payloadMass)
		TrainMaxPayloadMass.WithLabelValues(d.TrainName, frmAddress, sessionName).Set(maxPayloadMass)

		isDerailed := parseBool(d.Derailed)
		TrainDerailed.WithLabelValues(d.TrainName, frmAddress, sessionName).Set(isDerailed)

		d.handleTimingUpdates(c.TrackedTrains, frmAddress, sessionName)

		circuitGroupId := d.PowerInfo.CircuitGroupId
		val, ok := powerInfo[circuitGroupId]
		if ok {
			powerInfo[circuitGroupId] = val + trainPowerConsumed
		} else {
			powerInfo[circuitGroupId] = trainPowerConsumed
		}
		val, ok = maxPowerInfo[circuitGroupId]
		if ok {
			maxPowerInfo[circuitGroupId] = val + maxTrainPowerConsumed
		} else {
			maxPowerInfo[circuitGroupId] = maxTrainPowerConsumed
		}
	}
	for circuitGroupId, powerConsumed := range powerInfo {
		cid := strconv.FormatFloat(circuitGroupId, 'f', -1, 64)
		TrainCircuitPower.WithLabelValues(cid, frmAddress, sessionName).Set(powerConsumed)
	}
	for circuitGroupId, powerConsumed := range maxPowerInfo {
		cid := strconv.FormatFloat(circuitGroupId, 'f', -1, 64)
		TrainCircuitPowerMax.WithLabelValues(cid, frmAddress, sessionName).Set(powerConsumed)
	}
	c.metricsDropper.DropStaleMetricLabels()
}

func (c *TrainCollector) DropCache() {
	c.TrackedTrains = map[string]*TrainDetails{}
}
