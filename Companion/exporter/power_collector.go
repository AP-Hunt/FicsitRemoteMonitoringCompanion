package exporter

import (
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

type PowerInfo struct {
	CircuitGroupId   float64 `json:"CircuitGroupID"`
	PowerConsumed    float64 `json:"PowerConsumed"`
	MaxPowerConsumed float64 `json:"MaxPowerConsumed"`
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

// power max calculated via aggregate.
// actual max is either reported max, or calculated max. Take the higher one.
// see power_info.go for details about the bug we're working around.
func calculateMaxPowerCategory(gauge *prometheus.GaugeVec, circuitId string, frmAddress string, sessionName string) float64 {
	var m = &dto.Metric{}
	err := gauge.WithLabelValues(circuitId, frmAddress, sessionName).Write(m)
	if err != nil {
		return 0
	}
	return m.Gauge.GetValue()
}
func calculateMaxPower(reportedMaxPower float64, circuitId string, frmAddress string, sessionName string) float64 {
	maxConsumed := reportedMaxPower
	factoryPowerMax := calculateMaxPowerCategory(FactoryPowerMax, circuitId, frmAddress, sessionName)
	extractorPowerMax := calculateMaxPowerCategory(ExtractorPowerMax, circuitId, frmAddress, sessionName)
	dronePowerMax := calculateMaxPowerCategory(DronePortPowerMax, circuitId, frmAddress, sessionName)
	frackingPowerMax := calculateMaxPowerCategory(FrackingPowerMax, circuitId, frmAddress, sessionName)
	hypertubePowerMax := calculateMaxPowerCategory(HypertubePowerMax, circuitId, frmAddress, sessionName)
	portalPowerMax := calculateMaxPowerCategory(PortalPowerMax, circuitId, frmAddress, sessionName)
	pumpPowerMax := calculateMaxPowerCategory(PumpPowerMax, circuitId, frmAddress, sessionName)
	resourceSinkPowerMax := calculateMaxPowerCategory(ResourceSinkPowerMax, circuitId, frmAddress, sessionName)
	trainPowerMax := calculateMaxPowerCategory(TrainCircuitPowerMax, circuitId, frmAddress, sessionName)
	trainStationPowerMax := calculateMaxPowerCategory(TrainStationPowerMax, circuitId, frmAddress, sessionName)
	vehicleStationPowerMax := calculateMaxPowerCategory(VehicleStationPowerMax, circuitId, frmAddress, sessionName)
	calculatedMaxConsumed := factoryPowerMax +
		extractorPowerMax +
		dronePowerMax +
		frackingPowerMax +
		hypertubePowerMax +
		portalPowerMax +
		pumpPowerMax +
		resourceSinkPowerMax +
		trainPowerMax +
		trainStationPowerMax +
		vehicleStationPowerMax
	if calculatedMaxConsumed > maxConsumed {
		maxConsumed = calculatedMaxConsumed
	}
	return maxConsumed
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
			ResourceSinkPower,
			ResourceSinkPowerMax,
			DronePortPower,
			DronePortPowerMax,
			PumpPower,
			PumpPowerMax,
			ExtractorPower,
			ExtractorPowerMax,
			HypertubePower,
			HypertubePowerMax,
			PortalPower,
			PortalPowerMax,
			FrackingPower,
			FrackingPowerMax,
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

		maxConsumed := calculateMaxPower(d.PowerMaxConsumed, circuitId, frmAddress, sessionName)
		PowerMaxConsumed.WithLabelValues(circuitId, frmAddress, sessionName).Set(maxConsumed)

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

func (c *PowerCollector) DropCache() {}
