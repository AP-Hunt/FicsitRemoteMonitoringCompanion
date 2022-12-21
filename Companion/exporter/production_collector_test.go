package exporter_test

import (
	"fmt"

	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/m/v2/exporter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func prodPerMin(prodActual float32, prodCapacity float32, consActual float32, consCapacity float32) string {
	return fmt.Sprintf("P: %f/%f/min - C: %f/%f/min", prodActual, prodCapacity, consActual, consCapacity)
}

func prodPerMinWithoutCapacity(prodActual float32, consActual float32) string {
	// Inconsistent spacing in format string is correct according to the API
	return fmt.Sprintf("P:%f/min - C: %f/min", prodActual, consActual)
}

var _ = Describe("ProductionCollector", func() {
	var collector *exporter.ProductionCollector

	BeforeEach(func() {
		FRMServer.Reset()
		collector = exporter.NewProductionCollector("http://localhost:9080/getProdStats")

		FRMServer.ReturnsProductionData([]exporter.ProductionDetails{
			{
				ItemName:           "Iron Rod",
				ProdPerMin:         prodPerMin(10, 100, 40, 200),
				ProdPercent:        f(0.1),
				ConsPercent:        f(0.2),
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
			collector.Collect()

			val, err := gaugeValue(exporter.ItemsProducedPerMin, "Iron Rod")
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(10)))
		})

		It("sets the 'items_consumed_per_min' metric with the right labels", func() {
			collector.Collect()

			val, err := gaugeValue(exporter.ItemsConsumedPerMin, "Iron Rod")
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(40)))
		})
	})

	Describe("Item production & consumption efficiency metrics", func() {
		It("sets the 'item_production_capacity_pc' metric with the right labels", func() {
			collector.Collect()

			val, err := gaugeValue(exporter.ItemProductionCapacityPercent, "Iron Rod")
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(0.1)))
		})

		It("sets the 'item_consumption_capacity_pc' metric with the right labels", func() {
			collector.Collect()

			val, err := gaugeValue(exporter.ItemConsumptionCapacityPercent, "Iron Rod")
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(0.2)))
		})
	})

	Describe("Item production & consumption capacity metrics", func() {
		It("sets the 'item_production_capacity_per_min' metric with the right labels", func() {
			collector.Collect()

			val, err := gaugeValue(exporter.ItemProductionCapacityPerMinute, "Iron Rod")
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(100)))
		})

		It("sets the 'item_consumption_capacity_per_min' metric with the right labels", func() {
			collector.Collect()

			val, err := gaugeValue(exporter.ItemConsumptionCapacityPerMinute, "Iron Rod")
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(200)))
		})

		Describe("when capacity metrics aren't supplied", func() {

			BeforeEach(func() {
				FRMServer.ReturnsProductionData([]exporter.ProductionDetails{
					{
						ItemName:           "Iron Rod",
						ProdPerMin:         prodPerMinWithoutCapacity(10, 40),
						ProdPercent:        f(1.0),
						ConsPercent:        f(1.0),
						CurrentProduction:  10,
						CurrentConsumption: 40,
						MaxProd:            100.0,
						MaxConsumed:        200.0,
					},
				})
			})

			It("sets 'item_production_capacity_per_min' to zero", func() {
				collector.Collect()

				val, err := gaugeValue(exporter.ItemProductionCapacityPerMinute, "Iron Rod")
				Expect(err).ToNot(HaveOccurred())
				Expect(val).To(Equal(float64(0)))
			})

			It("sets 'item_consumption_capacity_per_min' to zero", func() {
				collector.Collect()

				val, err := gaugeValue(exporter.ItemConsumptionCapacityPerMinute, "Iron Rod")
				Expect(err).ToNot(HaveOccurred())
				Expect(val).To(Equal(float64(0)))
			})
		})
	})
})
