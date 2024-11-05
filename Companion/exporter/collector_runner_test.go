package exporter_test

import (
	"context"
	"time"

	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/exporter"
	. "github.com/onsi/ginkgo/v2"
	"github.com/coder/quartz"
	. "github.com/onsi/gomega"
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
			timeout, _ := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			testTime := quartz.NewMock(GinkgoTB())
			exporter.Clock = testTime
			trap := testTime.Trap().TickerFunc()
			defer trap.Close()

			c1 := NewTestCollector()
			c2 := NewTestCollector()
			runner := exporter.NewCollectorRunner(ctx, url, c1, c2)
			go runner.Start()
			call := trap.MustWait(timeout)
			call.Release()

			for i := 0; i < 2; i++ {
				_, w := testTime.AdvanceNext()
				w.MustWait(ctx)
			}
			Expect(c1.counter).To(Equal(3))
			Expect(c2.counter).To(Equal(3))
		})

		It("sanitizes session name", func() {
			Expect(exporter.SanitizeSessionName(`it's giving -- 123456!@#$%^&*() yo hollar "'"`)).To(Equal(`its giving  123456 yo hollar ` ))
		})
	})
})
