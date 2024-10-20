package exporter

import (
	"log"
	"strconv"
)

type PowerInfo struct {
	CircuitId     float64 `json:"CircuitID"`
	PowerConsumed float64 `json:"PowerConsumed"`
}

type PowerCollector struct {
	endpoint string
}

type PowerDetails struct {
	CircuitId           float64 `json:"CircuitID"`
	PowerConsumed       float64 `json:"PowerConsumed"`
	PowerCapacity       float64 `json:"PowerCapacity"`
	PowerMaxConsumed    float64 `json:"PowerMaxConsumed"`
	BatteryDifferential float64 `json:"BatteryDifferential"`
	BatteryPercent      float64 `json:"BatteryPercent"`
	BatteryCapacity     float64 `json:"BatteryCapacity"`
	BatteryTimeEmpty    string  `json:"BatteryTimeEmpty"`
	BatteryTimeFull     string  `json:"BatteryTimeFull"`
	FuseTriggered       bool    `json:"FuseTriggered"`
}

func NewPowerCollector(endpoint string) *PowerCollector {
	return &PowerCollector{
		endpoint: endpoint,
	}
}

func (c *PowerCollector) Collect(frmAddress string, saveName string) {
	details := []PowerDetails{}
	err := retrieveData(frmAddress+c.endpoint, &details)
	if err != nil {
		log.Printf("error reading power statistics from FRM: %s\n", err)
		return
	}

	for _, d := range details {
		circuitId := strconv.FormatFloat(d.CircuitId, 'f', -1, 64)
		PowerConsumed.WithLabelValues(circuitId, frmAddress, saveName).Set(d.PowerConsumed)
		PowerCapacity.WithLabelValues(circuitId, frmAddress, saveName).Set(d.PowerCapacity)
		PowerMaxConsumed.WithLabelValues(circuitId, frmAddress, saveName).Set(d.PowerMaxConsumed)
		BatteryDifferential.WithLabelValues(circuitId, frmAddress, saveName).Set(d.BatteryDifferential)
		BatteryPercent.WithLabelValues(circuitId, frmAddress, saveName).Set(d.BatteryPercent)
		BatteryCapacity.WithLabelValues(circuitId, frmAddress, saveName).Set(d.BatteryCapacity)
		batterySecondsRemaining := parseTimeSeconds(d.BatteryTimeEmpty)
		if batterySecondsRemaining != nil {
			BatterySecondsEmpty.WithLabelValues(circuitId, frmAddress, saveName).Set(*batterySecondsRemaining)
		}
		batterySecondsFull := parseTimeSeconds(d.BatteryTimeFull)
		if batterySecondsFull != nil {
			BatterySecondsFull.WithLabelValues(circuitId, frmAddress, saveName).Set(*batterySecondsFull)
		}
		fuseTriggered := parseBool(d.FuseTriggered)
		FuseTriggered.WithLabelValues(circuitId, frmAddress, saveName).Set(fuseTriggered)
	}
}
