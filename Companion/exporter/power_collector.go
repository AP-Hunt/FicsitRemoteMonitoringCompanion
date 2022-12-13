package exporter

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"strconv"
	"time"
)

type PowerCollector struct {
	FRMAddress string
	ctx        context.Context
	cancel     context.CancelFunc
}

var timeRegex = regexp.MustCompile(`\d\d:\d\d:\d\d`)

type PowerDetails struct {
	CircuitId           float64 `json:"CircuitID"`
	PowerConsumed       float64 `json:"PowerConsumed"`
	PowerCapacity       float64 `json:"PowerCapacity"`
	PowerMaxConsumed    float64 `json:"PowerMaxConsumed"`
	BatteryDifferential float64 `json:"BatteryDifferential"`
	BatteryPercent      float64 `json:"BatteryPercent"`
	BatteryCapacity     float64 `json:"BatteryCapacity"`
	BatteryTimeEmpty    string  `json:"BatteryTimeEmpty"`
	BatteryTimeFull     string  `json:"BatteryTimeFull"`
	FuseTriggered       bool    `json:"FuseTriggered"`
}

func parseTimeSeconds(timeStr string) (bool, float64) {
	match := timeRegex.FindStringSubmatch(timeStr)
	if len(match) < 1 {
		return false, 0
	}
	parts := strings.Split(match[0], ":")
	duration := parts[0] + "h" + parts[1] + "m" + parts[2] + "s"
	t, _ := time.ParseDuration(duration)
	return true, t.Seconds()
}

func (pd *PowerDetails) parseBatteryTimeEmptySeconds() *float64 {
	matched, params := parseTimeSeconds(pd.BatteryTimeEmpty)
	if !matched {
		return nil
	}
	return &params
}

func (pd *PowerDetails) parseBatteryTimeFullSeconds() *float64 {
	matched, params := parseTimeSeconds(pd.BatteryTimeFull)
	if !matched {
		return nil
	}
	return &params
}

func (pd *PowerDetails) parseFuseTriggered() float64 {
	if pd.FuseTriggered {
		return 1
	} else {
		return 0
	}
}

func NewPowerCollector(ctx context.Context, frmAddress string) *PowerCollector {
	ctx, cancel := context.WithCancel(ctx)
	return &PowerCollector{
		FRMAddress: frmAddress,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (c *PowerCollector) Start() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			c.Collect()
			time.Sleep(5 * time.Second)
		}
	}
}

func (c *PowerCollector) Stop() {
	c.cancel()
}

func (c *PowerCollector) Collect() {
	resp, err := http.Get(c.FRMAddress)

	if err != nil {
		log.Printf("error fetching power statistics from FRM: %s\n", err)
		return
	}

	defer resp.Body.Close()

	details := []PowerDetails{}
	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&details)
	if err != nil {
		log.Printf("error reading power statistics from FRM: %s\n", err)
		return
	}

	for _, d := range details {
		circuitId := strconv.FormatFloat(d.CircuitId, 'f', -1, 64)
		PowerConsumed.WithLabelValues(circuitId).Set(d.PowerConsumed)
		PowerCapacity.WithLabelValues(circuitId).Set(d.PowerCapacity)
		PowerMaxConsumed.WithLabelValues(circuitId).Set(d.PowerMaxConsumed)
		BatteryDifferential.WithLabelValues(circuitId).Set(d.BatteryDifferential)
		BatteryPercent.WithLabelValues(circuitId).Set(d.BatteryPercent)
		BatteryCapacity.WithLabelValues(circuitId).Set(d.BatteryCapacity)
		batterySecondsRemaining := d.parseBatteryTimeEmptySeconds()
		if batterySecondsRemaining != nil {
			BatterySecondsEmpty.WithLabelValues(circuitId).Set(*batterySecondsRemaining)
		}
		batterySecondsFull := d.parseBatteryTimeFullSeconds()
		if batterySecondsFull != nil {
			BatterySecondsFull.WithLabelValues(circuitId).Set(*batterySecondsFull)
		}
		fuseTriggered := d.parseFuseTriggered()
		FuseTriggered.WithLabelValues(circuitId).Set(fuseTriggered)
	}
}
