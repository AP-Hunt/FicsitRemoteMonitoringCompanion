package exporter_test

import (
	"context"

	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/exporter"
	"github.com/benbjohnson/clock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

type TestCollector struct {
	counter int
}

func NewTestCollector() *TestCollector {
	return &TestCollector{
		counter: 0,
	}
}
func (t *TestCollector) Collect(url string, sessionName string) {
	t.counter = t.counter + 1
}
func (t *TestCollector) DropCache() {}

var _ = Describe("CollectorRunner", func() {
	var url string

	BeforeEach(func() {
		FRMServer.Reset()
		url = FRMServer.server.URL
		FRMServer.ReturnsSessionInfoData(exporter.SessionInfo{
			SessionName: "test",
		})
	})

	Describe("Basic Functionality", func() {
		It("runs on init and on each timeout", func() {
			ctx, cancel := context.WithCancel(context.Background())
			testTime := clock.NewMock()
			exporter.Clock = testTime

			c1 := NewTestCollector()
			c2 := NewTestCollector()
			runner := exporter.NewCollectorRunner(ctx, url, c1, c2)
			go runner.Start()
			testTime.Add(5 * time.Second)
			testTime.Add(5 * time.Second)
			testTime.Add(5 * time.Second)
			testTime.Add(4 * time.Second)
			cancel()
			Eventually(c1.counter).Should(Equal(3))
			Eventually(c2.counter).Should(Equal(3))
		})

		It("does not run after being canceled", func() {
			ctx, cancel := context.WithCancel(context.Background())
			testTime := clock.NewMock()
			exporter.Clock = testTime

			c1 := NewTestCollector()
			runner := exporter.NewCollectorRunner(ctx, url, c1)
			go runner.Start()
			testTime.Add(5 * time.Second)
			Eventually(c1.counter).Should(Equal(1))
			cancel()
			testTime.Add(5 * time.Second)
			testTime.Add(5 * time.Second)
			Expect(c1.counter).To(Equal(1))
		})
	})
})
