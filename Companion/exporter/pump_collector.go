package exporter

import (
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	PumpPower = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "pump_power",
		Help: "pump power use in MW",
	}, []string{
		"circuit_id",
	})

	PumpPowerMax = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "pump_power_max",
		Help: "pump max power use in MW",
	}, []string{
		"circuit_id",
	})
)

type PumpCollector struct {
	endpoint string
}

type PumpDetails struct {
	Location  Location  `json:"location"`
	PowerInfo PowerInfo `json:"PowerInfo"`
}

func NewPumpCollector(endpoint string) *PumpCollector {
	return &PumpCollector{
		endpoint: endpoint,
	}
}

func (c *PumpCollector) Collect(frmAddress string, sessionName string) {
	details := []PumpDetails{}
	err := retrieveData(frmAddress+c.endpoint, &details)
	if err != nil {
		log.Printf("error reading pump statistics from FRM: %s\n", err)
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
			maxPowerInfo[d.PowerInfo.CircuitGroupId] = val + d.PowerInfo.MaxPowerConsumed
		} else {
			maxPowerInfo[d.PowerInfo.CircuitGroupId] = d.PowerInfo.MaxPowerConsumed
		}
	}
	for circuitId, powerConsumed := range powerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		PumpPower.WithLabelValues(cid, frmAddress, sessionName).Set(powerConsumed)
	}
	for circuitId, powerConsumed := range maxPowerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		PumpPowerMax.WithLabelValues(cid, frmAddress, sessionName).Set(powerConsumed)
	}
}

func (c *PumpCollector) DropCache() {}
