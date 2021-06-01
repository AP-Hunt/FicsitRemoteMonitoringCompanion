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
        protected override async void OnStartup(StartupEventArgs e)
        {
            base.OnStartup(e);
        }

        protected override void OnExit(ExitEventArgs e)
        {
            base.OnExit(e);

            GrafanaHost.Stop();
            PrometheusHost.Stop();
            PrometheusExporterHost.Stop();
        }

        private async void Application_Startup(object sender, StartupEventArgs e)
        {
            var mainWindow = new MainWindow();
            Current.ShutdownMode = ShutdownMode.OnMainWindowClose;
            Current.MainWindow = mainWindow;

            Config.ConfigWindow configWindow = new Config.ConfigWindow();

            if (configWindow.ShowDialog() == true)
            {
                Config.ConfigFile cfg = configWindow.Config;
                Config.FicsitRemoteMonitoringConfig ficsitConfig = await Config.ConfigIO.ReadFRMConfigFile(cfg.SatisfactoryGameDirectory);

                PrometheusExporterHost.Start(ficsitConfig.ListenAddress);
                PrometheusHost.Start();
                GrafanaHost.Start();
                mainWindow.Show();
            }
        }
    }
}
