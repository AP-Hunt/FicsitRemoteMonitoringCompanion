package exporter

import (
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type StorageInventoryCollector struct {
	endpoint       string
	metricsDropper *MetricsDropper
}

func NewStorageInventoryCollector(endpoint string) *StorageInventoryCollector {
	return &StorageInventoryCollector{
		endpoint: endpoint,
		metricsDropper: NewMetricsDropper(
			StorageInventory,
			StorageInventoryMax,
		),
	}
}

func (c *StorageInventoryCollector) Collect(frmAddress string, sessionName string) {
	details := []ContainerDetail{}
	err := retrieveData(frmAddress, c.endpoint, &details)
	if err != nil {
		c.metricsDropper.DropStaleMetricLabels()
		log.Printf("error reading inventory statistics from FRM: %s\n", err)
		return
	}

	for _, detail := range details {
		c.metricsDropper.CacheFreshMetricLabel(prometheus.Labels{
			"url":          frmAddress,
			"session_name": sessionName,
			"id":           detail.Id,
		})
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

			StorageInventoryMax.WithLabelValues(
				item.Name,
				detail.Name,
				strconv.FormatFloat(detail.Location.X, 'f', -1, 64),
				strconv.FormatFloat(detail.Location.Y, 'f', -1, 64),
				strconv.FormatFloat(detail.Location.Z, 'f', -1, 64),
				frmAddress,
				sessionName,
			).Set(float64(item.MaxAmount))
		}
	}
	c.metricsDropper.DropStaleMetricLabels()
}

func (c *StorageInventoryCollector) DropCache() {}
