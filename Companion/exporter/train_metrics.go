package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TrainRoundTrip = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_round_trip_seconds",
		Help: "Recorded train round trip time in seconds",
	}, []string{
		"name",
	})
	TrainSegmentTrip = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_segment_trip_seconds",
		Help: "Recorded train trip between two stations",
	}, []string{
		"name",
		"from",
		"to",
	})
	TrainDerailed = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_derailed",
		Help: "Is train derailed",
	}, []string{
		"name",
	})
	TrainPower = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_power_consumed",
		Help: "How much power train is consuming",
	}, []string{
		"name",
	})
)
