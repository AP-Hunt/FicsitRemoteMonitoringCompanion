package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ItemProductionCapacityPerMinute = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "item_production_capacity_per_min",
		Help: "The factory's capacity for the production of an item, per minute",
	}, []string{
		"item_name",
	})

	ItemProductionCapacityPercent = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "item_production_capacity_pc",
		Help: "The percentage of an item's production capacity being used",
	}, []string{
		"item_name",
	})

	ItemConsumptionCapacityPerMinute = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "item_consumption_capacity_per_min",
		Help: "The factory's capacity for the consumption of an item, per minute",
	}, []string{
		"item_name",
	})

	ItemConsumptionCapacityPercent = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "item_consumption_capacity_pc",
		Help: "The percentage of an item's consumption capacity being used",
	}, []string{
		"item_name",
	})

	ItemsProducedPerMin = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "items_produced_per_min",
		Help: "The number of an item being produced, per minute",
	}, []string{
		"item_name",
	})

	ItemsConsumedPerMin = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "items_consumed_per_min",
		Help: "The number of an item being consumed, per minute",
	}, []string{
		"item_name",
	})

	PowerConsumed = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "power_consumed",
		Help: "Power consumed on selected power circuit",
	}, []string{
		"circuit_id",
	})

	PowerCapacity = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "power_capacity",
		Help: "Power capacity on selected power circuit",
	}, []string{
		"circuit_id",
	})

	PowerMaxConsumed = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "power_max_consumed",
		Help: "Maximum Power that can be consumed on selected power circuit",
	}, []string{
		"circuit_id",
	})

	BatteryDifferential = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "battery_differential",
		Help: "Amount of power in excess/deficit going into or out of the battery bank(s). Positive = Charges batteries, Negative = Drains batteries",
	}, []string{
		"circuit_id",
	})

	BatteryPercent = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "battery_percent",
		Help: "Percentage of battery bank(s) charge",
	}, []string{
		"circuit_id",
	})

	BatteryCapacity = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "battery_capacity",
		Help: "Total capacity of battery bank(s)",
	}, []string{
		"circuit_id",
	})

	BatterySecondsEmpty = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "battery_seconds_empty",
		Help: "Seconds until Batteries are empty",
	}, []string{
		"circuit_id",
	})

	BatterySecondsFull = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "battery_seconds_full",
		Help: "Seconds until Batteries are full",
	}, []string{
		"circuit_id",
	})

	FuseTriggered = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "fuse_triggered",
		Help: "Has the fuse been triggered",
	}, []string{
		"circuit_id",
	})
)
