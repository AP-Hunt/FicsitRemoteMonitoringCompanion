package exporter

import (
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type DroneStationCollector struct {
	endpoint       string
	metricsDropper *MetricsDropper
}

type DroneFuelInventory struct {
	Name   string  `json:"Name"`
	Amount float64 `json:"Amount"`
}

type DroneActiveFuel struct {
	Name string  `json:"FuelName"`
	Rate float64 `json:"EstimatedFuelCostRate"`
}

type DroneStationDetails struct {
	Id                     string               `json:"ID"`
	HomeStation            string               `json:"Name"`
	PairedStation          string               `json:"PairedStation"`
	DroneStatus            string               `json:"DroneStatus"`
	AvgIncRate             float64              `json:"AvgIncRate"`
	AvgIncStack            float64              `json:"AvgIncStack"`
	AvgOutRate             float64              `json:"AvgOutRate"`
	AvgOutStack            float64              `json:"AvgOutStack"`
	AvgRndTrip             string               `json:"AvgRndTrip"`
	AvgTotalIncRate        float64              `json:"AvgTotalIncRate"`
	AvgTotalIncStack       float64              `json:"AvgTotalIncStack"`
	AvgTotalOutRate        float64              `json:"AvgTotalOutRate"`
	AvgTotalOutStack       float64              `json:"AvgTotalOutStack"`
	AvgTripIncAmt          float64              `json:"AvgTripIncAmt"`
	EstTotalTransRate      float64              `json:"EstTotalTransRate"`
	EstTransRate           float64              `json:"EstTransRate"`
	EstLatestTotalIncStack float64              `json:"EstLatestTotalIncStack"`
	EstLatestTotalOutStack float64              `json:"EstLatestTotalOutStack"`
	LatestIncStack         float64              `json:"LatestIncStack"`
	LatestOutStack         float64              `json:"LatestOutStack"`
	LatestRndTrip          float64              `json:"LatestRndTrip"`
	LatestTripIncAmt       float64              `json:"LatestTripIncAmt"`
	LatestTripOutAmt       float64              `json:"LatestTripOutAmt"`
	MedianRndTrip          string               `json:"MedianRndTrip"`
	MedianTripIncAmt       float64              `json:"MedianTripIncAmt"`
	MedianTripOutAmt       float64              `json:"MedianTripOutAmt"`
	PowerInfo              PowerInfo            `json:"PowerInfo"`
	Fuel                   []DroneFuelInventory `json:"FuelInventory"`
	ActiveFuel             DroneActiveFuel      `json:"ActiveFuel"`
}

func NewDroneStationCollector(endpoint string) *DroneStationCollector {
	return &DroneStationCollector{
		endpoint: endpoint,
		metricsDropper: NewMetricsDropper(
			DronePortFuelRate,
			DronePortFuelAmount,
			DronePortRndTrip,
		),
	}
}

func (c *DroneStationCollector) Collect(frmAddress string, sessionName string) {
	details := []DroneStationDetails{}
	err := retrieveData(frmAddress, c.endpoint, &details)
	if err != nil {
		c.metricsDropper.DropStaleMetricLabels()
		log.Printf("error reading drone station statistics from FRM: %s\n", err)
		return
	}

	powerInfo := map[float64]float64{}
	maxPowerInfo := map[float64]float64{}
	for _, d := range details {
		c.metricsDropper.CacheFreshMetricLabel(prometheus.Labels{"url": frmAddress, "session_name": sessionName, "id": d.Id})
		id := d.Id
		home := d.HomeStation
		paired := d.PairedStation

		if len(d.Fuel) > 0 {
			DronePortFuelAmount.WithLabelValues(id, home, d.Fuel[0].Name, frmAddress, sessionName).Set(d.Fuel[0].Amount)
			DronePortFuelRate.WithLabelValues(id, home, d.Fuel[0].Name, frmAddress, sessionName).Set(d.ActiveFuel.Rate)
			DronePortRndTrip.WithLabelValues(id, home, paired, frmAddress, sessionName).Set(d.LatestRndTrip)
		}

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
		DronePortPower.WithLabelValues(cid, frmAddress, sessionName).Set(powerConsumed)
	}
	for circuitId, powerConsumed := range maxPowerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		DronePortPowerMax.WithLabelValues(cid, frmAddress, sessionName).Set(powerConsumed)
	}

	c.metricsDropper.DropStaleMetricLabels()
}

func (c *DroneStationCollector) DropCache() {}
