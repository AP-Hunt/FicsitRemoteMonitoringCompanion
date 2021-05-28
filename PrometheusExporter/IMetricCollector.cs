using System;
using System.Collections.Generic;
using System.Text;
using System.Threading;
using System.Threading.Tasks;

namespace PrometheusExporter
{
    interface IMetricCollector
    {
        Task BeginCollecting(CancellationToken token); 
    }
}
