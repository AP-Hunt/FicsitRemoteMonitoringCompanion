package exporter

import (
	"log"
	"strconv"
)

var (
	StationPower       = 0.1 // should be 50, but currently bugged.
	CargoPlatformPower = 50.0
)

// TODO: drop tracked stations when save game is updated?
type TrainStationCollector struct {
	endpoint        string
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

func NewTrainStationCollector(endpoint string, trackedStations *map[string]TrainStationDetails) *TrainStationCollector {
	return &TrainStationCollector{
		endpoint:        endpoint,
		TrackedStations: trackedStations,
	}
}

func (c *TrainStationCollector) Collect(frmAddress string, saveName string) {
	details := []TrainStationDetails{}
	err := retrieveData(frmAddress+c.endpoint, &details)
	if err != nil {
		log.Printf("error reading train station statistics from FRM: %s\n", err)
		return
	}

	powerInfo := map[float64]float64{}
	maxPowerInfo := map[float64]float64{}
	for _, d := range details {
		val, ok := powerInfo[d.PowerInfo.CircuitGroupId]
		maxval, maxok := maxPowerInfo[d.PowerInfo.CircuitGroupId]

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
			powerInfo[d.PowerInfo.CircuitGroupId] = val + totalPowerConsumed
		} else {
			powerInfo[d.PowerInfo.CircuitGroupId] = totalPowerConsumed
		}

		if maxok {
			maxPowerInfo[d.PowerInfo.CircuitGroupId] = maxval + maxTotalPowerConsumed
		} else {
			maxPowerInfo[d.PowerInfo.CircuitGroupId] = maxTotalPowerConsumed
		}

		//also cache stations so other metrics can figure out a circuit id from a station name
		(*c.TrackedStations)[d.Name] = d
	}
	for circuitId, powerConsumed := range powerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		TrainStationPower.WithLabelValues(cid, frmAddress, saveName).Set(powerConsumed)
	}
	for circuitId, powerConsumed := range maxPowerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		TrainStationPowerMax.WithLabelValues(cid, frmAddress, saveName).Set(powerConsumed)
	}
}
