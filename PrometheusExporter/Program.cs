using System;
using System.Collections.Generic;
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

            try
            {
                Prometheus.Metrics.SuppressDefaultMetrics();
                var promServer = new Prometheus.MetricServer(hostname: "localhost", port: 9000);
                promServer.Start();

                Console.WriteLine("Begun exposing metrics at localhost:9000");

                _source = new CancellationTokenSource();
                CancellationToken token = _source.Token;
                PowerMetricsCollector powerCollector = new PowerMetricsCollector();

                Task.WaitAll(
                    powerCollector.BeginCollecting(token)
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

