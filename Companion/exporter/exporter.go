package exporter

import (
	"context"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusExporter struct {
	server              *http.Server
	ctx                 context.Context
	cancel              context.CancelFunc
	productionCollector *ProductionCollector
	powerCollector      *PowerCollector
	buildingCollector   *FactoryBuildingCollector
}

func NewPrometheusExporter(frmApiHost string) *PrometheusExporter {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Handler: mux,
		Addr:    ":9000",
	}

	ctx, cancel := context.WithCancel(context.Background())
	productionCollector := NewProductionCollector(ctx, frmApiHost+"/getProdStats")
	powerCollector := NewPowerCollector(ctx, frmApiHost+"/getPower")
	buildingCollector := NewFactoryBuildingCollector(ctx, frmApiHost+"/getFactory")

	return &PrometheusExporter{
		server:              server,
		ctx:                 ctx,
		cancel:              cancel,
		productionCollector: productionCollector,
		powerCollector:      powerCollector,
		buildingCollector:   buildingCollector,
	}
}

func (e *PrometheusExporter) Start() {
	go e.productionCollector.Start()
	go e.powerCollector.Start()
	go e.buildingCollector.Start()
	go func() {
		log.Fatal(e.server.ListenAndServe())
	}()
}

func (e *PrometheusExporter) Stop() error {
	e.cancel()
	return e.server.Close()
}
