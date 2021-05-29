using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Net;
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

                Task.WaitAll(
                    powerCollector.BeginCollecting(token),
                    productionCollector.BeginCollecting(token)
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
    }
}

