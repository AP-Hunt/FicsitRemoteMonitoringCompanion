package exporter

import (
	"log"
	"strconv"
)

var ResourceSinkPowerConsumption = 30.0

type ResourceSinkCollector struct {
	endpoint string
}

type ResourceSinkDetails struct {
	Location  Location  `json:"location"`
	PowerInfo PowerInfo `json:"PowerInfo"`
}

func NewResourceSinkCollector(endpoint string) *ResourceSinkCollector {
	return &ResourceSinkCollector {
		endpoint: endpoint,
	}
}

func (c *ResourceSinkCollector) Collect(frmAddress string, sessionName string) {
	details := []ResourceSinkDetails{}
	err := retrieveData(frmAddress+c.endpoint, &details)
	if err != nil {
		log.Printf("error reading vehicle station statistics from FRM: %s\n", err)
		return
	}

	powerInfo := map[float64]float64{}
	maxPowerInfo := map[float64]float64{}
	for _, d := range details {
		val, ok := powerInfo[d.PowerInfo.CircuitGroupId]
		if ok {
			powerInfo[d.PowerInfo.CircuitGroupId] = val + d.PowerInfo.PowerConsumed
		} else {
			powerInfo[d.PowerInfo.CircuitGroupId] = d.PowerInfo.PowerConsumed
		}
		val, ok = maxPowerInfo[d.PowerInfo.CircuitGroupId]
		if ok {
			maxPowerInfo[d.PowerInfo.CircuitGroupId] = val + ResourceSinkPowerConsumption
		} else {
			maxPowerInfo[d.PowerInfo.CircuitGroupId] = ResourceSinkPowerConsumption
		}
	}
	for circuitId, powerConsumed := range powerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		ResourceSinkPower.WithLabelValues(cid, frmAddress, sessionName).Set(powerConsumed)
	}
	for circuitId, powerConsumed := range maxPowerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		ResourceSinkPowerMax.WithLabelValues(cid, frmAddress, sessionName).Set(powerConsumed)
	}
}

func (c *ResourceSinkCollector) DropCache() {}
