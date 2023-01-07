package exporter

import (
	"context"
	"time"
)

type CollectorRunner struct {
	collectors []Collector
	ctx        context.Context
	cancel     context.CancelFunc
}

type Collector interface {
	Collect()
}

func NewCollectorRunner(ctx context.Context, collectors ...Collector) *CollectorRunner {
	ctx, cancel := context.WithCancel(ctx)
	return &CollectorRunner{
		ctx:        ctx,
		cancel:     cancel,
		collectors: collectors,
	}
}

func (c *CollectorRunner) Start() {
	c.Collect()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-Clock.After(5 * time.Second):
			c.Collect()
		}
	}
}

func (c *CollectorRunner) Stop() {
	c.cancel()
}

func (c *CollectorRunner) Collect() {
	for _, collector := range c.collectors {
		collector.Collect()
	}
}
