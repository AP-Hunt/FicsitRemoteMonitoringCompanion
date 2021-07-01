using Microsoft.Web.WebView2.Core;
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
        TextBlock loadingLabel;
        bool firstLoad;
        public MainWindow()
        {
            InitializeComponent();
            loadingLabel = new TextBlock()
            {
                Text = @$"
                Loading Grafana...
                
                This can take up to 5 minutes the very first time.",
                HorizontalAlignment = HorizontalAlignment.Center,
                VerticalAlignment = VerticalAlignment.Center
            };
            this.MainPanel.Children.Add(loadingLabel);

            firstLoad = true;
        }

        public async void ShowGrafana()
        {
            try
            {
                AppLogStgream.Instance.WriteLine("initialising web view");
                var webViewEnv = await CoreWebView2Environment.CreateAsync();
                await webView.EnsureCoreWebView2Async(webViewEnv);
                webView.NavigationCompleted += WebView_NavigationCompleted;
                webView.Source = new Uri("http://localhost:3000");
            }
            catch(Exception ex)
            {
                AppLogStgream.Instance.WriteLine("error initialising web view: {0}", ex.Message);
            }
        }

        private void WebView_NavigationCompleted(object sender, CoreWebView2NavigationCompletedEventArgs e)
        {
            if(firstLoad && !e.IsSuccess)
            {
                AppLogStgream.Instance.WriteLine("Naviating to Grafana failed with error '{0}'", e.WebErrorStatus.ToString());
                loadingLabel.Text = "Something went wrong loading Grafana. Take a look at the diagnostics tab";
                firstLoad = false;
            }
            else if(firstLoad)
            {
                AppLogStgream.Instance.WriteLine("grafana page has loaded, showing browser window");
                Application.Current.Dispatcher.Invoke(() => {
                    this.MainPanel.Children.Remove(loadingLabel);
                    webView.Visibility = Visibility.Visible;
                    AppLogStgream.Instance.WriteLine("browser window was shown");
                });
                firstLoad = false;
            }
        }

        private async void Window_ContentRendered(object sender, EventArgs e)
        {
            try
            {
                await GrafanaHost.WaitForReadiness();
            }
            catch(TaskCanceledException ex)
            {
                MessageBox.Show("Grafana failed to start up. Try restarting the program.");
                this.Close();
                return;
            }

            this.ShowGrafana();
        }

        private void Diagnostics_Click(object sender, RoutedEventArgs e)
        {
            LogLines logLinesWindow = new LogLines();
            logLinesWindow.Owner = this;
            logLinesWindow.Show();
        }
    }
}   
