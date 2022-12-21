package exporter

import (
	"log"
	"strconv"
)

type PowerCollector struct {
	FRMAddress string
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

func NewPowerCollector(frmAddress string) *PowerCollector {
	return &PowerCollector{
		FRMAddress: frmAddress,
	}
}

func (c *PowerCollector) Collect() {
	details := []PowerDetails{}
	err := retrieveData(c.FRMAddress, &details)
	if err != nil {
		log.Printf("error reading power statistics from FRM: %s\n", err)
		return
	}

	for _, d := range details {
		circuitId := strconv.FormatFloat(d.CircuitId, 'f', -1, 64)
		PowerConsumed.WithLabelValues(circuitId).Set(d.PowerConsumed)
		PowerCapacity.WithLabelValues(circuitId).Set(d.PowerCapacity)
		PowerMaxConsumed.WithLabelValues(circuitId).Set(d.PowerMaxConsumed)
		BatteryDifferential.WithLabelValues(circuitId).Set(d.BatteryDifferential)
		BatteryPercent.WithLabelValues(circuitId).Set(d.BatteryPercent)
		BatteryCapacity.WithLabelValues(circuitId).Set(d.BatteryCapacity)
		batterySecondsRemaining := parseTimeSeconds(d.BatteryTimeEmpty)
		if batterySecondsRemaining != nil {
			BatterySecondsEmpty.WithLabelValues(circuitId).Set(*batterySecondsRemaining)
		}
		batterySecondsFull := parseTimeSeconds(d.BatteryTimeFull)
		if batterySecondsFull != nil {
			BatterySecondsFull.WithLabelValues(circuitId).Set(*batterySecondsFull)
		}
		fuseTriggered := parseBool(d.FuseTriggered)
		FuseTriggered.WithLabelValues(circuitId).Set(fuseTriggered)
	}
}
