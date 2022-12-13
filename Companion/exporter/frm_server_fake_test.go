package exporter_test

import (
	"encoding/json"
	"net/http"

	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/m/v2/exporter"
)

type FRMServerFake struct {
	server           *http.Server
	productionData   []exporter.ProductionDetails
	powerData        []exporter.PowerDetails
	factoryBuildings []exporter.BuildingDetail
}

func NewFRMServerFake() *FRMServerFake {
	mux := http.NewServeMux()
	server := &http.Server{
		Handler: mux,
		Addr:    ":9080",
	}

	fake := &FRMServerFake{
		server: server,
	}

	mux.Handle("/getProdStats", http.HandlerFunc(fake.productionStatsHandler))
	mux.Handle("/getPower", http.HandlerFunc(fake.powerStatsHandler))
	mux.Handle("/getFactory", http.HandlerFunc(fake.factoryBuildingsHandler))

	return fake
}

func (f *FRMServerFake) Start() {

	go func() {
		f.server.ListenAndServe()
	}()
}

func (e *FRMServerFake) Stop() error {
	err := e.server.Close()
	e.Reset()
	return err
}

func (e *FRMServerFake) Reset() {
	e.productionData = nil
	e.powerData = nil

	exporter.ItemConsumptionCapacityPerMinute.Reset()
	exporter.ItemConsumptionCapacityPercent.Reset()
	exporter.ItemProductionCapacityPerMinute.Reset()
	exporter.ItemProductionCapacityPercent.Reset()
	exporter.ItemsConsumedPerMin.Reset()
	exporter.ItemsProducedPerMin.Reset()
	exporter.PowerConsumed.Reset()
	exporter.PowerCapacity.Reset()
	exporter.PowerMaxConsumed.Reset()
	exporter.BatteryDifferential.Reset()
	exporter.BatteryPercent.Reset()
	exporter.BatteryCapacity.Reset()
	exporter.BatterySecondsEmpty.Reset()
	exporter.BatterySecondsFull.Reset()
	exporter.FuseTriggered.Reset()
}

func (e *FRMServerFake) ReturnsProductionData(data []exporter.ProductionDetails) {
	e.productionData = data
}

func (e *FRMServerFake) ReturnsPowerData(data []exporter.PowerDetails) {
	e.powerData = data
}

func (e *FRMServerFake) ReturnsFactoryBuildings(data []exporter.BuildingDetail) {
	e.factoryBuildings = data
}

func (e *FRMServerFake) productionStatsHandler(w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := json.Marshal(e.productionData)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Write(jsonBytes)
}

func (e *FRMServerFake) powerStatsHandler(w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := json.Marshal(e.powerData)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Write(jsonBytes)
}

func (e *FRMServerFake) factoryBuildingsHandler(w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := json.Marshal(e.factoryBuildings)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Write(jsonBytes)
}
