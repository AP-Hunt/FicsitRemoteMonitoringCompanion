package exporter

import (
	"log"
	"regexp"
	"strconv"
)

type ProductionCollector struct {
	FRMAddress string
}

var prodPerMinRegex = regexp.MustCompile(`P: (?P<prod_current>[\d.]+)/(?P<prod_capacity>[\d.]+)/min - C: (?P<cons_current>[\d.]+)/(?P<cons_capacity>[\d.]+)/min`)

type ProductionDetails struct {
	ItemName           string   `json:"ItemName"`
	ProdPerMin         string   `json:"ProdPerMin"`
	ProdPercent        *float64 `json:"ProdPercent"`
	ConsPercent        *float64 `json:"ConsPercent"`
	CurrentProduction  float64  `json:"CurrentProd"`
	CurrentConsumption float64  `json:"CurrentConsumed"`
	MaxProd            float64  `json:"MaxProd"`
	MaxConsumed        float64  `json:"MaxConsumed"`
}

func (pd *ProductionDetails) parseProdPerMin() (bool, map[string]string) {
	match := prodPerMinRegex.FindStringSubmatch(pd.ProdPerMin)

	if len(match) < 1 {
		return false, nil
	}

	paramsMap := make(map[string]string)
	for i, name := range prodPerMinRegex.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}

	return true, paramsMap
}

func (pd *ProductionDetails) ItemProductionCapacity() *float64 {
	hasMatched, params := pd.parseProdPerMin()

	if !hasMatched {
		return nil
	}

	value := params["prod_capacity"]

	v, err := strconv.ParseFloat(value, 64)

	if err != nil {
		return nil
	}

	return &v
}

func (pd *ProductionDetails) ItemConsumptionCapacity() *float64 {
	hasMatched, params := pd.parseProdPerMin()

	if !hasMatched {
		return nil
	}

	value := params["cons_capacity"]

	v, err := strconv.ParseFloat(value, 64)

	if err != nil {
		return nil
	}

	return &v
}

func NewProductionCollector(frmAddress string) *ProductionCollector {
	return &ProductionCollector{
		FRMAddress: frmAddress,
	}
}

func (c *ProductionCollector) Collect() {
	details := []ProductionDetails{}
	err := retrieveData(c.FRMAddress, &details)
	if err != nil {
		log.Printf("error reading production statistics from FRM: %s\n", err)
		return
	}

	for _, d := range details {
		ItemsProducedPerMin.WithLabelValues(d.ItemName).Set(d.CurrentProduction)
		ItemsConsumedPerMin.WithLabelValues(d.ItemName).Set(d.CurrentConsumption)

		ItemProductionCapacityPercent.WithLabelValues(d.ItemName).Set(*d.ProdPercent)
		ItemConsumptionCapacityPercent.WithLabelValues(d.ItemName).Set(*d.ConsPercent)

		prodCapacity := d.ItemProductionCapacity()
		consCapacity := d.ItemConsumptionCapacity()

		if prodCapacity != nil {
			ItemProductionCapacityPerMinute.WithLabelValues(d.ItemName).Set(*prodCapacity)
		}

		if consCapacity != nil {
			ItemConsumptionCapacityPerMinute.WithLabelValues(d.ItemName).Set(*consCapacity)
		}
	}
}
