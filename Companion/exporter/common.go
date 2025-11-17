package exporter

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/coder/quartz"
)

// Struct for the JSON body of POST requests
type FrmApiRequest struct {
	Function string `json:"function"`
	Endpoint string `json:"endpoint"`
}

// Reusable HTTP client for HTTPS w/ self-signed certs
var tlsClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

var timeRegex = regexp.MustCompile(`\d\d:\d\d:\d\d`)

var Clock = quartz.NewReal()

func parseTimeSeconds(timeStr string) *float64 {
	match := timeRegex.FindStringSubmatch(timeStr)
	if len(match) < 1 {
		return nil
	}
	parts := strings.Split(match[0], ":")
	duration := parts[0] + "h" + parts[1] + "m" + parts[2] + "s"
	t, _ := time.ParseDuration(duration)
	seconds := t.Seconds()
	return &seconds
}

func parseBool(b bool) float64 {
	if b {
		return 1
	} else {
		return 0
	}
}

func retrieveDataViaGET(frmAddress string, details any) error {
	resp, err := http.Get(frmAddress)

	if err != nil {
		log.Printf("error fetching statistics from FRM: %s\n", err)
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("non-200 returned when retireving data: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&details)
	return err
}

func retrieveDataViaPOST(frmApiUrl string, endpointName string, details any) error {
	reqBody := FrmApiRequest{
		Function: "frm",
		Endpoint: endpointName,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %s\n", err.Error())
	}

	req, err := http.NewRequest("POST", frmApiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating POST request: %s\n", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := tlsClient.Do(req)
	if err != nil {
		return fmt.Errorf("error fetching statistics from FRM: %s\n", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("non-200 returned when retrieving data: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&details)
	return err
}

func retrieveData(frmAddress string, endpoint string, details any) error {
	u, err := url.Parse(frmAddress)
	if err != nil {
		return fmt.Errorf("invalid FRM address URL: %s", err)
	}

	// Check if we're using the Dedicated Server API (1.1)
	if strings.HasSuffix(u.Path, "/api/v1") {
		// Dedicated server mode.
		// frmAddress is "https://host:7777/api/v1"
		// endpoint is "/getPower", so we strip the slash
		endpointName := strings.TrimPrefix(endpoint, "/")
		return retrieveDataViaPOST(frmAddress, endpointName, details)
	} else {
		// Web server mode.
		// frmAddress is "http://host:8080"
		// endpoint is "/getPower"
		return retrieveDataViaGET(frmAddress+endpoint, details)
	}
}
