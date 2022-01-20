package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/m/v2/exporter"
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/m/v2/prometheus"
)

func main() {
	// Create exporter
	promExporter := exporter.NewPrometheusExporter("http://localhost:8090")

	// Create prometheus
	prom, err := prometheus.NewPrometheusWrapper()
	if err != nil {
		fmt.Printf("error preparing prometheus: %s", err)
		os.Exit(1)
	}

	// Start prometheus
	err = prom.Start()
	if err != nil {
		fmt.Printf("error starting prometheus: %s", err)
		os.Exit(1)
	}

	// Start exporter
	promExporter.Start()

	fmt.Print(`
Ficsit Remote Monitoring Companion

To access metrics in Prometheus visit:
http://localhost:9090/

Press Ctrl + C to exit.
`)

	// Wait for an interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	// Stop the exporter
	err = promExporter.Stop()
	if err != nil {
		fmt.Printf("error stopping prometheus exporter: %s", err)
	}

	// Stop prometheus
	err = prom.Stop()
	if err != nil {
		fmt.Printf("error stopping prometheus: %s", err)
	}

	fmt.Println("Exiting.")
	os.Exit(0)
}
