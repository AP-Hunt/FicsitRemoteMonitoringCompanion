using System;
using System.Collections.Generic;
using System.IO;
using System.Net;
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

    class Program
    {
        static void Main(string[] args)
        {
            try
            {
                Prometheus.Metrics.SuppressDefaultMetrics();
                var promServer = new Prometheus.MetricServer(hostname: "localhost", port: 9000);
                promServer.Start();

                Console.WriteLine("Begun exposing metrics at localhost:9000");

                BeginMonitoringPower();
                Console.ReadLine();
            }
            catch(Exception ex)
            {
                Console.WriteLine(ex.Message);
            }
        }

        static void BeginMonitoringPower() 
        {
            string powerUrl = "http://localhost:8090/getPowerData";

            double second = 1000;
            System.Timers.Timer timer = new System.Timers.Timer(5 * second);
            timer.Elapsed += ReadPowerMetrics(powerUrl);
            timer.AutoReset = true;
            timer.Enabled = true;

        }

        private static ElapsedEventHandler ReadPowerMetrics(string powerUrl)
        {
            return (Object source, ElapsedEventArgs e) =>
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
                catch(Exception ex)
                {
                    Console.WriteLine(ex.Message);
                }
            };
        }

        private static void UpdateCircuitMetrics(PowerData circuit)
        {
            string circuitId = circuit.CircuitID.ToString();
            Metrics.PowerBatteryCapacity.WithLabels(circuitId).Set(circuit.BatteryCapacity);
            Metrics.PowerBatteryDifferential.WithLabels(circuitId).Set(circuit.BatteryDifferential);
            Metrics.PowerCapacity.WithLabels(circuitId).Set(circuit.PowerCapacity);
            Metrics.PowerConsumed.WithLabels(circuitId).Set(circuit.PowerConsumed);
            Metrics.PowerFuseBlown.WithLabels(circuitId).Set(circuit.FuseBlown? 1d : 0d);
            Metrics.PowerMaxConsumed.WithLabels(circuitId).Set(circuit.PowerMaxConsumed);
            Metrics.PowerProduced.WithLabels(circuitId).Set(circuit.PowerProduced);

            UpdateBatteryPercentMetric(circuit);
            UpdateBatteryTimeEmptyMetric(circuit);
        }

        private static void UpdateBatteryPercentMetric(PowerData circuit)
        {
            // Battery percent can come through as the string "NaN" OR the double value with a percentage sign
            // Need to cast it all to lower case, so that "nan" isn't seen as Not A Number, which is a valid
            // Double type value
            double batteryPct = 0;
            string batteryPctValue = circuit.BatteryPercent.ToLower().Trim(new Char[] { '%' });
            bool success = double.TryParse(batteryPctValue, out batteryPct);

            if (success && !double.IsNaN(batteryPct))
            {
                Metrics.PowerBatteryPercent.WithLabels(circuit.CircuitID.ToString()).Set(batteryPct);
            }
            else
            {
                Metrics.PowerBatteryPercent.WithLabels(circuit.CircuitID.ToString()).Set(0d);
            }
        }

        private static void UpdateBatteryTimeEmptyMetric(PowerData circuit)
        {
            // Battery time empty comes through as a time span string
            // eh hh:mm:ss
            TimeSpan ts = TimeSpan.Zero;
            bool success = TimeSpan.TryParse(circuit.BatteryTimeEmpty, out ts);

            if(success) {
                Metrics.PowerBatteryTimeEmpty.WithLabels(circuit.CircuitID.ToString()).Set(ts.TotalSeconds);
            } else {
                Metrics.PowerBatteryTimeEmpty.WithLabels(circuit.CircuitID.ToString()).Set(0);
            }
        }
    }
}

