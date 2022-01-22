package exporter

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FactoryBuildingCollector struct {
	FRMAddress string
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewFactoryBuildingCollector(ctx context.Context, frmAddress string) *FactoryBuildingCollector {
	ctx, cancel := context.WithCancel(ctx)

	return &FactoryBuildingCollector{
		FRMAddress: frmAddress,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (c *FactoryBuildingCollector) Start() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			c.Collect()
			time.Sleep(5 * time.Second)
		}
	}
}

func (c *FactoryBuildingCollector) Stop() {
	c.cancel()
}

func (c *FactoryBuildingCollector) Collect() {
	resp, err := http.Get(c.FRMAddress)

	if err != nil {
		log.Printf("error fetching factory buildings from FRM: %s\n", err)
		return
	}

	defer resp.Body.Close()

	details := []BuildingDetail{}
	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&details)
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

			efficiency := float64(0.0)
			if e, err := strconv.ParseFloat(prod.ProdPercent, 64); err == nil {
				efficiency = e
			}

			MachineItemsProducedEffiency.WithLabelValues(
				prod.Name,
				building.Building,
				strconv.FormatFloat(building.Location.X, 'f', -1, 64),
				strconv.FormatFloat(building.Location.Y, 'f', -1, 64),
				strconv.FormatFloat(building.Location.Z, 'f', -1, 64),
			).Set(efficiency)
		}
	}
}
