package exporter

import (
	"log"
)

type ProductionCollector struct {
	FRMAddress string
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

func NewProductionCollector(frmAddress string) *ProductionCollector {
	return &ProductionCollector{
		FRMAddress: frmAddress,
	}
}

func (c *ProductionCollector) Collect() {
	details := []ProductionDetails{}
	err := retrieveData(c.FRMAddress, &details)
	if err != nil {
		log.Printf("error reading production statistics from FRM: %s\n", err)
		return
	}

	for _, d := range details {
		ItemsProducedPerMin.WithLabelValues(d.ItemName).Set(d.CurrentProduction)
		ItemsConsumedPerMin.WithLabelValues(d.ItemName).Set(d.CurrentConsumption)

		ItemProductionCapacityPercent.WithLabelValues(d.ItemName).Set(d.ProdPercent)
		ItemConsumptionCapacityPercent.WithLabelValues(d.ItemName).Set(d.ConsPercent)
		ItemProductionCapacityPerMinute.WithLabelValues(d.ItemName).Set(d.MaxProd)
		ItemConsumptionCapacityPerMinute.WithLabelValues(d.ItemName).Set(d.MaxConsumed)
	}
}
