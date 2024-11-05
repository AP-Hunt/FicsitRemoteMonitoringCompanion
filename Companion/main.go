package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"

	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/exporter"
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/prometheus"
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/realtime_map"
)

var Version = "0.0.0-dev"

func lookupEnvWithDefault(variable string, defaultVal string) string {
	val, exist := os.LookupEnv(variable)
	if exist {
		return val
	}
	return defaultVal
}

func main() {

	var frmHostname string
	flag.StringVar(&frmHostname, "hostname", "localhost", "hostname of Ficsit Remote Monitoring webserver")
	var frmPort string
	flag.StringVar(&frmPort, "port", "8080", "port of Ficsit Remote Monitoring webserver")

	var frmHostnames string
	flag.StringVar(&frmHostnames, "hostnames", "", "comma separated values of multiple Ficsit Remote Monitoring webservers, of the form http://myserver1:8080,http://myserver2:8080. If defined, this will be used instead of hostname+port")

	var genReadme bool
	flag.BoolVar(&genReadme, "GenerateReadme", false, "Generate readme and exit")
	var noProm bool
	flag.BoolVar(&noProm, "noprom", false, "Do not run prometheus with the app.")
	flag.Parse()

	frmHostname = lookupEnvWithDefault("FRM_HOST", frmHostname)
	frmPort = lookupEnvWithDefault("FRM_PORT", frmPort)
	frmHostnames = lookupEnvWithDefault("FRM_HOSTS", frmHostnames)

	if genReadme {
		generateReadme()
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
		frmUrls = append(frmUrls, "http://"+frmHostname+":"+frmPort)
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

func generateReadme() {

	tpl := template.New("readme.tpl.md")
	tpl.Funcs(template.FuncMap{
		"List": strings.Join,
	})

	tpl = template.Must(tpl.ParseFiles("../readme/readme.tpl.md"))

	err := tpl.Execute(
		os.Stdout,
		struct {
			Metrics []exporter.MetricVectorDetails
		}{
			Metrics: exporter.RegisteredMetricVectors,
		},
	)
	if err != nil {
		fmt.Printf("Error writing readme", err)
	}
}
