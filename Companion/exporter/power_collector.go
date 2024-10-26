package exporter

import (
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type PowerInfo struct {
	CircuitGroupId float64 `json:"CircuitGroupID"`
	PowerConsumed  float64 `json:"PowerConsumed"`
}

type PowerCollector struct {
	endpoint       string
	metricsDropper *MetricsDropper
}

type PowerDetails struct {
	CircuitGroupId      float64 `json:"CircuitGroupID"`
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
		metricsDropper: NewMetricsDropper(
			PowerConsumed,
			PowerCapacity,
			PowerMaxConsumed,
			BatteryDifferential,
			BatteryPercent,
			BatteryCapacity,
			BatterySecondsEmpty,
			BatterySecondsFull,
			FuseTriggered,
			TrainCircuitPower,
			TrainCircuitPowerMax,
			TrainStationPower,
			TrainStationPowerMax,
			VehicleStationPower,
			VehicleStationPowerMax,
			FactoryPower,
			FactoryPowerMax,
		),
	}
}

func (c *PowerCollector) Collect(frmAddress string, sessionName string) {
	details := []PowerDetails{}
	err := retrieveData(frmAddress+c.endpoint, &details)
	if err != nil {
		c.metricsDropper.DropStaleMetricLabels()
		log.Printf("error reading power statistics from FRM: %s\n", err)
		return
	}

	for _, d := range details {
		circuitId := strconv.FormatFloat(d.CircuitGroupId, 'f', -1, 64)
		c.metricsDropper.CacheFreshMetricLabel(prometheus.Labels{"url": frmAddress, "session_name": sessionName, "circuit_id": circuitId})
		PowerConsumed.WithLabelValues(circuitId, frmAddress, sessionName).Set(d.PowerConsumed)
		PowerCapacity.WithLabelValues(circuitId, frmAddress, sessionName).Set(d.PowerCapacity)
		PowerMaxConsumed.WithLabelValues(circuitId, frmAddress, sessionName).Set(d.PowerMaxConsumed)
		BatteryDifferential.WithLabelValues(circuitId, frmAddress, sessionName).Set(d.BatteryDifferential)
		BatteryPercent.WithLabelValues(circuitId, frmAddress, sessionName).Set(d.BatteryPercent)
		BatteryCapacity.WithLabelValues(circuitId, frmAddress, sessionName).Set(d.BatteryCapacity)
		batterySecondsRemaining := parseTimeSeconds(d.BatteryTimeEmpty)
		if batterySecondsRemaining != nil {
			BatterySecondsEmpty.WithLabelValues(circuitId, frmAddress, sessionName).Set(*batterySecondsRemaining)
		}
		batterySecondsFull := parseTimeSeconds(d.BatteryTimeFull)
		if batterySecondsFull != nil {
			BatterySecondsFull.WithLabelValues(circuitId, frmAddress, sessionName).Set(*batterySecondsFull)
		}
		fuseTriggered := parseBool(d.FuseTriggered)
		FuseTriggered.WithLabelValues(circuitId, frmAddress, sessionName).Set(fuseTriggered)
	}
	c.metricsDropper.DropStaleMetricLabels()
}
