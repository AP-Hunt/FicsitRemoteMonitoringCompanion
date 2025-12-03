package exporter_test

import (
	"github.com/AP-Hunt/FicsitRemoteMonitoringCompanion/Companion/exporter"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("InventoryCollectors", func() {
	var url string
	var sessionName = "default"

	BeforeEach(func() {
		FRMServer.Reset()
		url = FRMServer.server.URL
	})

	Describe("CloudInventoryCollector", func() {
		var collector *exporter.CloudInventoryCollector

		BeforeEach(func() {
			collector = exporter.NewCloudInventoryCollector("/getCloudInventory")
		})

		AfterEach(func() {
			collector = nil
		})

		Describe("Cloud inventory metrics", func() {
			BeforeEach(func() {
				FRMServer.ReturnsCloudInventoryData([]exporter.InventoryItem{
					{
						Name:      "Iron Ingot",
						Amount:    500,
						MaxAmount: 1000,
					},
					{
						Name:      "Copper Ingot",
						Amount:    250,
						MaxAmount: 500,
					},
				})
			})

			It("records metrics with labels for item name", func() {
				collector.Collect(url, sessionName)
				metric, err := getMetric(exporter.CloudInventory, "Iron Ingot", url, sessionName)
				Expect(err).ToNot(HaveOccurred())
				Expect(metric).ToNot(BeNil())
			})

			It("records the current amount as the metric value", func() {
				collector.Collect(url, sessionName)

				val, err := gaugeValue(exporter.CloudInventory, "Iron Ingot", url, sessionName)
				Expect(err).ToNot(HaveOccurred())
				Expect(val).To(Equal(500.0))
			})

			It("records the max amount metric", func() {
				collector.Collect(url, sessionName)

				val, err := gaugeValue(exporter.CloudInventoryMax, "Iron Ingot", url, sessionName)
				Expect(err).ToNot(HaveOccurred())
				Expect(val).To(Equal(1000.0))
			})

			Describe("when there are multiple items", func() {
				It("records a metric per item", func() {
					collector.Collect(url, sessionName)

					ironVal, err := gaugeValue(exporter.CloudInventory, "Iron Ingot", url, sessionName)
					Expect(err).ToNot(HaveOccurred())
					Expect(ironVal).To(Equal(500.0))

					copperVal, err := gaugeValue(exporter.CloudInventory, "Copper Ingot", url, sessionName)
					Expect(err).ToNot(HaveOccurred())
					Expect(copperVal).To(Equal(250.0))
				})
			})
		})
	})

	Describe("WorldInventoryCollector", func() {
		var collector *exporter.WorldInventoryCollector

		BeforeEach(func() {
			collector = exporter.NewWorldInventoryCollector("/getWorldInventory")
		})

		AfterEach(func() {
			collector = nil
		})

		Describe("World inventory metrics", func() {
			BeforeEach(func() {
				FRMServer.ReturnsWorldInventoryData([]exporter.InventoryItem{
					{
						Name:      "Concrete",
						Amount:    10000,
						MaxAmount: 50000,
					},
					{
						Name:      "Steel Beam",
						Amount:    2500,
						MaxAmount: 10000,
					},
				})
			})

			It("records metrics with labels for item name", func() {
				collector.Collect(url, sessionName)
				metric, err := getMetric(exporter.WorldInventory, "Concrete", url, sessionName)
				Expect(err).ToNot(HaveOccurred())
				Expect(metric).ToNot(BeNil())
			})

			It("records the current amount as the metric value", func() {
				collector.Collect(url, sessionName)

				val, err := gaugeValue(exporter.WorldInventory, "Concrete", url, sessionName)
				Expect(err).ToNot(HaveOccurred())
				Expect(val).To(Equal(10000.0))
			})

			It("records the max amount metric", func() {
				collector.Collect(url, sessionName)

				val, err := gaugeValue(exporter.WorldInventoryMax, "Concrete", url, sessionName)
				Expect(err).ToNot(HaveOccurred())
				Expect(val).To(Equal(50000.0))
			})

			Describe("when there are multiple items", func() {
				It("records a metric per item", func() {
					collector.Collect(url, sessionName)

					concreteVal, err := gaugeValue(exporter.WorldInventory, "Concrete", url, sessionName)
					Expect(err).ToNot(HaveOccurred())
					Expect(concreteVal).To(Equal(10000.0))

					steelVal, err := gaugeValue(exporter.WorldInventory, "Steel Beam", url, sessionName)
					Expect(err).ToNot(HaveOccurred())
					Expect(steelVal).To(Equal(2500.0))
				})
			})
		})
	})

	Describe("StorageInventoryCollector", func() {
		var collector *exporter.StorageInventoryCollector

		BeforeEach(func() {
			collector = exporter.NewStorageInventoryCollector("/getStorage")
		})

		AfterEach(func() {
			collector = nil
		})

		Describe("Storage inventory metrics", func() {
			BeforeEach(func() {
				FRMServer.ReturnsStorageContainerData([]exporter.ContainerDetail{
					{
						Name: "Storage Container",
						Location: exporter.Location{
							X: 150.0,
							Y: 250.0,
							Z: -350.0,
						},
						Inventory: []exporter.InventoryItem{
							{
								Name:      "Iron Ore",
								Amount:    1200,
								MaxAmount: 2400,
							},
							{
								Name:      "Copper Ore",
								Amount:    600,
								MaxAmount: 2400,
							},
						},
					},
				})
			})

			It("records metrics with labels for item name, container name, and coordinates", func() {
				collector.Collect(url, sessionName)
				metric, err := getMetric(exporter.StorageInventory, "Iron Ore", "Storage Container", "150", "250", "-350", url, sessionName)
				Expect(err).ToNot(HaveOccurred())
				Expect(metric).ToNot(BeNil())
			})

			It("records the current amount as the metric value", func() {
				collector.Collect(url, sessionName)

				val, err := gaugeValue(exporter.StorageInventory, "Iron Ore", "Storage Container", "150", "250", "-350", url, sessionName)
				Expect(err).ToNot(HaveOccurred())
				Expect(val).To(Equal(1200.0))
			})

			It("records the max amount metric", func() {
				collector.Collect(url, sessionName)

				val, err := gaugeValue(exporter.StorageInventoryMax, "Iron Ore", "Storage Container", "150", "250", "-350", url, sessionName)
				Expect(err).ToNot(HaveOccurred())
				Expect(val).To(Equal(2400.0))
			})

			Describe("when a container has multiple items", func() {
				It("records a metric per item", func() {
					collector.Collect(url, sessionName)

					ironVal, err := gaugeValue(exporter.StorageInventory, "Iron Ore", "Storage Container", "150", "250", "-350", url, sessionName)
					Expect(err).ToNot(HaveOccurred())
					Expect(ironVal).To(Equal(1200.0))

					copperVal, err := gaugeValue(exporter.StorageInventory, "Copper Ore", "Storage Container", "150", "250", "-350", url, sessionName)
					Expect(err).ToNot(HaveOccurred())
					Expect(copperVal).To(Equal(600.0))
				})
			})

			Describe("when there are multiple containers", func() {
				BeforeEach(func() {
					FRMServer.ReturnsStorageContainerData([]exporter.ContainerDetail{
						{
							Name: "Storage Container",
							Location: exporter.Location{
								X: 100.0,
								Y: 200.0,
								Z: -300.0,
							},
							Inventory: []exporter.InventoryItem{
								{
									Name:      "Iron Plate",
									Amount:    400,
									MaxAmount: 800,
								},
							},
						},
						{
							Name: "Storage Container",
							Location: exporter.Location{
								X: 150.0,
								Y: 250.0,
								Z: -350.0,
							},
							Inventory: []exporter.InventoryItem{
								{
									Name:      "Iron Ore",
									Amount:    1200,
									MaxAmount: 2400,
								},
							},
						},
					})
				})

				It("records a metric per container with distinct coordinates", func() {
					collector.Collect(url, sessionName)

					val1, err := gaugeValue(exporter.StorageInventory, "Iron Plate", "Storage Container", "100", "200", "-300", url, sessionName)
					Expect(err).ToNot(HaveOccurred())
					Expect(val1).To(Equal(400.0))

					val2, err := gaugeValue(exporter.StorageInventory, "Iron Ore", "Storage Container", "150", "250", "-350", url, sessionName)
					Expect(err).ToNot(HaveOccurred())
					Expect(val2).To(Equal(1200.0))
				})
			})
		})
	})

	Describe("CrateInventoryCollector", func() {
		var collector *exporter.CrateInventoryCollector

		BeforeEach(func() {
			collector = exporter.NewCrateInventoryCollector("/getCrates")
		})

		AfterEach(func() {
			collector = nil
		})

		Describe("Crate inventory metrics", func() {
			BeforeEach(func() {
				FRMServer.ReturnsCrateData([]exporter.ContainerDetail{
					{
						Name: "Death Crate",
						Location: exporter.Location{
							X: 75.0,
							Y: 125.0,
							Z: -175.0,
						},
						Inventory: []exporter.InventoryItem{
							{
								Name:      "Rifle Ammo",
								Amount:    50,
								MaxAmount: 100,
							},
							{
								Name:      "Health Inhaler",
								Amount:    5,
								MaxAmount: 10,
							},
						},
					},
				})
			})

			It("records metrics with labels for item name, container name, and coordinates", func() {
				collector.Collect(url, sessionName)
				metric, err := getMetric(exporter.CrateInventory, "Rifle Ammo", "Death Crate", "75", "125", "-175", url, sessionName)
				Expect(err).ToNot(HaveOccurred())
				Expect(metric).ToNot(BeNil())
			})

			It("records the current amount as the metric value", func() {
				collector.Collect(url, sessionName)

				val, err := gaugeValue(exporter.CrateInventory, "Rifle Ammo", "Death Crate", "75", "125", "-175", url, sessionName)
				Expect(err).ToNot(HaveOccurred())
				Expect(val).To(Equal(50.0))
			})

			It("records the max amount metric", func() {
				collector.Collect(url, sessionName)

				val, err := gaugeValue(exporter.CrateInventoryMax, "Rifle Ammo", "Death Crate", "75", "125", "-175", url, sessionName)
				Expect(err).ToNot(HaveOccurred())
				Expect(val).To(Equal(100.0))
			})

			Describe("when a crate has multiple items", func() {
				It("records a metric per item", func() {
					collector.Collect(url, sessionName)

					ammoVal, err := gaugeValue(exporter.CrateInventory, "Rifle Ammo", "Death Crate", "75", "125", "-175", url, sessionName)
					Expect(err).ToNot(HaveOccurred())
					Expect(ammoVal).To(Equal(50.0))

					inhalerVal, err := gaugeValue(exporter.CrateInventory, "Health Inhaler", "Death Crate", "75", "125", "-175", url, sessionName)
					Expect(err).ToNot(HaveOccurred())
					Expect(inhalerVal).To(Equal(5.0))
				})
			})

			Describe("when there are multiple crates", func() {
				BeforeEach(func() {
					FRMServer.ReturnsCrateData([]exporter.ContainerDetail{
						{
							Name: "Death Crate",
							Location: exporter.Location{
								X: 75.0,
								Y: 125.0,
								Z: -175.0,
							},
							Inventory: []exporter.InventoryItem{
								{
									Name:      "Rifle Ammo",
									Amount:    50,
									MaxAmount: 100,
								},
							},
						},
						{
							Name: "Dismantle Crate",
							Location: exporter.Location{
								X: 200.0,
								Y: 300.0,
								Z: -400.0,
							},
							Inventory: []exporter.InventoryItem{
								{
									Name:      "Modular Frame",
									Amount:    10,
									MaxAmount: 20,
								},
							},
						},
					})
				})

				It("records a metric per crate with distinct coordinates", func() {
					collector.Collect(url, sessionName)

					val1, err := gaugeValue(exporter.CrateInventory, "Rifle Ammo", "Death Crate", "75", "125", "-175", url, sessionName)
					Expect(err).ToNot(HaveOccurred())
					Expect(val1).To(Equal(50.0))

					val2, err := gaugeValue(exporter.CrateInventory, "Modular Frame", "Dismantle Crate", "200", "300", "-400", url, sessionName)
					Expect(err).ToNot(HaveOccurred())
					Expect(val2).To(Equal(10.0))
				})
			})
		})
	})
})
