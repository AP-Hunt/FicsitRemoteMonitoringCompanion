using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Shapes;

namespace Companion
{
    /// <summary>
    /// Interaction logic for LogLines.xaml
    /// </summary>
    public partial class LogLines : Window
    {
        BackgroundWorker exporterWorker;
        BackgroundWorker prometheusWorker;
        BackgroundWorker grafanaWorker;
        BackgroundWorker companionWorker;

        public LogLines()
        {
            InitializeComponent();
        }

        private void Window_Loaded(object sender, RoutedEventArgs e)
        {
            exporterWorker = new BackgroundWorker
            {
                WorkerSupportsCancellation = true
            };
            exporterWorker.DoWork += ExporterWorker_Work;
            exporterWorker.RunWorkerAsync();

            prometheusWorker = new BackgroundWorker
            {
                WorkerSupportsCancellation = true
            };
            prometheusWorker.DoWork += PrometheusWorker_Work;
            prometheusWorker.RunWorkerAsync();

            grafanaWorker = new BackgroundWorker
            {
                WorkerSupportsCancellation = true
            };
            grafanaWorker.DoWork += GrafanaWorker_Work;
            grafanaWorker.RunWorkerAsync();

            companionWorker = new BackgroundWorker
            {
                WorkerSupportsCancellation = true
            };
            companionWorker.DoWork += CompanionWorker_Work;
            companionWorker.RunWorkerAsync();
        }

        private void ExporterWorker_Work(object sender, DoWorkEventArgs e)
        {
            BackgroundWorker bg = sender as BackgroundWorker;
            if (!bg.CancellationPending)
            {
                string fullContent = PrometheusExporterHost.LogStream.FullLogOutput;
                this.tbExporter.Dispatcher.Invoke(() => { this.tbExporter.AppendText(fullContent); });

                LogLineArrived handler = CreateOnLogLineArrivedHandler(this.tbExporter);
                PrometheusExporterHost.LogStream.OnLogLine += handler;

                while (true)
                {
                    if (bg.CancellationPending)
                    {
                        break;
                    }

                    System.Threading.Thread.Sleep(2000);
                }

                PrometheusExporterHost.LogStream.OnLogLine -= handler;
            }
        }

        private void PrometheusWorker_Work(object sender, DoWorkEventArgs e)
        {
            BackgroundWorker bg = sender as BackgroundWorker;
            if (!bg.CancellationPending)
            {
                string fullContent = PrometheusHost.LogStream.FullLogOutput;
                this.tbPrometheus.Dispatcher.Invoke(() => { this.tbPrometheus.AppendText(fullContent); });

                LogLineArrived handler = CreateOnLogLineArrivedHandler(this.tbPrometheus);
                PrometheusHost.LogStream.OnLogLine += handler;

                while (true)
                {
                    if (bg.CancellationPending)
                    {
                        break;
                    }

                    System.Threading.Thread.Sleep(2000);
                }

                PrometheusHost.LogStream.OnLogLine -= handler;
            }
        }

        private void GrafanaWorker_Work(object sender, DoWorkEventArgs e)
        {
            BackgroundWorker bg = sender as BackgroundWorker;
            if (!bg.CancellationPending)
            {
                string fullContent = GrafanaHost.LogStream.FullLogOutput;
                this.tbGrafana.Dispatcher.Invoke(() => { this.tbGrafana.AppendText(fullContent); });

                LogLineArrived handler = CreateOnLogLineArrivedHandler(this.tbGrafana);
                GrafanaHost.LogStream.OnLogLine += handler;

                while (true)
                {
                    if (bg.CancellationPending)
                    {
                        break;
                    }

                    System.Threading.Thread.Sleep(2000);
                }

                GrafanaHost.LogStream.OnLogLine -= handler;
            }
        }

        private void CompanionWorker_Work(object sender, DoWorkEventArgs e)
        {
            BackgroundWorker bg = sender as BackgroundWorker;
            if (!bg.CancellationPending)
            {
                string fullContent = AppLogStgream.Instance.FullLogOutput;

                this.tbCompanion.Dispatcher.Invoke(() => { this.tbCompanion.AppendText(fullContent); });

                LogLineArrived handler = CreateOnLogLineArrivedHandler(this.tbCompanion);
                AppLogStgream.Instance.OnLogLine += handler;

                while (true)
                {
                    if (bg.CancellationPending)
                    {
                        break;
                    }

                    System.Threading.Thread.Sleep(2000);
                }

                AppLogStgream.Instance.OnLogLine -= handler;
            }
        }
        private LogLineArrived CreateOnLogLineArrivedHandler(TextBox target)
        {
            return (string data) =>
            {
                target.Dispatcher.Invoke(() =>
                {
                    target.AppendText(data + "\n");
                });
            };
        }

        private void Window_Closing(object sender, CancelEventArgs e)
        {
            if(exporterWorker != null)
            {
                exporterWorker.CancelAsync();
            }

            if (prometheusWorker != null)
            {
                prometheusWorker.CancelAsync();
            }

            if (grafanaWorker != null)
            {
                grafanaWorker.CancelAsync();
            }
        }
    }
}
