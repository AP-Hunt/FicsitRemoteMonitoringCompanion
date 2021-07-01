using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Net;
using System.Net.Http;
using System.Reflection;
using System.Text;
using System.Threading;
using System.Threading.Tasks;

namespace Companion
{
    static class GrafanaHost
    {
        static Process _grafanaProcess;
        static ProcessLogStream _logStream;
        internal static ILogStream LogStream
        {
            get
            {
                return _logStream;
            }
        }

        internal static void Start()
        {
            string grafanaWorkingDir = Path.Combine(Paths.RootDirectory, "grafana", "bin");
            string grafanaExePath = Path.Combine(grafanaWorkingDir, "grafana-server.exe");

            string grafanaConfigPath = WriteGrafanaConfig(grafanaWorkingDir);

            _grafanaProcess = Process.Start(new ProcessStartInfo()
            {
                FileName = grafanaExePath,
                WorkingDirectory = grafanaWorkingDir,
                UseShellExecute = false,
                CreateNoWindow = true,
                RedirectStandardOutput = true,

                Arguments = $"-config \"{grafanaConfigPath}\""
            });

            _logStream = new ProcessLogStream(_grafanaProcess);

            _grafanaProcess.BeginOutputReadLine();
        }

        internal async static Task WaitForReadiness()
        {
            await WaitAPIHealthy();
            await ConfigureGrafana();
            await WaitUIReady();
        }

        private async static Task WaitAPIHealthy()
        {
            CancellationTokenSource source = new CancellationTokenSource();
            source.CancelAfter(TimeSpan.FromMinutes(5));

            bool statusOK = false;
            while (!statusOK && !source.Token.IsCancellationRequested)
            {
                try
                {
                    using (var client = new HttpClient())
                    {
                        client.Timeout = TimeSpan.FromSeconds(3);
                        AppLogStgream.Instance.WriteLine("waiting for Grafana API to be available");
                        var response = await client.GetAsync("http://localhost:3000/api/health");
                        statusOK = (response.StatusCode == HttpStatusCode.OK);
                    }
                }
                catch (TaskCanceledException)
                {
                    continue;
                }
                catch (HttpRequestException)
                {
                    continue;
                }
            }

            if (statusOK)
            {
                AppLogStgream.Instance.WriteLine("Grafana API is ready");
            }

            if (source.Token.IsCancellationRequested)
            {
                throw new TaskCanceledException();
            }
        }

        private async static Task WaitUIReady()
        {
            CancellationTokenSource source = new CancellationTokenSource();
            source.CancelAfter(TimeSpan.FromMinutes(5));

            bool statusOK = false;
            while (!statusOK && !source.Token.IsCancellationRequested)
            {
                try
                {
                    using (var client = new HttpClient())
                    {
                        client.Timeout = TimeSpan.FromSeconds(3);
                        AppLogStgream.Instance.WriteLine("waiting for Grafana UI to be ready");
                        var response = await client.GetAsync("http://localhost:3000/");

                        long? contentLength = 0;
                        if(response.Content != null)
                        {
                            await response.Content.LoadIntoBufferAsync();
                            contentLength = response.Content.Headers.ContentLength;
                        }

                        statusOK = (response.StatusCode == HttpStatusCode.OK && contentLength > 0);
                    }
                }
                catch (TaskCanceledException)
                {
                    continue;
                }
                catch (HttpRequestException)
                {
                    continue;
                }
            }

            if (statusOK)
            {
                AppLogStgream.Instance.WriteLine("Grafana UI is ready");
            }

            if (source.Token.IsCancellationRequested)
            {
                throw new TaskCanceledException();
            }
        }

        private async static Task ConfigureGrafana()
        {
            HttpClient client = new HttpClient();
            client.DefaultRequestHeaders.Add("Authorization", string.Format(
                "Basic {0}",
                Convert.ToBase64String(Encoding.Default.GetBytes("ficsit:pioneer"))
            ));
            client.DefaultRequestHeaders.Add("Accept", "application/json");

            await CreatePrometheusDataSource(client);
            await PopulateDashboards(client);
        }

        private async static Task CreatePrometheusDataSource(HttpClient client)
        {
            bool datasourceFound;
            try
            {
                var response = await client.GetAsync("http://localhost:3000/api/datasources/uid/prometheus");
                datasourceFound = (response.StatusCode != HttpStatusCode.NotFound);
            }
            catch(HttpRequestException)
            {
                datasourceFound = false;
            }

            if(!datasourceFound){
                AppLogStgream.Instance.WriteLine("Prometheus data source was not found in Grafana. Creating.");
                string bodyJson = @"
                        {
                            ""name"": ""prometheus"",
                            ""type"": ""prometheus"",
                            ""uid"": ""prometheus"",
                            ""url"": ""http://localhost:9090"",
                            ""access"": ""proxy""
                        }
                    ";
                HttpContent content = new StringContent(bodyJson);
                content.Headers.ContentType = new System.Net.Http.Headers.MediaTypeHeaderValue("application/json");

                try
                {
                    await client.PostAsync("http://localhost:3000/api/datasources", content);
                    AppLogStgream.Instance.WriteLine("Created Prometheus data source in Grafana");
                }
                catch(HttpRequestException ex)
                {
                    AppLogStgream.Instance.WriteLine("Failed creating Prometheus data source: {0}", ex.Message);
                }
            }
        }

        private async static Task PopulateDashboards(HttpClient client)
        {
            var dashboards = new Dictionary<string, string>{
                { "Power", Dashboards.Power },
                { "Production", Dashboards.Production }
            };

            foreach(var kvp in dashboards)
            {
                HttpContent content = new StringContent(kvp.Value);
                content.Headers.ContentType = new System.Net.Http.Headers.MediaTypeHeaderValue("application/json");

                try
                {
                    await client.PostAsync("http://localhost:3000/api/dashboards/db ", content);
                    AppLogStgream.Instance.WriteLine("Created Grafana dashboard: {0}", kvp.Key);
                }
                catch (HttpRequestException ex)
                {
                    AppLogStgream.Instance.WriteLine("Failed creating Grafana dashboard {0}: {1}", kvp.Key, ex.Message);
                }
            }
        }


        private static string WriteGrafanaConfig(string grafanaWorkingDir)
        {
            string configPath = Path.Combine(grafanaWorkingDir, "config.ini");
            File.WriteAllText(configPath, ConfigFileResources.GranaConfiguration);
            return configPath;
        }

        internal static void Stop()
        {
            if(_grafanaProcess != null && !_grafanaProcess.HasExited)
            {
                _grafanaProcess.Kill();
            }
        }
    }
}
