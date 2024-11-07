package exporter

import (
	"log"
	"strconv"
)

type TrainStationCollector struct {
	endpoint string
}

type CargoPlatform struct {
	LoadingDock   string    `json:"LoadingDock"`
	TransferRate  float64   `json:"TransferRate"`
	LoadingStatus string    `json:"LoadingStatus"` // Idle, Loading, Unloading
	LoadingMode   string    `json:"LoadingMode"`
	PowerInfo     PowerInfo `json:"PowerInfo"`
}

type TrainStationDetails struct {
	Name           string          `json:"Name"`
	Location       Location        `json:"location"`
	CargoPlatforms []CargoPlatform `json:"CargoInventory"`
	PowerInfo      PowerInfo       `json:"PowerInfo"`
}

func NewTrainStationCollector(endpoint string) *TrainStationCollector {
	return &TrainStationCollector{
		endpoint: endpoint,
	}
}

func (c *TrainStationCollector) Collect(frmAddress string, sessionName string) {
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

		// some additional calculations: power listed here is only for the station.
		// each of the cargo platforms have power stats returned. add up power metrics for total power use.
		totalPowerConsumed := d.PowerInfo.PowerConsumed
		maxTotalPowerConsumed := d.PowerInfo.MaxPowerConsumed
		for _, p := range d.CargoPlatforms {
			totalPowerConsumed = totalPowerConsumed + p.PowerInfo.PowerConsumed
			maxTotalPowerConsumed = maxTotalPowerConsumed + p.PowerInfo.MaxPowerConsumed
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
	}
	for circuitId, powerConsumed := range powerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		TrainStationPower.WithLabelValues(cid, frmAddress, sessionName).Set(powerConsumed)
	}
	for circuitId, powerConsumed := range maxPowerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		TrainStationPowerMax.WithLabelValues(cid, frmAddress, sessionName).Set(powerConsumed)
	}
}

func (c *TrainStationCollector) DropCache() {}
