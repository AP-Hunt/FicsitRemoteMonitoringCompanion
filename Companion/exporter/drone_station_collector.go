package exporter

import (
	"log"
)

type DroneStationCollector struct {
	FRMAddress string
}

type DroneStationDetails struct {
	Id                     string  `json:"ID"`
	HomeStation            string  `json:"HomeStation"`
	PairedStation          string  `json:"PairedStation"`
	DroneStatus            string  `json:"DroneStatus"`
	AvgIncRate             float64 `json:"AvgIncRate"`
	AvgIncStack            float64 `json:"AvgIncStack"`
	AvgOutRate             float64 `json:"AvgOutRate"`
	AvgOutStack            float64 `json"AvgOutStack"`
	AvgRndTrip             string  `json:"AvgRndTrip"`
	AvgTotalIncRate        float64 `json:"AvgTotalIncRate"`
	AvgTotalIncStack       float64 `json:"AvgTotalIncStack"`
	AvgTotalOutRate        float64 `json:"AvgTotalOutRate"`
	AvgTotalOutStack       float64 `json:"AvgTotalOutStack"`
	AvgTripIncAmt          float64 `json:"AvgTripIncAmt"`
	EstRndTrip             string  `json:"EstRndTrip"`
	EstTotalTransRate      float64 `json:"EstTotalTransRate"`
	EstTransRate           float64 `json:"EstTransRate"`
	EstLatestTotalIncStack float64 `json:"EstLatestTotalIncStack"`
	EstLatestTotalOutStack float64 `json:"EstLatestTotalOutStack"`
	LatestIncStack         float64 `json:"LatestIncStack"`
	LatestOutStack         float64 `json:"LatestOutStack"`
	LatestRndTrip          string  `json:"LatestRndTrip"`
	LatestTripIncAmt       float64 `json:"LatestTripIncAmt"`
	LatestTripOutAmt       float64 `json:"LatestTripOutAmt"`
	MedianRndTrip          string  `json:"MedianRndTrip"`
	MedianTripIncAmt       float64 `json:"MedianTripIncAmt"`
	MedianTripOutAmt       float64 `json:"MedianTripOutAmt"`
	EstBatteryRate         float64 `json:"EstBatteryRate"`
}

func NewDroneStationCollector(frmAddress string) *DroneStationCollector {
	return &DroneStationCollector{
		FRMAddress: frmAddress,
	}
}

func (c *DroneStationCollector) Collect() {
	details := []DroneStationDetails{}
	err := retrieveData(c.FRMAddress, &details)
	if err != nil {
		log.Printf("error reading drone station statistics from FRM: %s\n", err)
		return
	}
	for _, d := range details {
		id := d.Id
		home := d.HomeStation
		paired := d.PairedStation

		DronePortBatteryRate.WithLabelValues(id, home, paired).Set(d.EstBatteryRate)

		roundTrip := parseTimeSeconds(d.LatestRndTrip)
		if roundTrip != nil {
			DronePortRndTrip.WithLabelValues(id, home, paired).Set(*roundTrip)
		}
	}
}
