package exporter

import (
	"log"
	"strconv"
	"time"
)

var MaxTrainPowerConsumption = 110.0

type TrainCollector struct {
	endpoint        string
	TrackedTrains   map[string]*TrainDetails
	TrackedStations *map[string]TrainStationDetails
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

func NewTrainCollector(endpoint string, trackedStations *map[string]TrainStationDetails) *TrainCollector {
	return &TrainCollector{
		endpoint:        endpoint,
		TrackedTrains:   make(map[string]*TrainDetails),
		TrackedStations: trackedStations,
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
	}
	trackedTrains[t.TrainName] = &trackedTrain
}

func (d *TrainDetails) handleTimingUpdates(trackedTrains map[string]*TrainDetails, frmAddress string, sessionName string) {
	// track self driving train timing
	if d.Status == "Self-Driving" {
		train, exists := trackedTrains[d.TrainName]
		if exists && !train.FirstArrivalTime.IsZero() {
			train.recordNextStation(d, frmAddress, sessionName)
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

func (c *TrainCollector) Collect(frmAddress string, sessionName string) {
	details := []TrainDetails{}
	err := retrieveData(frmAddress+c.endpoint, &details)
	if err != nil {
		log.Printf("error reading train statistics from FRM: %s\n", err)
		return
	}

	powerInfo := map[float64]float64{}
	maxPowerInfo := map[float64]float64{}
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

		// for now, the total power consumed is a multiple of the reported power consumed by the number of locomotives
		trainPowerConsumed := d.PowerConsumed * locomotives
		maxTrainPowerConsumed := MaxTrainPowerConsumption * locomotives

		TrainPower.WithLabelValues(d.TrainName, frmAddress, sessionName).Set(trainPowerConsumed)
		TrainTotalMass.WithLabelValues(d.TrainName, frmAddress, sessionName).Set(totalMass)
		TrainPayloadMass.WithLabelValues(d.TrainName, frmAddress, sessionName).Set(payloadMass)
		TrainMaxPayloadMass.WithLabelValues(d.TrainName, frmAddress, sessionName).Set(maxPayloadMass)

		isDerailed := parseBool(d.Derailed)
		TrainDerailed.WithLabelValues(d.TrainName, frmAddress, sessionName).Set(isDerailed)

		d.handleTimingUpdates(c.TrackedTrains, frmAddress, sessionName)

		if len(d.TimeTable) > 0 {
			station, ok := (*c.TrackedStations)[d.TimeTable[0].StationName]
			if ok {
				circuitGroupId := station.PowerInfo.CircuitGroupId
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
}
