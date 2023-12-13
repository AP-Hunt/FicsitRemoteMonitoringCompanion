package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TrainStationPower = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_station_power",
		Help: "Train station power in MW",
	}, []string{
		"circuit_id",
	})
)
