using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Net;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using System.Timers;

namespace PrometheusExporter
{
    class Program
    {
        private static CancellationTokenSource _source;

        static void Main(string[] args)
        {
            Console.CancelKeyPress += Console_CancelKeyPress;

            if(args.Length < 1){
                Console.Error.WriteLine("First argument must be the address on which Ficsit Remote Monitoring is listening");
                Environment.Exit(1);
            }

            // Output the names of all the metrics exposed by this exporter
            // in HTML format, as a helper for the README
            if(args[0].ToUpper() == "-SHOWMETRICS")
            {
                ShowExportedMetrics();
                return;
            }

            string satisfactoryAddress = args[0];
            Uri satisfactoryUri;
            if(!Uri.TryCreate(satisfactoryAddress, UriKind.Absolute, out satisfactoryUri)){
                Console.Error.WriteLine($"'{satisfactoryAddress}' is not a valid URI");
                Environment.Exit(2);
            }

            Console.WriteLine($"Will contact Ficsit Remote Montioring on {satisfactoryUri.ToString()}");

            try
            {
                Prometheus.Metrics.SuppressDefaultMetrics();
                var promServer = new Prometheus.MetricServer(hostname: "localhost", port: 9000);
                promServer.Start();

                Console.WriteLine("Begun exposing metrics at localhost:9000");

                _source = new CancellationTokenSource();
                CancellationToken token = _source.Token;

                PowerMetricsCollector powerCollector = new PowerMetricsCollector(satisfactoryUri);
                ProductionMetricsCollector productionCollector = new ProductionMetricsCollector(satisfactoryUri);
                FactoryBuildingMetricsCollector factoryCollector = new FactoryBuildingMetricsCollector(satisfactoryUri);

                Task.WaitAll(
                    powerCollector.BeginCollecting(token),
                    productionCollector.BeginCollecting(token),
                    factoryCollector.BeginCollecting(token)
                );
                Console.WriteLine("Exiting");
            }
            catch (Exception ex)
            {
                Console.WriteLine(ex.Message);
            }
        }

        private static void Console_CancelKeyPress(object sender, ConsoleCancelEventArgs e)
        {
            if (_source == null)
            {
                return;
            }

            _source.Cancel();
        }

        private static void ShowExportedMetrics()
        {
            IEnumerable<Type> metricCollectorTypes = typeof(Program)
                .Assembly
                .GetTypes()
                .Where(t => t.IsClass)
                .Where(t => t.GetInterfaces().Contains(typeof(IMetricCollector)));

            var collectors = new List<Prometheus.Collector>();
            foreach (Type collectorType in metricCollectorTypes)
            {
                var exported = (IEnumerable<Prometheus.Collector>)collectorType
                    .GetProperty(nameof(IMetricCollector.ExposedMetrics), System.Reflection.BindingFlags.Static | System.Reflection.BindingFlags.Public)
                    .GetValue(null);
                collectors.AddRange(exported);
            }

            StringBuilder builder = new StringBuilder();
            builder.Append(@"
<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Description</th>
            <th>Labels</th>
        </tr>
    </thead>
    <tbody>");

            foreach(Prometheus.Collector collector in collectors)
            {
                builder.Append($@"
        <tr>
            <td>{collector.Name}</td>
            <td>{collector.Help}</td>
            <td>{string.Join(", ", collector.LabelNames)}</td>
        </tr>");
            }

            builder.Append(@"
    </tbody>
</table>");

            Console.WriteLine(builder.ToString());
        }

    }
}

