package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	VehicleStationPower = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "vehicle_station_power",
		Help: "Vehicle station power in MW",
	}, []string{
		"circuit_id",
	})
)
