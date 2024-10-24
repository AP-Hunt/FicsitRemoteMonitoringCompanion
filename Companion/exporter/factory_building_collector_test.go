package exporter_test

import (
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/exporter"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FactoryBuildingCollector", func() {
	var url = "http://localhost:9080"
	var saveName = "default"
	var collector *exporter.FactoryBuildingCollector

	BeforeEach(func() {
		FRMServer.Reset()
		collector = exporter.NewFactoryBuildingCollector("/getFactory")

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
						ProdPercent: 0.5,
					},
					{
						Name:        "Iron Nothing",
						CurrentProd: 1000.0,
						MaxProd:     4000.0,
						ProdPercent: 0.25,
					},
				},
				Ingredients: []exporter.Ingredient{
					{
						Name:            "Iron Ore",
						CurrentConsumed: 5.0,
						MaxConsumed:     5.0,
						ConsPercent:     1.0,
					},
				},
				ManuSpeed:    100.0,
				IsConfigured: false,
				IsProducing:  false,
				IsPaused:     false,
				CircuitGroupId:    0,
				PowerInfo: exporter.PowerInfo{
					CircuitGroupId:     1,
					PowerConsumed: 23,
				},
			},
		})
	})

	AfterEach(func() {
		collector = nil
	})

	Describe("Factory Power", func() {

		It("Records power per circuit", func() {
			collector.Collect(url, saveName)
			val, err := gaugeValue(exporter.FactoryPower, "1", url, saveName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(23.0))
			val2, _ := gaugeValue(exporter.FactoryPowerMax, "1", url, saveName)
			Expect(val2).To(Equal(exporter.SmelterPower))
		})
	})

	Describe("Machine item production metrics", func() {
		It("records a metric with labels for the produced item name, machine type, and x, y, z coordinates", func() {
			collector.Collect(url, saveName)
			metric, err := getMetric(exporter.MachineItemsProducedPerMin, "Iron Ingot", "Smelter", "100", "200", "-300", url, saveName)
			Expect(err).ToNot(HaveOccurred())
			Expect(metric).ToNot(BeNil())
		})

		It("records the current production figure as the metric value", func() {
			collector.Collect(url, saveName)

			val, err := gaugeValue(exporter.MachineItemsProducedPerMin, "Iron Ingot", "Smelter", "100", "200", "-300", url, saveName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(10.0)))
		})

		Describe("when a machine has multiple outputs", func() {
			It("records a metric per item", func() {
				collector.Collect(url, saveName)

				ironIngots, err := gaugeValue(exporter.MachineItemsProducedPerMin, "Iron Ingot", "Smelter", "100", "200", "-300", url, saveName)
				Expect(err).ToNot(HaveOccurred())
				Expect(ironIngots).To(Equal(float64(10.0)))

				ironNothing, err := gaugeValue(exporter.MachineItemsProducedPerMin, "Iron Nothing", "Smelter", "100", "200", "-300", url, saveName)
				Expect(err).ToNot(HaveOccurred())
				Expect(ironNothing).To(Equal(float64(1000.0)))
			})
		})
	})

	Describe("Machine item production efficiency metrics", func() {
		It("records a metric with labels for the produced item name, machine type, and x, y, z coordinates", func() {
			collector.Collect(url, saveName)
			metric, err := getMetric(exporter.MachineItemsProducedEffiency, "Iron Ingot", "Smelter", "100", "200", "-300", url, saveName)
			Expect(err).ToNot(HaveOccurred())
			Expect(metric).ToNot(BeNil())
		})

		It("records the current production efficiency as the metric value", func() {
			collector.Collect(url, saveName)

			val, err := gaugeValue(exporter.MachineItemsProducedEffiency, "Iron Ingot", "Smelter", "100", "200", "-300", url, saveName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(0.5)))
		})

		Describe("when a machine has multiple outputs", func() {
			It("records a metric per item", func() {
				collector.Collect(url, saveName)

				ironIngots, err := gaugeValue(exporter.MachineItemsProducedEffiency, "Iron Ingot", "Smelter", "100", "200", "-300", url, saveName)
				Expect(err).ToNot(HaveOccurred())
				Expect(ironIngots).To(Equal(float64(0.5)))

				ironNothing, err := gaugeValue(exporter.MachineItemsProducedEffiency, "Iron Nothing", "Smelter", "100", "200", "-300", url, saveName)
				Expect(err).ToNot(HaveOccurred())
				Expect(ironNothing).To(Equal(float64(0.25)))
			})
		})
	})
})
