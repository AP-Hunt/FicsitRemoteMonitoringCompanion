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

	ResourceSinkTotalPoints = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "resource_sink_total_points",
		Help: "AWESOME sink total points",
	}, []string{
		"sink_type",
	})

	ResourceSinkPointsToCoupon = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "resource_sink_points_to_coupon",
		Help: "AWESOME sink points to next coupon",
	}, []string{
		"sink_type",
	})

	ResourceSinkPercent = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "resource_sink_percent",
		Help: "AWESOME sink percent to next coupon",
	}, []string{
		"sink_type",
	})

	ResourceSinkCollectedCoupons = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "resource_sink_collected_coupons",
		Help: "AWESOME sink collected coupons",
	}, []string{})
)
