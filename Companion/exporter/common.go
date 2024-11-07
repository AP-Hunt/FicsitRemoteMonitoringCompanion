package exporter

import (
	"encoding/json"
	"fmt"
	"github.com/coder/quartz"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

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

func retrieveData(frmAddress string, details any) error {
	resp, err := http.Get(frmAddress)

	if err != nil {
		log.Printf("error fetching statistics from FRM: %s\n", err)
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("non-200 returned when retireving data: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&details)
	return err
}
