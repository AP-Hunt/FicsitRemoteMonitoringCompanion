package exporter

import (
	"log"
	"strconv"
)

type TrainStationCollector struct {
	FRMAddress string
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

func NewTrainStationCollector(frmAddress string) *TrainStationCollector {
	return &TrainStationCollector{
		FRMAddress: frmAddress,
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
	for _, d := range details {
		val, ok := powerInfo[d.PowerInfo.CircuitId]

		// some additional calculations: for now, power listed here is only for the station.
		// add each of the cargo platforms' power info: 0.1MW if Idle, 50MW otherwise
		totalPowerConsumed := d.PowerInfo.PowerConsumed
		for _, p := range d.CargoPlatforms {
			if p.LoadingStatus == "Idle" {
				totalPowerConsumed = totalPowerConsumed + 0.1
			} else {
				totalPowerConsumed = totalPowerConsumed + 50
			}
		}

		if ok {
			powerInfo[d.PowerInfo.CircuitId] = val + totalPowerConsumed
		} else {
			powerInfo[d.PowerInfo.CircuitId] = totalPowerConsumed
		}
	}
	for circuitId, powerConsumed := range powerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		TrainStationPower.WithLabelValues(cid).Set(powerConsumed)
	}
}
