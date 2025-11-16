package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Cloud Inventory (Dimensional Depot)
	CloudInventory = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "cloud_inventory",
		Help: "Items stored in the dimensional depot",
	}, []string{
		"item_name",
	})

	CloudInventoryMax = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "cloud_inventory_max",
		Help: "Stack size for items in the dimensional depot",
	}, []string{
		"item_name",
	})

	// World Inventory
	WorldInventory = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "world_inventory",
		Help: "Inventory of the world regardless of location (All buildings whom purpose is to provide storage)",
	}, []string{
		"item_name",
	})

	WorldInventoryMax = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "world_inventory_max",
		Help: "Stack size for items in the world invetory",
	}, []string{
		"item_name",
	})

	// Storage Container Inventory
	StorageInventory = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "storage_inventory",
		Help: "Items stored inside storage containers",
	}, []string{
		"item_name",
		"container_name",
		"x",
		"y",
		"z",
	})

	StorageInventoryMax = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "storage_inventory_max",
		Help: "Stack size for items stored in storage containers",
	}, []string{
		"item_name",
		"container_name",
		"x",
		"y",
		"z",
	})

	// Crate Inventory (Dismantle and Death Crates)
	CrateInventory = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "crate_inventory",
		Help: "Items stored inside crates",
	}, []string{
		"item_name",
		"container_name",
		"x",
		"y",
		"z",
	})

	CrateInventoryMax = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "crate_inventory_max",
		Help: "Stack size for items stored in crates",
	}, []string{
		"item_name",
		"container_name",
		"x",
		"y",
		"z",
	})
)
