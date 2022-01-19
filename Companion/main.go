package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/m/v2/exporter"
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/m/v2/prometheus"
)

func main() {
	prom, err := prometheus.NewPrometheusWrapper()

	if err != nil {
		fmt.Printf("error starting prometheus: %s", err)
		os.Exit(1)
	}

	err = prom.Start()
	if err != nil {
		fmt.Printf("error starting prometheus: %s", err)
		os.Exit(1)
	}

	promExporter := exporter.NewPrometheusExporter("http://localhost:8090")
	promExporter.Start()

	fmt.Println("Processes are running. Press enter to stop.")
	rdr := bufio.NewReader(os.Stdin)
	_, _ = rdr.ReadString('\n')
	fmt.Println("Exiting.")

	err = promExporter.Stop()
	if err != nil {
		fmt.Printf("error stopping prometheus exporter: %s", err)
	}

	err = prom.Stop()
	if err != nil {
		fmt.Printf("error stopping prometheus: %s", err)
	}

	os.Exit(0)
}
