package exporter

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type CollectorRunner struct {
	collectors  []Collector
	ctx         context.Context
	cancel      context.CancelFunc
	frmBaseUrl  string
	sessionName string
}

type Collector interface {
	Collect(string, string)
}

func NewCollectorRunner(ctx context.Context, frmBaseUrl string, collectors ...Collector) *CollectorRunner {
	ctx, cancel := context.WithCancel(ctx)
	return &CollectorRunner{
		ctx:         ctx,
		cancel:      cancel,
		collectors:  collectors,
		frmBaseUrl:  frmBaseUrl,
		sessionName: "default",
	}
}

func (c *CollectorRunner) updateSessionName() {
	//TODO: update session name
	// TODO: on update, we also should reset calculated route timings on trains and vehicles
	newSessionName := "default"
	if newSessionName != c.sessionName {
		for _, metric := range RegisteredMetrics {
			metric.DeletePartialMatch(prometheus.Labels{"url": c.frmBaseUrl, "session_name": c.sessionName})
		}
		c.sessionName = newSessionName
	}
}

func (c *CollectorRunner) Start() {
	c.updateSessionName()
	c.Collect(c.frmBaseUrl, c.sessionName)
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-Clock.After(5 * time.Second):
			c.updateSessionName()
			c.Collect(c.frmBaseUrl, c.sessionName)
		}
	}
}

func (c *CollectorRunner) Stop() {
	c.cancel()
}

func (c *CollectorRunner) Collect(server string, sessionName string) {
	for _, collector := range c.collectors {
		collector.Collect(server, sessionName)
	}
}
