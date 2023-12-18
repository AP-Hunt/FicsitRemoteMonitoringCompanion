package exporter

import (
	"log"
	"strconv"
)

var(
	StationPower = 0.1 // should be 50, but currently bugged.
	CargoPlatformPower = 50.0
)

type TrainStationCollector struct {
	FRMAddress      string
	TrackedStations *map[string]TrainStationDetails
}

type CargoPlatform struct {
	LoadingDock   string  `json:"LoadingDock"`
	TransferRate  float64 `json:"TransferRate"`
	LoadingStatus string  `json:"LoadingStatus"` // Idle, Loading, Unloading
	LoadingMode   string  `json:"LoadingMode"`
}

type TrainStationDetails struct {
	Name           string          `json:"Name"`
	Location       Location        `json:"location"`
	CargoPlatforms []CargoPlatform `json:"CargoPlatforms"`
	PowerInfo      PowerInfo       `json:"PowerInfo"`
}

func NewTrainStationCollector(frmAddress string, trackedStations *map[string]TrainStationDetails) *TrainStationCollector {
	return &TrainStationCollector{
		FRMAddress:      frmAddress,
		TrackedStations: trackedStations,
	}
}

func (c *TrainStationCollector) Collect() {
	details := []TrainStationDetails{}
	err := retrieveData(c.FRMAddress, &details)
	if err != nil {
		log.Printf("error reading train station statistics from FRM: %s\n", err)
		return
	}

	powerInfo := map[float64]float64{}
	maxPowerInfo := map[float64]float64{}
	for _, d := range details {
		val, ok := powerInfo[d.PowerInfo.CircuitId]
		maxval, maxok := maxPowerInfo[d.PowerInfo.CircuitId]

		// some additional calculations: for now, power listed here is only for the station.
		// add each of the cargo platforms' power info: 0.1MW if Idle, 50MW otherwise
		totalPowerConsumed := d.PowerInfo.PowerConsumed
		maxTotalPowerConsumed := StationPower
		for _, p := range d.CargoPlatforms {
			maxTotalPowerConsumed = maxTotalPowerConsumed + CargoPlatformPower
			if p.LoadingStatus == "Idle" {
				totalPowerConsumed = totalPowerConsumed + 0.1
			} else {
				totalPowerConsumed = totalPowerConsumed + CargoPlatformPower
			}
		}

		if ok {
			powerInfo[d.PowerInfo.CircuitId] = val + totalPowerConsumed
		} else {
			powerInfo[d.PowerInfo.CircuitId] = totalPowerConsumed
		}

		if maxok {
			maxPowerInfo[d.PowerInfo.CircuitId] = maxval + maxTotalPowerConsumed
		} else {
			maxPowerInfo[d.PowerInfo.CircuitId] = maxTotalPowerConsumed
		}

		//also cache stations so other metrics can figure out a circuit id from a station name
		(*c.TrackedStations)[d.Name] = d
	}
	for circuitId, powerConsumed := range powerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		TrainStationPower.WithLabelValues(cid).Set(powerConsumed)
	}
	for circuitId, powerConsumed := range maxPowerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		TrainStationPowerMax.WithLabelValues(cid).Set(powerConsumed)
	}
}
