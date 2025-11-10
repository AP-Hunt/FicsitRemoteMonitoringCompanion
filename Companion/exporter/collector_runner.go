package exporter

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func SanitizeSessionName(sessionName string) string {
	re := regexp.MustCompile(`[^\w\s]`)
	return re.ReplaceAllString(sessionName, "")
}

type CollectorRunner struct {
	collectors  []Collector
	ctx         context.Context
	cancel      context.CancelFunc
	frmBaseUrl  string
	sessionName string
}

type Collector interface {
	Collect(string, string)
	DropCache()
}

type SessionInfo struct {
	SessionName string `json:"SessionName"`
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
	details := SessionInfo{}
	err := retrieveData(c.frmBaseUrl, "/getSessionInfo", &details)
	if err != nil {
		log.Printf("error reading session name from FRM: %s\n", err)
		return
	}
	newSessionName := SanitizeSessionName(details.SessionName)
	if newSessionName != "" && newSessionName != c.sessionName {
		log.Printf("%s has a new session name: %s\n", c.frmBaseUrl, newSessionName)
		for _, metric := range RegisteredMetrics {
			metric.DeletePartialMatch(prometheus.Labels{"url": c.frmBaseUrl, "session_name": c.sessionName})
		}
		for _, collector := range c.collectors {
			collector.DropCache()
		}
		c.sessionName = newSessionName
	}
}

func (c *CollectorRunner) Start() error {
	c.updateSessionName()
	c.Collect(c.frmBaseUrl, c.sessionName)
	t := Clock.TickerFunc(c.ctx, 5*time.Second, func() error {
		c.updateSessionName()
		c.Collect(c.frmBaseUrl, c.sessionName)
		return nil
	})
	return t.Wait()
}

func (c *CollectorRunner) Stop() {
	c.cancel()
}

func (c *CollectorRunner) Collect(server string, sessionName string) {
	for _, collector := range c.collectors {
		collector.Collect(server, sessionName)
	}
}
