package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	VehicleStationPower = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "vehicle_station_power",
		Help: "Vehicle station power use in MW",
	}, []string{
		"circuit_id",
	})

	VehicleStationPowerMax = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "vehicle_station_power_max",
		Help: "Vehicle station max power use in MW",
	}, []string{
		"circuit_id",
	})
)
