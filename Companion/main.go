package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"

	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/m/v2/exporter"
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/m/v2/prometheus"
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/m/v2/realtime_map"
)

func main() {
	logFile, err := createLogFile()
	if err != nil {
		fmt.Printf("error creating log file: %s", err)
		os.Exit(1)
	}
	log.Default().SetOutput(logFile)

	// Create exporter
	promExporter := exporter.NewPrometheusExporter("http://localhost:8090")

	// Create prometheus
	prom, err := prometheus.NewPrometheusWrapper()
	if err != nil {
		fmt.Printf("error preparing prometheus: %s", err)
		os.Exit(1)
	}

	// Create map server
	mapServ, err := realtime_map.NewMapServer()
	if err != nil {
		fmt.Printf("error preparing dynamic map: %s", err)
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

	// Start map
	mapServ.Start()

	fmt.Print(`
Ficsit Remote Monitoring Companion

To access the realtime map visit:
http://localhost:8000/

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

	// Stop map
	mapServ.Stop()

	fmt.Println("Exiting.")
	os.Exit(0)
}

func createLogFile() (*os.File, error) {
	curExePath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	curExeDir := filepath.Dir(curExePath)

	if err != nil {
		return nil, err
	}

	return os.Create(path.Join(curExeDir, "frmc.log"))
}
