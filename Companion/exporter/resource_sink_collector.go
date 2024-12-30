package exporter

import (
	"log"
	"strconv"
)

type ResourceSinkCollector struct {
	buildingEndpoint string
	globalEndpoint   string
}

type ResourceSinkDetails struct {
	Location  Location  `json:"location"`
	PowerInfo PowerInfo `json:"PowerInfo"`
}

type ResourceSinkGlobalDetails struct {
	NumCoupon      int `json:"NumCoupon"`
	TotalPoints    int `json:"TotalPoints"`
	PointsToCoupon int `json:"PointsToCoupon"`
}

func NewResourceSinkCollector(buildingEndpoint, globalEndpoint string) *ResourceSinkCollector {
	return &ResourceSinkCollector{
		buildingEndpoint: buildingEndpoint,
		globalEndpoint:   globalEndpoint,
	}
}

func (c *ResourceSinkCollector) Collect(frmAddress string, sessionName string) {
	buildingDetails := []ResourceSinkDetails{}
	err := retrieveData(frmAddress+c.buildingEndpoint, &buildingDetails)
	if err != nil {
		log.Printf("error reading resource sink details statistics from FRM: %s\n", err)
		return
	}

	globalDetails := []ResourceSinkGlobalDetails{}
	err = retrieveData(frmAddress+c.globalEndpoint, &globalDetails)
	if err != nil {
		log.Printf("error reading global resource sink statistics from FRM: %s\n", err)
		return
	}

	for _, d := range globalDetails {
		ResourceSinkTotalPoints.WithLabelValues(frmAddress, sessionName).Set(float64(d.TotalPoints))
		ResourceSinkPointsToCoupon.WithLabelValues(frmAddress, sessionName).Set(float64(d.PointsToCoupon))
		ResourceSinkCollectedCoupons.WithLabelValues(frmAddress, sessionName).Set(float64(d.NumCoupon))
	}

	powerInfo := map[float64]float64{}
	maxPowerInfo := map[float64]float64{}
	for _, d := range buildingDetails {
		val, ok := powerInfo[d.PowerInfo.CircuitGroupId]
		if ok {
			powerInfo[d.PowerInfo.CircuitGroupId] = val + d.PowerInfo.PowerConsumed
		} else {
			powerInfo[d.PowerInfo.CircuitGroupId] = d.PowerInfo.PowerConsumed
		}
		val, ok = maxPowerInfo[d.PowerInfo.CircuitGroupId]
		if ok {
			maxPowerInfo[d.PowerInfo.CircuitGroupId] = val + d.PowerInfo.MaxPowerConsumed
		} else {
			maxPowerInfo[d.PowerInfo.CircuitGroupId] = d.PowerInfo.MaxPowerConsumed
		}
	}
	for circuitId, powerConsumed := range powerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		ResourceSinkPower.WithLabelValues(cid, frmAddress, sessionName).Set(powerConsumed)
	}
	for circuitId, powerConsumed := range maxPowerInfo {
		cid := strconv.FormatFloat(circuitId, 'f', -1, 64)
		ResourceSinkPowerMax.WithLabelValues(cid, frmAddress, sessionName).Set(powerConsumed)
	}
}

func (c *ResourceSinkCollector) DropCache() {}
