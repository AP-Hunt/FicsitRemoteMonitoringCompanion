package exporter

import (
	"log"
	"strconv"
)

var VehicleStationPowerConsumption = 20.0

type VehicleStationCollector struct {
	endpoint string
}

type VehicleStationDetails struct {
	Name      string    `json:"Name"`
	Location  Location  `json:"location"`
	PowerInfo PowerInfo `json:"PowerInfo"`
}

func NewVehicleStationCollector(endpoint string) *VehicleStationCollector {
	return &VehicleStationCollector{
		endpoint: endpoint,
	}
}

func (c *VehicleStationCollector) Collect(frmAddress string, saveName string) {
	details := []VehicleStationDetails{}
	err := retrieveData(frmAddress+c.endpoint, &details)
	if err != nil {
		log.Printf("error reading vehicle station statistics from FRM: %s\n", err)
		return
	}

	powerInfo := map[float64]float64{}
	maxPowerInfo := map[float64]float64{}
	for _, d := range details {
		val, ok := powerInfo[d.PowerInfo.CircuitId]
		if ok {
			powerInfo[d.PowerInfo.CircuitId] = val + d.PowerInfo.PowerConsumed
		} else {
			powerInfo[d.PowerInfo.CircuitId] = d.PowerInfo.PowerConsumed
		}
		val, ok = maxPowerInfo[d.PowerInfo.CircuitId]
		if ok {
			maxPowerInfo[d.PowerInfo.CircuitId] = val + VehicleStationPowerConsumption
		} else {
			maxPowerInfo[d.PowerInfo.CircuitId] = VehicleStationPowerConsumption
		}
	}
	for circuitId, powerConsumed := range powerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		VehicleStationPower.WithLabelValues(cid, frmAddress, saveName).Set(powerConsumed)
	}
	for circuitId, powerConsumed := range maxPowerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		VehicleStationPowerMax.WithLabelValues(cid, frmAddress, saveName).Set(powerConsumed)
	}
}
