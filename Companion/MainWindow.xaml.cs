using Microsoft.Web.WebView2.Wpf;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Navigation;
using System.Windows.Shapes;

namespace Companion
{
    /// <summary>
    /// Interaction logic for MainWindow.xaml
    /// </summary>
    public partial class MainWindow : Window
    {
        Label loadingLabel;
        public MainWindow()
        {
            InitializeComponent();
            loadingLabel = new Label()
            {
                Content = "Loading Grafana...",
                HorizontalAlignment = HorizontalAlignment.Center,
                VerticalAlignment = VerticalAlignment.Center
            };
            this.MainPanel.Children.Add(loadingLabel);
        }

        public void ShowGrafana()
        {
            Application.Current.Dispatcher.Invoke(() =>
            {
                this.MainPanel.Children.Remove(loadingLabel);
                this.MainPanel.Children.Add(new WebView2()
                {
                    Source = new Uri("http://localhost:3000")
                });
            });
        }

        private async void Window_ContentRendered(object sender, EventArgs e)
        {
            Config.ConfigWindow configWindow = new Config.ConfigWindow();
            configWindow.Owner = this;
            configWindow.ShowDialog();

            Config.ConfigFile cfg = configWindow.Config;

            PrometheusExporterHost.Start();
            PrometheusHost.Start();
            GrafanaHost.Start();

            await GrafanaHost.WaitForReadiness();
            this.ShowGrafana();
        }
    }
}
