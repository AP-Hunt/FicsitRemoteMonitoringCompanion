package exporter

import (
	"context"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusExporter struct {
	server          *http.Server
	ctx             context.Context
	cancel          context.CancelFunc
	collectorRunner *CollectorRunner
}

func NewPrometheusExporter(frmApiHost string) *PrometheusExporter {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Handler: mux,
		Addr:    ":9000",
	}

	ctx, cancel := context.WithCancel(context.Background())
	productionCollector := NewProductionCollector(frmApiHost + "/getProdStats")
	powerCollector := NewPowerCollector(frmApiHost + "/getPower")
	buildingCollector := NewFactoryBuildingCollector(frmApiHost + "/getFactory")
	vehicleCollector := NewVehicleCollector(frmApiHost + "/getVehicles")
	droneCollector := NewDroneStationCollector(frmApiHost + "/getDroneStation")
	vehicleStationCollector := NewVehicleStationCollector(frmApiHost + "/getTruckStation")


	trackedStations := &(map[string]TrainStationDetails{})
	trainCollector := NewTrainCollector(frmApiHost + "/getTrains", trackedStations)
	trainStationCollector := NewTrainStationCollector(frmApiHost + "/getTrainStation", trackedStations)
	collectorRunner := NewCollectorRunner(ctx, productionCollector, powerCollector, buildingCollector, vehicleCollector, trainCollector, droneCollector, vehicleStationCollector, trainStationCollector)

	return &PrometheusExporter{
		server:          server,
		ctx:             ctx,
		cancel:          cancel,
		collectorRunner: collectorRunner,
	}
}

func (e *PrometheusExporter) Start() {
	go e.collectorRunner.Start()
	go func() {
		e.server.ListenAndServe()
		log.Println("stopping exporter")
	}()
}

func (e *PrometheusExporter) Stop() error {
	e.cancel()
	return e.server.Close()
}
