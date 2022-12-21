package exporter

import (
	"log"
	"strconv"
)

type FactoryBuildingCollector struct {
	FRMAddress string
}

func NewFactoryBuildingCollector(frmAddress string) *FactoryBuildingCollector {
	return &FactoryBuildingCollector{
		FRMAddress: frmAddress,
	}
}

func (c *FactoryBuildingCollector) Collect() {
	details := []BuildingDetail{}
	err := retrieveData(c.FRMAddress, &details)
	if err != nil {
		log.Printf("error reading factory buildings from FRM: %s\n", err)
		return
	}

	for _, building := range details {
		for _, prod := range building.Production {
			MachineItemsProducedPerMin.WithLabelValues(
				prod.Name,
				building.Building,
				strconv.FormatFloat(building.Location.X, 'f', -1, 64),
				strconv.FormatFloat(building.Location.Y, 'f', -1, 64),
				strconv.FormatFloat(building.Location.Z, 'f', -1, 64),
			).Set(prod.CurrentProd)

			MachineItemsProducedEffiency.WithLabelValues(
				prod.Name,
				building.Building,
				strconv.FormatFloat(building.Location.X, 'f', -1, 64),
				strconv.FormatFloat(building.Location.Y, 'f', -1, 64),
				strconv.FormatFloat(building.Location.Z, 'f', -1, 64),
			).Set(prod.ProdPercent)
		}
	}
}
