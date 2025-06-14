package exporter

import (
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	PortalPower = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "portal_power",
		Help: "portal power use in MW",
	}, []string{
		"circuit_id",
	})

	PortalPowerMax = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "portal_power_max",
		Help: "portal max power use in MW",
	}, []string{
		"circuit_id",
	})
)

type PortalCollector struct {
	endpoint string
}

type PortalDetails struct {
	Location  Location  `json:"location"`
	PowerInfo PowerInfo `json:"PowerInfo"`
}

func NewPortalCollector(endpoint string) *PortalCollector {
	return &PortalCollector{
		endpoint: endpoint,
	}
}

func (c *PortalCollector) Collect(frmAddress string, sessionName string) {
	details := []PortalDetails{}
	err := retrieveData(frmAddress+c.endpoint, &details)
	if err != nil {
		log.Printf("error reading portal statistics from FRM: %s\n", err)
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
		// TODO: max portal power is bugged in the base game.
		// Replace with reported values when they are correct.
		if ok {
			maxPowerInfo[d.PowerInfo.CircuitGroupId] = val + MaxPortalPower
		} else {
			maxPowerInfo[d.PowerInfo.CircuitGroupId] = MaxPortalPower
		}
	}
	for circuitId, powerConsumed := range powerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		PortalPower.WithLabelValues(cid, frmAddress, sessionName).Set(powerConsumed)
	}
	for circuitId, powerConsumed := range maxPowerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		PortalPowerMax.WithLabelValues(cid, frmAddress, sessionName).Set(powerConsumed)
	}
}

func (c *PortalCollector) DropCache() {}
