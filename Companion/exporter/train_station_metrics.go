package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TrainStationPower = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_station_power",
		Help: "Train station power consumed in MW",
	}, []string{
		"circuit_id",
	})

	TrainStationPowerMax = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_station_power_max",
		Help: "Train station power max consumed in MW",
	}, []string{
		"circuit_id",
	})
)
