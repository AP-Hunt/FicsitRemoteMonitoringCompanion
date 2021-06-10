using System;
using System.Collections.Generic;
using System.Text;

namespace PrometheusExporter
{
    static class PrometheusExtensions
    {
        public static void TrySet(this Prometheus.Gauge.Child gauge, double? value)
        {
            if(value != null)
            {
                gauge.Set(value.Value);
            }
        }
    }
}
