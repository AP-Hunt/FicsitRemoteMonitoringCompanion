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
)

func ReportAllMetrics() []*prometheus.MetricVec {
	return []*prometheus.MetricVec{
		ItemProductionCapacityPerMinute.MetricVec,
		ItemProductionCapacityPercent.MetricVec,
		ItemConsumptionCapacityPerMinute.MetricVec,
		ItemConsumptionCapacityPercent.MetricVec,
		ItemsProducedPerMin.MetricVec,
		ItemsConsumedPerMin.MetricVec,
		MachineItemsProducedPerMin.MetricVec,
		MachineItemsProducedEffiency.MetricVec,
	}
}
