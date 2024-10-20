package exporter

import (
	"log"
)

type ProductionCollector struct {
	endpoint string
}

type ProductionDetails struct {
	ItemName           string   `json:"Name"`
	ProdPercent        float64 `json:"ProdPercent"`
	ConsPercent        float64 `json:"ConsPercent"`
	CurrentProduction  float64  `json:"CurrentProd"`
	CurrentConsumption float64  `json:"CurrentConsumed"`
	MaxProd            float64  `json:"MaxProd"`
	MaxConsumed        float64  `json:"MaxConsumed"`
}

func NewProductionCollector(endpoint string) *ProductionCollector {
	return &ProductionCollector{
		endpoint: endpoint,
	}
}

func (c *ProductionCollector) Collect(frmAddress string, saveName string) {
	details := []ProductionDetails{}
	err := retrieveData(frmAddress + c.endpoint, &details)
	if err != nil {
		log.Printf("error reading production statistics from FRM: %s\n", err)
		return
	}

	for _, d := range details {
		ItemsProducedPerMin.WithLabelValues(d.ItemName, frmAddress, saveName).Set(d.CurrentProduction)
		ItemsConsumedPerMin.WithLabelValues(d.ItemName, frmAddress, saveName).Set(d.CurrentConsumption)

		ItemProductionCapacityPercent.WithLabelValues(d.ItemName, frmAddress, saveName).Set(d.ProdPercent)
		ItemConsumptionCapacityPercent.WithLabelValues(d.ItemName, frmAddress, saveName).Set(d.ConsPercent)
		ItemProductionCapacityPerMinute.WithLabelValues(d.ItemName, frmAddress, saveName).Set(d.MaxProd)
		ItemConsumptionCapacityPerMinute.WithLabelValues(d.ItemName, frmAddress, saveName).Set(d.MaxConsumed)
	}
}
