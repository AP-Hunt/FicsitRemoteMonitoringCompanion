package exporter

import (
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	ExtractorPower = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "extractor_power",
		Help: "extractor power use in MW",
	}, []string{
		"circuit_id",
	})

	ExtractorPowerMax = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "extractor_power_max",
		Help: "extractor max power use in MW",
	}, []string{
		"circuit_id",
	})
)

type ExtractorCollector struct {
	endpoint string
}

type ExtractorDetails struct {
	Location  Location  `json:"location"`
	PowerInfo PowerInfo `json:"PowerInfo"`
}

func NewExtractorCollector(endpoint string) *ExtractorCollector {
	return &ExtractorCollector{
		endpoint: endpoint,
	}
}

func (c *ExtractorCollector) Collect(frmAddress string, sessionName string) {
	details := []ExtractorDetails{}
	err := retrieveData(frmAddress+c.endpoint, &details)
	if err != nil {
		log.Printf("error reading extractor statistics from FRM: %s\n", err)
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
		ExtractorPower.WithLabelValues(cid, frmAddress, sessionName).Set(powerConsumed)
	}
	for circuitId, powerConsumed := range maxPowerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		ExtractorPowerMax.WithLabelValues(cid, frmAddress, sessionName).Set(powerConsumed)
	}
}

func (c *ExtractorCollector) DropCache() {}
