package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	MachineItemsProducedPerMin = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "machine_items_produced_per_min",
		Help: "How much of an item a building is producing",
	}, []string{
		"item_name",
		"machine_name",
		"x",
		"y",
		"z",
	})

	MachineItemsProducedEffiency = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "machine_items_produced_pc",
		Help: "The efficiency with which a building is producing an item",
	}, []string{
		"item_name",
		"machine_name",
		"x",
		"y",
		"z",
	})
	FactoryPower = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "factory_power",
		Help: "Power draw from factory machines in MW. Does not include extractors.",
	}, []string{
		"circuit_id",
	})

	FactoryPowerMax = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "factory_power_max",
		Help: "Max power draw from factory machines in MW. Does not include extractors.",
	}, []string{
		"circuit_id",
	})
)
