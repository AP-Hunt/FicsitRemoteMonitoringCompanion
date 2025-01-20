package exporter_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/exporter"
)

type FRMServerFake struct {
	server                   *httptest.Server
	productionData           []exporter.ProductionDetails
	powerData                []exporter.PowerDetails
	factoryBuildings         []exporter.BuildingDetail
	vehicleData              []exporter.VehicleDetails
	trainData                []exporter.TrainDetails
	droneData                []exporter.DroneStationDetails
	vehicleStationData       []exporter.VehicleStationDetails
	trainStationData         []exporter.TrainStationDetails
	resourceSinkData         []exporter.ResourceSinkDetails
	gloalResourceSinkData    []exporter.GlobalSinkDetails
	gloalExplorationSinkData []exporter.GlobalSinkDetails
	sessionInfoData          exporter.SessionInfo
	pumpData                 []exporter.PumpDetails
	extractorData            []exporter.ExtractorDetails
	portalData               []exporter.PortalDetails
	hypertubeData            []exporter.HypertubeDetails
	frackingData             []exporter.FrackingDetails
}

func NewFRMServerFake() *FRMServerFake {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	fake := &FRMServerFake{
		server: server,
	}

	mux.Handle("/getProdStats", http.HandlerFunc(getStatsHandler(&fake.productionData)))
	mux.Handle("/getPower", http.HandlerFunc(getStatsHandler(&fake.powerData)))
	mux.Handle("/getFactory", http.HandlerFunc(getStatsHandler(&fake.factoryBuildings)))
	mux.Handle("/getDroneStation", http.HandlerFunc(getStatsHandler(&fake.droneData)))
	mux.Handle("/getTrains", http.HandlerFunc(getStatsHandler(&fake.trainData)))
	mux.Handle("/getVehicles", http.HandlerFunc(getStatsHandler(&fake.vehicleData)))
	mux.Handle("/getTruckStation", http.HandlerFunc(getStatsHandler(&fake.vehicleStationData)))
	mux.Handle("/getTrainStation", http.HandlerFunc(getStatsHandler(&fake.trainStationData)))
	mux.Handle("/getResourceSink", http.HandlerFunc(getStatsHandler(&fake.gloalResourceSinkData)))
	mux.Handle("/getExplorationSink", http.HandlerFunc(getStatsHandler(&fake.gloalExplorationSinkData)))
	mux.Handle("/getResourceSinkBuilding", http.HandlerFunc(getStatsHandler(&fake.resourceSinkData)))
	mux.Handle("/getSessionInfo", http.HandlerFunc(getStatsHandler(&fake.sessionInfoData)))
	mux.Handle("/getPump", http.HandlerFunc(getStatsHandler(&fake.pumpData)))
	mux.Handle("/getExtractor", http.HandlerFunc(getStatsHandler(&fake.extractorData)))
	mux.Handle("/getPortal", http.HandlerFunc(getStatsHandler(&fake.portalData)))
	mux.Handle("/getHypertube", http.HandlerFunc(getStatsHandler(&fake.hypertubeData)))
	mux.Handle("/getFrackingActivator", http.HandlerFunc(getStatsHandler(&fake.frackingData)))

	return fake
}

func (e *FRMServerFake) Stop() {
	e.server.Close()
	e.Reset()
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

func (e *FRMServerFake) ReturnsVehicleStationData(data []exporter.VehicleStationDetails) {
	e.vehicleStationData = data
}

func (e *FRMServerFake) ReturnsTrainStationData(data []exporter.TrainStationDetails) {
	e.trainStationData = data
}

func (e *FRMServerFake) ReturnsResourceSinkData(data []exporter.ResourceSinkDetails) {
	e.resourceSinkData = data
}

func (e *FRMServerFake) ReturnsGlobalResourceSinkData(data []exporter.GlobalSinkDetails) {
	e.gloalResourceSinkData = data
}

func (e *FRMServerFake) ReturnsGlobalExplorationSinkData(data []exporter.GlobalSinkDetails) {
	e.gloalExplorationSinkData = data
}

func (e *FRMServerFake) ReturnsPumpData(data []exporter.PumpDetails) {
	e.pumpData = data
}

func (e *FRMServerFake) ReturnsExtractorData(data []exporter.ExtractorDetails) {
	e.extractorData = data
}

func (e *FRMServerFake) ReturnsPortalData(data []exporter.PortalDetails) {
	e.portalData = data
}

func (e *FRMServerFake) ReturnsHypertubeData(data []exporter.HypertubeDetails) {
	e.hypertubeData = data
}

func (e *FRMServerFake) ReturnsFrackingData(data []exporter.FrackingDetails) {
	e.frackingData = data
}

func (e *FRMServerFake) ReturnsSessionInfoData(data exporter.SessionInfo) {
	e.sessionInfoData = data
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
