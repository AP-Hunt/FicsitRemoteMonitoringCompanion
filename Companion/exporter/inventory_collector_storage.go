package exporter

import (
	"log"
	"strconv"
)

type StorageInventoryCollector struct {
	endpoint string
}

func NewStorageInventoryCollector(endpoint string) *StorageInventoryCollector {
	return &StorageInventoryCollector{
		endpoint: endpoint,
	}
}

func (c *StorageInventoryCollector) Collect(frmAddress string, sessionName string) {
	details := []ContainerDetail{}
	err := retrieveData(frmAddress+c.endpoint, &details)
	if err != nil {
		log.Printf("error reading inventory statistics from FRM: %s\n", err)
		return
	}

	for _, detail := range details {
		for _, item := range detail.Inventory {
			StorageInventory.WithLabelValues(
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

func (c *StorageInventoryCollector) DropCache() {}
