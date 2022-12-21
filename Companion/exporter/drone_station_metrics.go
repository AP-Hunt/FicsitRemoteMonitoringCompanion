package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	DronePortBatteryRate = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "drone_port_battery_rate",
		Help: "Rate of batteries used",
	}, []string{
		"id",
		"home_station",
		"paired_station",
	})
	DronePortRndTrip = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "drone_port_round_trip_seconds",
		Help: "Recorded drone round trip time in seconds",
	}, []string{
		"id",
		"home_station",
		"paired_station",
	})
)
