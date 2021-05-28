using System;
using System.Collections.Generic;
using System.IO;
using System.Net;
using System.Threading;
using System.Threading.Tasks;
using System.Timers;

namespace PrometheusExporter
{
    public class PowerData
    {
        public int CircuitID { get; set; }
        public double PowerCapacity { get; set; }
        public double PowerConsumed { get; set; }
        [System.Text.Json.Serialization.JsonPropertyName("PowerProductiuon")]
        public double PowerProduced { get; set; }
        public double PowerMaxConsumed { get; set; }
        public double BatteryDifferential { get; set; }
        public string BatteryPercent { get; set; }
        public double BatteryCapacity { get; set; }
        public string BatteryTimeEmpty { get; set; }
        public bool FuseBlown { get; set; }

    }

    class PowerMetricsCollector : IMetricCollector
    {
        public PowerMetricsCollector()
        {
        }

        public Task BeginCollecting(CancellationToken token)
        {
            string powerUrl = "http://localhost:8090/getPowerData";
            return Task.Run(async () =>
            {
                try
                {
                    while (!token.IsCancellationRequested)
                    {
                        await Task.Delay(5 * 1000, token);
                        if (!token.IsCancellationRequested)
                        {
                            ReadPowerMetrics(powerUrl);
                        }
                    }
                }
                catch(TaskCanceledException)
                {
                }
            }, token);
        }

        private void ReadPowerMetrics(string powerUrl)
        {
            try
            {
                WebRequest req = WebRequest.Create(powerUrl);
                WebResponse resp = req.GetResponse();
                Stream responseStream = resp.GetResponseStream();
                StreamReader rdr = new StreamReader(responseStream);
                string responseJson = rdr.ReadToEnd();
                resp.Close();

                var options = new System.Text.Json.JsonSerializerOptions
                {
                    AllowTrailingCommas = true,
                    PropertyNameCaseInsensitive = true
                };

                List<PowerData> powerCircuits = System.Text.Json.JsonSerializer.Deserialize<List<PowerData>>(responseJson, options);

                foreach (PowerData circuit in powerCircuits)
                {
                    UpdateCircuitMetrics(circuit);

                }
            }
            catch (Exception ex)
            {
                Console.WriteLine(ex.Message);
            }
        }

        private void UpdateCircuitMetrics(PowerData circuit)
        {
            string circuitId = circuit.CircuitID.ToString();
            PowerBatteryCapacity.WithLabels(circuitId).Set(circuit.BatteryCapacity);
            PowerBatteryDifferential.WithLabels(circuitId).Set(circuit.BatteryDifferential);
            PowerCapacity.WithLabels(circuitId).Set(circuit.PowerCapacity);
            PowerConsumed.WithLabels(circuitId).Set(circuit.PowerConsumed);
            PowerFuseBlown.WithLabels(circuitId).Set(circuit.FuseBlown ? 1d : 0d);
            PowerMaxConsumed.WithLabels(circuitId).Set(circuit.PowerMaxConsumed);
            PowerProduced.WithLabels(circuitId).Set(circuit.PowerProduced);

            UpdateBatteryPercentMetric(circuit);
            UpdateBatteryTimeEmptyMetric(circuit);
        }

        private void UpdateBatteryPercentMetric(PowerData circuit)
        {
            // Battery percent can come through as the string "NaN" OR the double value with a percentage sign
            // Need to cast it all to lower case, so that "nan" isn't seen as Not A Number, which is a valid
            // Double type value
            double batteryPct = 0;
            string batteryPctValue = circuit.BatteryPercent.ToLower().Trim(new Char[] { '%' });
            bool success = double.TryParse(batteryPctValue, out batteryPct);

            if (success && !double.IsNaN(batteryPct))
            {
                PowerBatteryPercent.WithLabels(circuit.CircuitID.ToString()).Set(batteryPct);
            }
            else
            {
                PowerBatteryPercent.WithLabels(circuit.CircuitID.ToString()).Set(0d);
            }
        }

        private void UpdateBatteryTimeEmptyMetric(PowerData circuit)
        {
            // Battery time empty comes through as a time span string
            // eh hh:mm:ss
            TimeSpan ts = TimeSpan.Zero;
            bool success = TimeSpan.TryParse(circuit.BatteryTimeEmpty, out ts);

            if (success)
            {
                PowerBatteryTimeEmpty.WithLabels(circuit.CircuitID.ToString()).Set(ts.TotalSeconds);
            }
            else
            {
                PowerBatteryTimeEmpty.WithLabels(circuit.CircuitID.ToString()).Set(0);
            }
        }

        private static readonly Prometheus.Gauge PowerCapacity = Prometheus.Metrics.CreateGauge(
        "power_capacity_mw", "Total capacity of all generators in a circuit",
        new Prometheus.GaugeConfiguration()
        {
            LabelNames = new string[] {
                    "circuit_id"
            },
        });

        private static readonly Prometheus.Gauge PowerConsumed = Prometheus.Metrics.CreateGauge(
            "power_consumed_mw", "Amount of power being consumed in a circuit",
            new Prometheus.GaugeConfiguration()
            {
                LabelNames = new string[] {
                    "circuit_id"
                },
            });

        private static readonly Prometheus.Gauge PowerMaxConsumed = Prometheus.Metrics.CreateGauge(
            "power_max_consumed_mw", "Highest amount of power consumed in a circuit",
            new Prometheus.GaugeConfiguration()
            {
                LabelNames = new string[] {
                    "circuit_id"
                },
            });

        private static readonly Prometheus.Gauge PowerProduced = Prometheus.Metrics.CreateGauge(
            "power_produced_mw", "Amount of power being produced in a circuit",
            new Prometheus.GaugeConfiguration()
            {
                LabelNames = new string[] {
                    "circuit_id"
                },
            });

        private static readonly Prometheus.Gauge PowerBatteryDifferential = Prometheus.Metrics.CreateGauge(
             "power_battery_differential_pc", "Amount of battery capacity as a percentage of total capacity in a circuit",
             new Prometheus.GaugeConfiguration()
             {
                 LabelNames = new string[] {
                    "circuit_id"
                 },
             });

        private static readonly Prometheus.Gauge PowerBatteryPercent = Prometheus.Metrics.CreateGauge(
             "power_battery_pc", "Amount of battery power available as a percentage of total battery capacity in a circuit",
             new Prometheus.GaugeConfiguration()
             {
                 LabelNames = new string[] {
                    "circuit_id"
                 },
             });

        private static readonly Prometheus.Gauge PowerBatteryCapacity = Prometheus.Metrics.CreateGauge(
             "power_battery_capacity_mw", "Total battery capacity in a circuit",
             new Prometheus.GaugeConfiguration()
             {
                 LabelNames = new string[] {
                    "circuit_id"
                 },
             });

        private static readonly Prometheus.Gauge PowerBatteryTimeEmpty = Prometheus.Metrics.CreateGauge(
             "power_battery_time_empty_seconds", "Amount of time batteries have been empty in a circuit",
             new Prometheus.GaugeConfiguration()
             {
                 LabelNames = new string[] {
                    "circuit_id"
                 },
             });

        private static readonly Prometheus.Gauge PowerFuseBlown = Prometheus.Metrics.CreateGauge(
             "power_fuse_blown", "Whether the fuse has been blown in a circuit",
             new Prometheus.GaugeConfiguration()
             {
                 LabelNames = new string[] {
                    "circuit_id"
                 },
             });
    }
}
