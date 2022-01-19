package exporter_test

import (
	"context"

	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/m/v2/exporter"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FactoryBuildingCollector", func() {
	var collector *exporter.FactoryBuildingCollector

	BeforeEach(func() {
		FRMServer.Reset()
		collector = exporter.NewFactoryBuildingCollector(context.Background(), "http://localhost:9080/getFactory")

		FRMServer.ReturnsFactoryBuildings([]exporter.BuildingDetail{
			{
				Building: "Smelter",
				Location: exporter.Location{
					X:        100.0,
					Y:        200.0,
					Z:        -300.0,
					Rotation: 60,
				},
				Recipe: "Iron Ingot",
				Production: []exporter.Production{
					{
						Name:        "Iron Ingot",
						CurrentProd: 10.0,
						MaxProd:     10.0,
						ProdPercent: "0.5",
					},
					{
						Name:        "Iron Nothing",
						CurrentProd: 1000.0,
						MaxProd:     4000.0,
						ProdPercent: "0.25",
					},
				},
				Ingredients: []exporter.Ingredient{
					{
						Name:            "Iron Ore",
						CurrentConsumed: 5.0,
						MaxConsumed:     5.0,
						ConsPercent:     "1.0",
					},
				},
				ManuSpeed:    1.0,
				IsConfigured: false,
				IsProducing:  false,
				IsPaused:     false,
				CircuitID:    0,
			},
		})
	})

	AfterEach(func() {
		collector = nil
	})

	Describe("Machine item production metrics", func() {
		It("records a metric with labels for the produced item name, machine type, and x, y, z coordinates", func() {
			collector.Collect()
			metric, err := getMetric(exporter.MachineItemsProducedPerMin, "Iron Ingot", "Smelter", "100.00000", "200.00000", "-300.00000")
			Expect(err).ToNot(HaveOccurred())
			Expect(metric).ToNot(BeNil())
		})

		It("records the current production figure as the metric value", func() {
			collector.Collect()

			val, err := gaugeValue(exporter.MachineItemsProducedPerMin, "Iron Ingot", "Smelter", "100.00000", "200.00000", "-300.00000")
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(10.0)))
		})

		Describe("when a machine has multiple outputs", func() {
			It("records a metric per item", func() {
				collector.Collect()

				ironIngots, err := gaugeValue(exporter.MachineItemsProducedPerMin, "Iron Ingot", "Smelter", "100.00000", "200.00000", "-300.00000")
				Expect(err).ToNot(HaveOccurred())
				Expect(ironIngots).To(Equal(float64(10.0)))

				ironNothing, err := gaugeValue(exporter.MachineItemsProducedPerMin, "Iron Nothing", "Smelter", "100.00000", "200.00000", "-300.00000")
				Expect(err).ToNot(HaveOccurred())
				Expect(ironNothing).To(Equal(float64(1000.0)))
			})
		})
	})

	Describe("Machine item production efficiency metrics", func() {
		It("records a metric with labels for the produced item name, machine type, and x, y, z coordinates", func() {
			collector.Collect()
			metric, err := getMetric(exporter.MachineItemsProducedEffiency, "Iron Ingot", "Smelter", "100.00000", "200.00000", "-300.00000")
			Expect(err).ToNot(HaveOccurred())
			Expect(metric).ToNot(BeNil())
		})

		It("records the current production efficiency as the metric value", func() {
			collector.Collect()

			val, err := gaugeValue(exporter.MachineItemsProducedEffiency, "Iron Ingot", "Smelter", "100.00000", "200.00000", "-300.00000")
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(0.5)))
		})

		Describe("when a machine has multiple outputs", func() {
			It("records a metric per item", func() {
				collector.Collect()

				ironIngots, err := gaugeValue(exporter.MachineItemsProducedEffiency, "Iron Ingot", "Smelter", "100.00000", "200.00000", "-300.00000")
				Expect(err).ToNot(HaveOccurred())
				Expect(ironIngots).To(Equal(float64(0.5)))

				ironNothing, err := gaugeValue(exporter.MachineItemsProducedEffiency, "Iron Nothing", "Smelter", "100.00000", "200.00000", "-300.00000")
				Expect(err).ToNot(HaveOccurred())
				Expect(ironNothing).To(Equal(float64(0.25)))
			})
		})
	})
})
