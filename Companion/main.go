package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"

	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/m/v2/exporter"
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/m/v2/prometheus"
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/m/v2/realtime_map"
)

var Version = "0.0.0-dev"

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "-ShowMetrics" {
		exportMetrics()
		os.Exit(0)
	}

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

	fmt.Printf(`
Ficsit Remote Monitoring Companion (v%s)

To access the realtime map visit:
http://localhost:8000/

To access metrics in Prometheus visit:
http://localhost:9090/

Press Ctrl + C to exit.
`, Version)

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

func exportMetrics() {
	tpl := template.New("metrics_table")
	tpl.Funcs(template.FuncMap{
		"List": strings.Join,
	})

	tpl, err := tpl.Parse(`
<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Description</th>
            <th>Labels</th>
        </tr>
    </thead>
    <tbody>
		{{range .Metrics}}
        <tr>
            <td>{{.Name}}</td>
            <td>{{.Help}}</td>
            <td>{{List .Labels ", "}}</td>
        </tr>
		{{ end -}}
	</tbody>
</table>
`)
	if err != nil {
		fmt.Printf("Error generating metrics table: %s", err)
		os.Exit(1)
	}

	tpl.Execute(
		os.Stdout,
		struct {
			Metrics []exporter.MetricVectorDetails
		}{
			Metrics: exporter.RegisteredMetricVectors,
		},
	)
}
