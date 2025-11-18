package exporter_test

import (
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/exporter"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"math"
)

var _ = Describe("FactoryBuildingCollector", func() {
	var url string
	var sessionName = "default"
	var collector *exporter.FactoryBuildingCollector

	BeforeEach(func() {
		FRMServer.Reset()
		url = FRMServer.server.URL
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
				ManuSpeed:      100.0,
				IsConfigured:   false,
				IsProducing:    false,
				IsPaused:       false,
				CircuitGroupId: 0,
				PowerInfo: exporter.PowerInfo{
					CircuitGroupId:   1,
					PowerConsumed:    23,
					MaxPowerConsumed: 4,
				},
			},
		})
	})

	AfterEach(func() {
		collector = nil
	})

	Describe("Factory Power", func() {

		It("Records power per circuit", func() {
			collector.Collect(url, sessionName)
			val, err := gaugeValue(exporter.FactoryPower, "1", url, sessionName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(23.0))
			val2, _ := gaugeValue(exporter.FactoryPowerMax, "1", url, sessionName)
			Expect(val2).To(Equal(4.0))
		})
	})

	Describe("Machine item production metrics", func() {
		It("records a metric with labels for the produced item name, machine type, and x, y, z coordinates", func() {
			collector.Collect(url, sessionName)
			metric, err := getMetric(exporter.MachineItemsProducedPerMin, "Iron Ingot", "Smelter", "100", "200", "-300", url, sessionName)
			Expect(err).ToNot(HaveOccurred())
			Expect(metric).ToNot(BeNil())
		})

		It("records the current production figure as the metric value", func() {
			collector.Collect(url, sessionName)

			val, err := gaugeValue(exporter.MachineItemsProducedPerMin, "Iron Ingot", "Smelter", "100", "200", "-300", url, sessionName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(10.0)))
		})

		Describe("when a machine has multiple outputs", func() {
			It("records a metric per item", func() {
				collector.Collect(url, sessionName)

				ironIngots, err := gaugeValue(exporter.MachineItemsProducedPerMin, "Iron Ingot", "Smelter", "100", "200", "-300", url, sessionName)
				Expect(err).ToNot(HaveOccurred())
				Expect(ironIngots).To(Equal(float64(10.0)))

				ironNothing, err := gaugeValue(exporter.MachineItemsProducedPerMin, "Iron Nothing", "Smelter", "100", "200", "-300", url, sessionName)
				Expect(err).ToNot(HaveOccurred())
				Expect(ironNothing).To(Equal(float64(1000.0)))
			})
		})
	})

	Describe("Machine item production efficiency metrics", func() {
		It("records a metric with labels for the produced item name, machine type, and x, y, z coordinates", func() {
			collector.Collect(url, sessionName)
			metric, err := getMetric(exporter.MachineItemsProducedEffiency, "Iron Ingot", "Smelter", "100", "200", "-300", url, sessionName)
			Expect(err).ToNot(HaveOccurred())
			Expect(metric).ToNot(BeNil())
		})

		It("records the current production efficiency as the metric value", func() {
			collector.Collect(url, sessionName)

			val, err := gaugeValue(exporter.MachineItemsProducedEffiency, "Iron Ingot", "Smelter", "100", "200", "-300", url, sessionName)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(float64(0.5)))
		})

		Describe("when a machine has multiple outputs", func() {
			It("records a metric per item", func() {
				collector.Collect(url, sessionName)

				ironIngots, err := gaugeValue(exporter.MachineItemsProducedEffiency, "Iron Ingot", "Smelter", "100", "200", "-300", url, sessionName)
				Expect(err).ToNot(HaveOccurred())
				Expect(ironIngots).To(Equal(float64(0.5)))

				ironNothing, err := gaugeValue(exporter.MachineItemsProducedEffiency, "Iron Nothing", "Smelter", "100", "200", "-300", url, sessionName)
				Expect(err).ToNot(HaveOccurred())
				Expect(ironNothing).To(Equal(float64(0.25)))
			})
		})

		Describe("with power particle accelerator making diamonds", func() {
			BeforeEach(func() {
				FRMServer.Reset()
				url = FRMServer.server.URL
				collector = exporter.NewFactoryBuildingCollector("/getFactory")

				FRMServer.ReturnsFactoryBuildings([]exporter.BuildingDetail{
					{
						Building:       "Particle Accelerator",
						Recipe:         "Diamonds",
						ManuSpeed:      100.0,
						CircuitGroupId: 0,
						PowerInfo: exporter.PowerInfo{
							CircuitGroupId: 1,
							PowerConsumed:  23,
							// value will be ignored for this - the recipe here is set to 750 in power_info.go
							MaxPowerConsumed: 4,
						},
					},
				})
			})
			It("recalculates max power use", func() {
				collector.Collect(url, sessionName)

				val, err := gaugeValue(exporter.FactoryPowerMax, "1", url, sessionName)
				Expect(err).ToNot(HaveOccurred())
				Expect(val).To(Equal(750.0))
			})
		})
		Describe("with an overclocked accelerator", func() {
			BeforeEach(func() {
				FRMServer.Reset()
				url = FRMServer.server.URL
				collector = exporter.NewFactoryBuildingCollector("/getFactory")

				FRMServer.ReturnsFactoryBuildings([]exporter.BuildingDetail{
					{
						Building:       "Particle Accelerator",
						Recipe:         "Nuclear Pasta",
						ManuSpeed:      250.0,
						CircuitGroupId: 0,
						PowerInfo: exporter.PowerInfo{
							CircuitGroupId: 1,
							PowerConsumed:  23,
							// value will be ignored for this - the recipe here is set to 750 in power_info.go
							MaxPowerConsumed: 4,
						},
					},
				})
			})
			It("recalculates max power use", func() {
				collector.Collect(url, sessionName)

				val, err := gaugeValue(exporter.FactoryPowerMax, "1", url, sessionName)
				Expect(err).ToNot(HaveOccurred())
				Expect(val).To(Equal(math.Pow((250.0/100), exporter.ClockspeedExponent) * 1500.0))
			})
		})
		Describe("with an underclocked converter", func() {
			BeforeEach(func() {
				FRMServer.Reset()
				url = FRMServer.server.URL
				collector = exporter.NewFactoryBuildingCollector("/getFactory")

				FRMServer.ReturnsFactoryBuildings([]exporter.BuildingDetail{
					{
						Building:       "Converter",
						Recipe:         "Coal",
						ManuSpeed:      50.0,
						CircuitGroupId: 0,
						PowerInfo: exporter.PowerInfo{
							CircuitGroupId: 1,
							PowerConsumed:  23,
							// value will be ignored for this - the recipe here is set to 400 * clockspeed in power_info.go
							MaxPowerConsumed: 4,
						},
					},
				})
			})
			It("recalculates max power use", func() {
				collector.Collect(url, sessionName)

				val, err := gaugeValue(exporter.FactoryPowerMax, "1", url, sessionName)
				Expect(err).ToNot(HaveOccurred())
				Expect(val).To(Equal(math.Pow((50.0/100), exporter.ClockspeedExponent) * 400.0))
			})
		})
		Describe("with an overclocked quantum encoder", func() {
			BeforeEach(func() {
				FRMServer.Reset()
				url = FRMServer.server.URL
				collector = exporter.NewFactoryBuildingCollector("/getFactory")

				FRMServer.ReturnsFactoryBuildings([]exporter.BuildingDetail{
					{
						Building:       "Quantum Encoder",
						Recipe:         "Power Shard",
						ManuSpeed:      250.0,
						CircuitGroupId: 0,
						PowerInfo: exporter.PowerInfo{
							CircuitGroupId: 1,
							PowerConsumed:  23,
							// value will be ignored for this - the recipe here is set to 2000 * clockspeed in power_info.go
							MaxPowerConsumed: 4,
						},
					},
				})
			})
			It("recalculates max power use", func() {
				collector.Collect(url, sessionName)

				val, err := gaugeValue(exporter.FactoryPowerMax, "1", url, sessionName)
				Expect(err).ToNot(HaveOccurred())
				Expect(val).To(Equal(math.Pow((250.0/100), exporter.ClockspeedExponent) * 2000.0))
			})
		})
		Describe("with a somerslooped quantum encoder", func() {
			BeforeEach(func() {
				FRMServer.Reset()
				url = FRMServer.server.URL
				collector = exporter.NewFactoryBuildingCollector("/getFactory")

				FRMServer.ReturnsFactoryBuildings([]exporter.BuildingDetail{
					{
						Building:       "Quantum Encoder",
						Recipe:         "Power Shard",
						ManuSpeed:      100.0,
						CircuitGroupId: 0,
						Somersloops:    4,
						PowerInfo: exporter.PowerInfo{
							CircuitGroupId: 1,
							PowerConsumed:  23,
							// value will be ignored for this - the recipe here is set to 2000 * clockspeed in power_info.go
							MaxPowerConsumed: 4,
						},
					},
				})
			})
			It("recalculates max power use", func() {
				collector.Collect(url, sessionName)

				val, err := gaugeValue(exporter.FactoryPowerMax, "1", url, sessionName)
				Expect(err).ToNot(HaveOccurred())
				Expect(val).To(Equal(4 * 2000.0))
			})
		})
	})
})
