package exporter

import (
	"log"
)

type CloudInventoryCollector struct {
	endpoint string
}

func NewCloudInventoryCollector(endpoint string) *CloudInventoryCollector {
	return &CloudInventoryCollector{
		endpoint: endpoint,
	}
}

func (c *CloudInventoryCollector) Collect(frmAddress string, sessionName string) {
	items := []InventoryItem{}
	err := retrieveData(frmAddress, c.endpoint, &items)
	if err != nil {
		log.Printf("error reading inventory statistics from FRM: %s\n", err)
		return
	}

	for _, item := range items {
		CloudInventory.WithLabelValues(item.Name, frmAddress, sessionName).Set(float64(item.Amount))
		CloudInventoryMax.WithLabelValues(item.Name, frmAddress, sessionName).Set(float64(item.MaxAmount))
	}
}

func (c *CloudInventoryCollector) DropCache() {}
