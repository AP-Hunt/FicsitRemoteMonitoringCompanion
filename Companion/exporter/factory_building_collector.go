package exporter

import (
	"log"
	"strconv"
	"math"
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

	powerInfo := map[float64]float64{}
	maxPowerInfo := map[float64]float64{}
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

		val, ok := powerInfo[building.PowerInfo.CircuitId]
		if ok {
			powerInfo[building.PowerInfo.CircuitId] = val + building.PowerInfo.PowerConsumed
		} else {
			powerInfo[building.PowerInfo.CircuitId] = building.PowerInfo.PowerConsumed
		}
		val, ok = maxPowerInfo[building.PowerInfo.CircuitId]
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
		maxBuildingPower = maxBuildingPower * (math.Pow(building.ManuSpeed / 100, 1.321928))
		if ok {
			maxPowerInfo[building.PowerInfo.CircuitId] = val + maxBuildingPower
		} else {
			maxPowerInfo[building.PowerInfo.CircuitId] = maxBuildingPower
		}
	}
	for circuitId, powerConsumed := range powerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		FactoryPower.WithLabelValues(cid).Set(powerConsumed)
	}
	for circuitId, powerConsumed := range maxPowerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		FactoryPowerMax.WithLabelValues(cid).Set(powerConsumed)
	}
}
