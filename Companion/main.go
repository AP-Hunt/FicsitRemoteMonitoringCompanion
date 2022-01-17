package main

import (
	"bufio"
	"fmt"
	"os"

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

	fmt.Println("Processes are running. Press enter to stop.")
	rdr := bufio.NewReader(os.Stdin)
	_, _ = rdr.ReadString('\n')
	fmt.Println("Exiting.")

	err = prom.Stop()
	if err != nil {
		fmt.Printf("error stopping prometheus: %s", err)
	}

	os.Exit(0)
}
