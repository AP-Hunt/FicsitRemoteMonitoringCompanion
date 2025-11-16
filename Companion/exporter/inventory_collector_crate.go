package exporter

import (
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type CrateInventoryCollector struct {
	endpoint       string
	metricsDropper *MetricsDropper
}

func NewCrateInventoryCollector(endpoint string) *CrateInventoryCollector {
	return &CrateInventoryCollector{
		endpoint: endpoint,
		metricsDropper: NewMetricsDropper(
			StorageInventory,
			StorageInventoryMax,
		),
	}
}

func (c *CrateInventoryCollector) Collect(frmAddress string, sessionName string) {
	details := []ContainerDetail{}
	err := retrieveData(frmAddress, c.endpoint, &details)
	if err != nil {
		c.metricsDropper.DropStaleMetricLabels()
		log.Printf("error reading inventory statistics from FRM: %s\n", err)
		return
	}

	for _, detail := range details {
		c.metricsDropper.CacheFreshMetricLabel(prometheus.Labels{
			"url":            frmAddress,
			"session_name":   sessionName,
			"container_name": detail.Name,
			"x":              strconv.FormatFloat(detail.Location.X, 'f', -1, 64),
			"y":              strconv.FormatFloat(detail.Location.Y, 'f', -1, 64),
			"z":              strconv.FormatFloat(detail.Location.Z, 'f', -1, 64),
		})
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

			CrateInventoryMax.WithLabelValues(
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

func (c *CrateInventoryCollector) DropCache() {}
