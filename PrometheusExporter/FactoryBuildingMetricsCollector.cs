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
    class FactoryBuildingMetricsCollector : IMetricCollector
    {
        private readonly Uri frmAddress;
        private readonly HttpClient httpClient;

        public static IEnumerable<Prometheus.Collector> ExposedMetrics
        {
            get
            {
                return new List<Prometheus.Collector>
                {
                    MachineItemsProducedPerMinute,
                    MachineItemsProducedPerMinute
                };
            }
        }

        public FactoryBuildingMetricsCollector(Uri frmAddress)
        {
            this.frmAddress = frmAddress;
            this.httpClient = new HttpClient();
        }

        public Task BeginCollecting(CancellationToken token)
        {
            string factoryUrl = new Uri(this.frmAddress, "/getFactory").ToString();
            Console.WriteLine($"Will collect factory metrics from {factoryUrl}");
            return Task.Run(async () =>
            {
                try
                {
                    while (!token.IsCancellationRequested)
                    {
                        if (!token.IsCancellationRequested)
                        {
                            await ReadFactoryMetrics(factoryUrl);
                            await Task.Delay(60 * 1000, token);
                        }
                    }
                }
                catch (TaskCanceledException)
                {
                }
            }, token);
        }

        private async Task ReadFactoryMetrics(string productionUrl)
        {
            try
            {
                string responseJson = await this.httpClient.GetStringAsync(productionUrl);

                var options = new System.Text.Json.JsonSerializerOptions
                {
                    AllowTrailingCommas = true,
                    PropertyNameCaseInsensitive = true
                };
                List<FactoryBuildingDetail> factoryDetails = System.Text.Json.JsonSerializer.Deserialize<List<FactoryBuildingDetail>>(responseJson, options);

                foreach (FactoryBuildingDetail buildingDetail in factoryDetails)
                {
                    UpdateMachineMetrics(buildingDetail);
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine(ex.Message);
            }
        }

        private void UpdateMachineMetrics(FactoryBuildingDetail detail)
        {
            foreach(RecipeProductionDetail output in detail.Production)
            {
                MachineItemsProducedPerMinute.WithLabels(ProductionLabelSet(output, detail)).TrySet(output.CurrentProduced);
                MachineItemsProducedEfficiency.WithLabels(ProductionLabelSet(output, detail)).TrySet(double.Parse(output.ProductionPercent));
            }
        }

        private string[] ProductionLabelSet(RecipeProductionDetail output, FactoryBuildingDetail buldingDetail)
        {
            return new string[]{
                output.Name, 
                buldingDetail.Building, 
                buldingDetail.Location.X.ToString(), 
                buldingDetail.Location.Y.ToString(), 
                buldingDetail.Location.Z.ToString()
            };
        }

        private static readonly string[] productionLabels = new string[] {
            "item_name",
            "machine_name",
            "x",
            "y",
            "z"
        };

        private static readonly Prometheus.Gauge MachineItemsProducedPerMinute = Prometheus.Metrics.CreateGauge(
            "machine_items_produced_per_min", "How much of an item a building is producing",
            new Prometheus.GaugeConfiguration()
            {
                LabelNames = productionLabels
            });

        private static readonly Prometheus.Gauge MachineItemsProducedEfficiency = Prometheus.Metrics.CreateGauge(
            "machine_items_produced_pc", "The efficiency with which a building is producing an item",
            new Prometheus.GaugeConfiguration()
            {
                LabelNames = productionLabels
            });
    }
}
