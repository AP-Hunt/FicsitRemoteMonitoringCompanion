package exporter

import (
	"log"
)

type WorldInventoryCollector struct {
	endpoint string
}

func NewWorldInventoryCollector(endpoint string) *WorldInventoryCollector {
	return &WorldInventoryCollector{
		endpoint: endpoint,
	}
}

func (c *WorldInventoryCollector) Collect(frmAddress string, sessionName string) {
	items := []InventoryItem{}
	err := retrieveData(frmAddress, c.endpoint, &items)
	if err != nil {
		log.Printf("error reading inventory statistics from FRM: %s\n", err)
		return
	}

	for _, item := range items {
		WorldInventory.WithLabelValues(item.Name, frmAddress, sessionName).Set(float64(item.Amount))
		WorldInventoryMax.WithLabelValues(item.Name, frmAddress, sessionName).Set(float64(item.MaxAmount))
	}
}

func (c *WorldInventoryCollector) DropCache() {}
