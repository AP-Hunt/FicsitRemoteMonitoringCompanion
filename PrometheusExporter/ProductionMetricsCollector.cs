using System;
using System.Collections.Generic;
using System.IO;
using System.Net;
using System.Net.Http;
using System.Text;
using System.Threading;
using System.Threading.Tasks;

namespace PrometheusExporter
{
    class ProductionMetricsCollector : IMetricCollector
    {
        private readonly Uri frmAddress;
        private readonly HttpClient httpClient;

        public static IEnumerable<Prometheus.Collector> ExposedMetrics
        {
            get
            {
                return new List<Prometheus.Collector>
                {
                    ItemProductionCapacityPerMinute,
                    ItemProductionCapacityPercent,
                    ItemConsumptionCapacityPerMinute,
                    ItemConsumptionCapacityPercent,
                    ItemsProduced,
                    ItemsConsumed
                };
            }
        }

        public ProductionMetricsCollector(Uri frmAddress)
        {
            this.frmAddress = frmAddress;
            this.httpClient = new HttpClient();
        }

        public Task BeginCollecting(CancellationToken token)
        {
            string productionUrl = new Uri(this.frmAddress, "/getProdStats").ToString();
            Console.WriteLine($"Will collect production metrics from {productionUrl}");
            return Task.Run(async () =>
            {
                try
                {
                    while (!token.IsCancellationRequested)
                    {
                        await Task.Delay(5 * 1000, token);
                        if (!token.IsCancellationRequested)
                        {
                            await ReadProductionMetrics(productionUrl);
                        }
                    }
                }
                catch (TaskCanceledException)
                {
                }
            }, token);
        }

        private async Task ReadProductionMetrics(string productionUrl)
        {
            try
            {
                string responseJson = await this.httpClient.GetStringAsync(productionUrl);

                List<ProductionDetail> productionDetails = Parsers.ProductionDetailsParser.ParseJSON(responseJson);

                foreach (ProductionDetail detail in productionDetails)
                {
                    UpdateProductionMetrics(detail);

                }
            }
            catch (Exception ex)
            {
                Console.WriteLine(ex.Message);
            }
        }

        private void UpdateProductionMetrics(ProductionDetail detail)
        {
            ItemProductionCapacityPerMinute.WithLabels(detail.ItemName).TrySet(detail.ProductionCapacity);
            ItemProductionCapacityPercent.WithLabels(detail.ItemName).TrySet(detail.ProductionPercent);
            ItemConsumptionCapacityPerMinute.WithLabels(detail.ItemName).TrySet(detail.ConsumptionCapacity);
            ItemConsumptionCapacityPercent.WithLabels(detail.ItemName).TrySet(detail.ConsumptionPercent);
            ItemsProduced.WithLabels(detail.ItemName).TrySet(detail.CurrentProduction);
            ItemsConsumed.WithLabels(detail.ItemName).TrySet(detail.CurrentConsumption);
        }

        private static readonly Prometheus.Gauge ItemProductionCapacityPerMinute = Prometheus.Metrics.CreateGauge(
            "item_production_capacity_per_min", "The factory's capacity for the production of an item, per minute",
            new Prometheus.GaugeConfiguration()
            {
                LabelNames = new string[] {
                    "item_name"
                },
            });

        private static readonly Prometheus.Gauge ItemProductionCapacityPercent = Prometheus.Metrics.CreateGauge(
            "item_production_capacity_pc", "The percentage of an item's production capacity being used",
            new Prometheus.GaugeConfiguration()
            {
                LabelNames = new string[] {
                    "item_name"
                },
            });

        private static readonly Prometheus.Gauge ItemConsumptionCapacityPerMinute = Prometheus.Metrics.CreateGauge(
            "item_consumption_capacity_per_min", "The factory's capacity for the consumption of an item, per minute",
            new Prometheus.GaugeConfiguration()
            {
                LabelNames = new string[] {
                    "item_name"
                },
            });

        private static readonly Prometheus.Gauge ItemConsumptionCapacityPercent = Prometheus.Metrics.CreateGauge(
            "item_consumption_capacity_pc", "The percentage of an item's consumption capacity being used",
            new Prometheus.GaugeConfiguration()
            {
                LabelNames = new string[] {
                    "item_name"
                },
            });

        private static readonly Prometheus.Gauge ItemsProduced = Prometheus.Metrics.CreateGauge(
            "items_produced_per_min", "The number of an item being produced, per minute",
            new Prometheus.GaugeConfiguration()
            {
                LabelNames = new string[] {
                    "item_name"
                },
            });

        private static readonly Prometheus.Gauge ItemsConsumed = Prometheus.Metrics.CreateGauge(
            "items_consumed_per_min", "The number of an item being consumed, per minute",
            new Prometheus.GaugeConfiguration()
            {
                LabelNames = new string[] {
                    "item_name"
                },
            });
    }
}
