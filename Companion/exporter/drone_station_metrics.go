package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	DronePortFuelRate = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "drone_port_fuel_rate",
		Help: "Rate of fuel used",
	}, []string{
		"id",
		"home_station",
		"fuel_name",
	})
	DronePortFuelAmount = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "drone_port_fuel_amount",
		Help: "Amount of fuel in inventory",
	}, []string{
		"id",
		"home_station",
		"fuel_name",
	})
	DronePortRndTrip = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "drone_port_round_trip_seconds",
		Help: "Recorded drone round trip time in seconds",
	}, []string{
		"id",
		"home_station",
		"paired_station",
	})
	DronePortPower = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "drone_port_power",
		Help: "Drone port power in MW",
	}, []string{
		"circuit_id",
	})
)
