package exporter

import (
	"log"
	"strconv"
)

type CrateInventoryCollector struct {
	endpoint string
}

func NewCrateInventoryCollector(endpoint string) *CrateInventoryCollector {
	return &CrateInventoryCollector{
		endpoint: endpoint,
	}
}

func (c *CrateInventoryCollector) Collect(frmAddress string, sessionName string) {
	details := []ContainerDetail{}
	err := retrieveData(frmAddress+c.endpoint, &details)
	if err != nil {
		log.Printf("error reading inventory statistics from FRM: %s\n", err)
		return
	}

	for _, detail := range details {
		for _, item := range detail.Inventory {
			CrateInventory.WithLabelValues(
				item.Name,
				detail.Name,
				strconv.FormatFloat(detail.Location.X, 'f', -1, 64),
				strconv.FormatFloat(detail.Location.Y, 'f', -1, 64),
				strconv.FormatFloat(detail.Location.Z, 'f', -1, 64),
				frmAddress,
				sessionName,
			).Set(float64(item.Amount))
		}
	}
}

func (c *CrateInventoryCollector) DropCache() {}
