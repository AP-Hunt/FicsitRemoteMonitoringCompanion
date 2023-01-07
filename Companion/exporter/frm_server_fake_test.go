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
	vehicleData      []exporter.VehicleDetails
	trainData        []exporter.TrainDetails
	droneData        []exporter.DroneStationDetails
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

	mux.Handle("/getProdStats", http.HandlerFunc(getStatsHandler(&fake.productionData)))
	mux.Handle("/getPower", http.HandlerFunc(getStatsHandler(&fake.powerData)))
	mux.Handle("/getFactory", http.HandlerFunc(getStatsHandler(&fake.factoryBuildings)))
	mux.Handle("/getDroneStation", http.HandlerFunc(getStatsHandler(&fake.droneData)))
	mux.Handle("/getTrains", http.HandlerFunc(getStatsHandler(&fake.trainData)))
	mux.Handle("/getVehicles", http.HandlerFunc(getStatsHandler(&fake.vehicleData)))

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

	for _, metric := range exporter.RegisteredMetrics {
		metric.Reset()
	}
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

func (e *FRMServerFake) ReturnsVehicleData(data []exporter.VehicleDetails) {
	e.vehicleData = data
}

func (e *FRMServerFake) ReturnsTrainData(data []exporter.TrainDetails) {
	e.trainData = data
}

func (e *FRMServerFake) ReturnsDroneStationData(data []exporter.DroneStationDetails) {
	e.droneData = data
}

func getStatsHandler(data any) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonBytes, err := json.Marshal(data)

		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.Write(jsonBytes)
	}
}
