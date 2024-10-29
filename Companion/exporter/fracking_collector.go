package exporter

import (
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	FrackingPower = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "fracking_power",
		Help: "fracking power use in MW",
	}, []string{
		"circuit_id",
	})

	FrackingPowerMax = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "fracking_power_max",
		Help: "fracking max power use in MW",
	}, []string{
		"circuit_id",
	})
)

type FrackingCollector struct {
	endpoint string
}

type FrackingDetails struct {
	Location  Location  `json:"location"`
	PowerInfo PowerInfo `json:"PowerInfo"`
}

func NewFrackingCollector(endpoint string) *FrackingCollector {
	return &FrackingCollector{
		endpoint: endpoint,
	}
}

func (c *FrackingCollector) Collect(frmAddress string, sessionName string) {
	details := []FrackingDetails{}
	err := retrieveData(frmAddress+c.endpoint, &details)
	if err != nil {
		log.Printf("error reading fracking statistics from FRM: %s\n", err)
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
		FrackingPower.WithLabelValues(cid, frmAddress, sessionName).Set(powerConsumed)
	}
	for circuitId, powerConsumed := range maxPowerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		FrackingPowerMax.WithLabelValues(cid, frmAddress, sessionName).Set(powerConsumed)
	}
}

func (c *FrackingCollector) DropCache() {}
