package exporter

import (
	"log"
	"math"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type FactoryBuildingCollector struct {
	endpoint       string
	metricsDropper *MetricsDropper
}

func NewFactoryBuildingCollector(endpoint string) *FactoryBuildingCollector {
	return &FactoryBuildingCollector{
		endpoint: endpoint,
		metricsDropper: NewMetricsDropper(
			MachineItemsProducedPerMin,
			MachineItemsProducedEffiency,
		),
	}
}

func (c *FactoryBuildingCollector) Collect(frmAddress string, sessionName string) {
	details := []BuildingDetail{}
	err := retrieveData(frmAddress+c.endpoint, &details)
	if err != nil {
		c.metricsDropper.DropStaleMetricLabels()
		log.Printf("error reading factory buildings from FRM: %s\n", err)
		return
	}

	powerInfo := map[float64]float64{}
	maxPowerInfo := map[float64]float64{}
	for _, building := range details {
		c.metricsDropper.CacheFreshMetricLabel(prometheus.Labels{"url": frmAddress, "session_name": sessionName, "machine_name": building.Building,
			"x": strconv.FormatFloat(building.Location.X, 'f', -1, 64),
			"y": strconv.FormatFloat(building.Location.Y, 'f', -1, 64),
			"z": strconv.FormatFloat(building.Location.Z, 'f', -1, 64),
		})
		for _, prod := range building.Production {
			MachineItemsProducedPerMin.WithLabelValues(
				prod.Name,
				building.Building,
				strconv.FormatFloat(building.Location.X, 'f', -1, 64),
				strconv.FormatFloat(building.Location.Y, 'f', -1, 64),
				strconv.FormatFloat(building.Location.Z, 'f', -1, 64),
				frmAddress, sessionName,
			).Set(prod.CurrentProd)

			MachineItemsProducedEffiency.WithLabelValues(
				prod.Name,
				building.Building,
				strconv.FormatFloat(building.Location.X, 'f', -1, 64),
				strconv.FormatFloat(building.Location.Y, 'f', -1, 64),
				strconv.FormatFloat(building.Location.Z, 'f', -1, 64),
				frmAddress, sessionName,
			).Set(prod.ProdPercent)
		}

		val, ok := powerInfo[building.PowerInfo.CircuitGroupId]
		if ok {
			powerInfo[building.PowerInfo.CircuitGroupId] = val + building.PowerInfo.PowerConsumed
		} else {
			powerInfo[building.PowerInfo.CircuitGroupId] = building.PowerInfo.PowerConsumed
		}
		val, ok = maxPowerInfo[building.PowerInfo.CircuitGroupId]
		maxBuildingPower := 0.0
		switch building.Building {
		case "Smelter":
			maxBuildingPower = SmelterPower
			break
		case "Constructor":
			maxBuildingPower = ConstructorPower
			break
		case "Assembler":
			maxBuildingPower = AssemblerPower
			break
		case "Manufacturer":
			maxBuildingPower = ManufacturerPower
			break
		case "Blender":
			maxBuildingPower = BlenderPower
			break
		case "Refinery":
			maxBuildingPower = RefineryPower
			break
		case "Particle Accelerator":
			maxBuildingPower = ParticleAcceleratorPower
			break
		}
		//update max power from clock speed
		// see https://satisfactory.wiki.gg/wiki/Clock_speed#Clock_speed_for_production_buildings for power info
		maxBuildingPower = maxBuildingPower * (math.Pow(building.ManuSpeed/100, 1.321928))
		if ok {
			maxPowerInfo[building.PowerInfo.CircuitGroupId] = val + maxBuildingPower
		} else {
			maxPowerInfo[building.PowerInfo.CircuitGroupId] = maxBuildingPower
		}
	}
	c.metricsDropper.DropStaleMetricLabels()
	for circuitId, powerConsumed := range powerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		FactoryPower.WithLabelValues(cid, frmAddress, sessionName).Set(powerConsumed)
	}
	for circuitId, powerConsumed := range maxPowerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		FactoryPowerMax.WithLabelValues(cid, frmAddress, sessionName).Set(powerConsumed)
	}
}
