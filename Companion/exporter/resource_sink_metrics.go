package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ResourceSinkPower = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "resource_sink_power",
		Help: "AWESOME sink power use in MW",
	}, []string{
		"circuit_id",
	})

	ResourceSinkPowerMax = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "resource_sink_power_max",
		Help: "AWESOME sink max power use in MW",
	}, []string{
		"circuit_id",
	})
)
