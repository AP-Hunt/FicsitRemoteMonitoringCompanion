package exporter_test

import (
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/exporter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProductionCollector", func() {
	var collector *exporter.ProductionCollector
	var url = "http://localhost:9080"
	var saveName = "default"

	BeforeEach(func() {
		FRMServer.Reset()
		collector = exporter.NewProductionCollector("/getProdStats")

		FRMServer.ReturnsProductionData([]exporter.ProductionDetails{
			{
				ItemName:           "Iron Rod",
				ProdPercent:        0.1,
				ConsPercent:        0.2,
				CurrentProduction:  10,
				CurrentConsumption: 40,
				MaxProd:            100.0,
				MaxConsumed:        200.0,
			},
		})
	})

	AfterEach(func() {
		collector = nil
	})

	Describe("Current item production & consumption metrics", func() {
		It("sets the 'items_produced_per_min' metric with the right labels", func() {
			collector.Collect(url, saveName)

			val, err := gaugeValue(exporter.ItemsProducedPerMin, "Iron Rod", url, saveName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(10)))
		})

		It("sets the 'items_consumed_per_min' metric with the right labels", func() {
			collector.Collect(url, saveName)

			val, err := gaugeValue(exporter.ItemsConsumedPerMin, "Iron Rod", url, saveName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(40)))
		})
	})

	Describe("Item production & consumption efficiency metrics", func() {
		It("sets the 'item_production_capacity_pc' metric with the right labels", func() {
			collector.Collect(url, saveName)

			val, err := gaugeValue(exporter.ItemProductionCapacityPercent, "Iron Rod", url, saveName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(0.1)))
		})

		It("sets the 'item_consumption_capacity_pc' metric with the right labels", func() {
			collector.Collect(url, saveName)

			val, err := gaugeValue(exporter.ItemConsumptionCapacityPercent, "Iron Rod", url, saveName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(0.2)))
		})
	})

	Describe("Item production & consumption capacity metrics", func() {
		It("sets the 'item_production_capacity_per_min' metric with the right labels", func() {
			collector.Collect(url, saveName)

			val, err := gaugeValue(exporter.ItemProductionCapacityPerMinute, "Iron Rod", url, saveName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(100)))
		})

		It("sets the 'item_consumption_capacity_per_min' metric with the right labels", func() {
			collector.Collect(url, saveName)

			val, err := gaugeValue(exporter.ItemConsumptionCapacityPerMinute, "Iron Rod", url, saveName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(200)))
		})

	})
})
