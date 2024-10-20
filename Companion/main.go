package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/exporter"
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/prometheus"
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/realtime_map"
)

var Version = "0.0.0-dev"

func main() {

	var frmHostname string
	flag.StringVar(&frmHostname, "hostname", "localhost", "hostname of Ficsit Remote Monitoring webserver")
	var frmPort int
	flag.IntVar(&frmPort, "port", 8080, "port of Ficsit Remote Monitoring webserver")

	var frmHostnames string
	flag.StringVar(&frmHostnames, "hostnames", "", "comma separated values of multiple Ficsit Remote Monitoring webservers, of the form http://myserver1:8080,http://myserver2:8080. If defined, this will be used instead of hostname+port")

	var showMetrics bool
	flag.BoolVar(&showMetrics, "ShowMetrics", false, "Show metrics and exit")
	var noProm bool
	flag.BoolVar(&noProm, "noprom", false, "Do not run prometheus with the app.")
	flag.Parse()

	if showMetrics {
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
	frmUrls := []string{}
	if frmHostnames == "" {
		frmUrls = append(frmUrls, "http://" + frmHostname + ":" + strconv.Itoa(frmPort))
	} else {
		for _, frmServer := range strings.Split(frmHostnames, ",") {
			if !strings.HasPrefix(frmServer, "http://") && !strings.HasPrefix(frmServer, "https://") {
				frmServer = "http://" + frmServer
			}
			frmUrls = append(frmUrls, frmServer)
		}
	}
	var promExporter *exporter.PrometheusExporter
	promExporter = exporter.NewPrometheusExporter(frmUrls)

	var prom *prometheus.PrometheusWrapper
	if !noProm {
		// Create prometheus
		prom, err = prometheus.NewPrometheusWrapper()
		if err != nil {
			fmt.Printf("error preparing prometheus: %s", err)
			os.Exit(1)
		}
	}

	// Create map server
	mapServ, err := realtime_map.NewMapServer()
	if err != nil {
		fmt.Printf("error preparing dynamic map: %s", err)
		os.Exit(1)
	}

	// Start prometheus
	if !noProm {
		err = prom.Start()
		if err != nil {
			fmt.Printf("error starting prometheus: %s", err)
			os.Exit(1)
		}
	}

	// Start exporter
	promExporter.Start()

	// Start map
	mapServ.Start()

	fmt.Printf(`
Ficsit Remote Monitoring Companion (v%s)

To access the realtime map visit:
http://localhost:8000/?frmport=8080

    If you have configured Ficsit Remote Monitoring 
    to use a port other than 8080 for its web server, 
    change the "frmport" query string parameter to 
    match the port you chose and refresh the page.

To access metrics in Prometheus visit:
http://localhost:9090/

Press Ctrl + C to exit.
`, Version)

	// Wait for an interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	// Stop the exporter
	err = promExporter.Stop()
	if err != nil {
		fmt.Printf("error stopping prometheus exporter: %s", err)
	}

	// Stop prometheus
	if !noProm {
		err = prom.Stop()
		if err != nil {
			fmt.Printf("error stopping prometheus: %s", err)
		}
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
