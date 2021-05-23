using System;
using System.Collections.Generic;
using System.Configuration;
using System.Data;
using System.Linq;
using System.Threading.Tasks;
using System.Windows;

namespace Companion
{
    /// <summary>
    /// Interaction logic for App.xaml
    /// </summary>
    public partial class App : Application
    {
        protected override void OnStartup(StartupEventArgs e)
        {
            base.OnStartup(e);

            PrometheusExporterHost.Start();
            PrometheusHost.Start();
            GrafanaHost.Start();
        }

        protected override void OnExit(ExitEventArgs e)
        {
            base.OnExit(e);

            GrafanaHost.Stop();
            PrometheusHost.Stop();
            PrometheusExporterHost.Stop();
        }
    }
}
