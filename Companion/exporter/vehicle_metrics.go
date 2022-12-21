package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	VehicleRoundTrip = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "vehicle_round_trip_seconds",
		Help: "Recorded vehicle round trip time in seconds",
	}, []string{
		"id",
		"vehicle_type",
		"path_name",
	})
	VehicleFuel = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "vehicle_fuel",
		Help: "Amount of fuel remaining",
	}, []string{
		"id",
		"vehicle_type",
		"fuel_type",
	})
)
