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
        static LogStream _logStream;
        internal static LogStream LogStream
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

            _logStream = new LogStream(_grafanaProcess);

            _grafanaProcess.BeginOutputReadLine();
        }

        internal async static Task WaitForReadiness()
        {
            await WaitAPIHealthy();
            await ConfigureGrafana();
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

            if(source.Token.IsCancellationRequested)
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
            await ConfigureOrgs(client);
            await CreateFolder(client);
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
                }
                catch(HttpRequestException ex)
                {
                    Console.WriteLine("Failed creating Prometheus data source: {0}", ex.Message);
                }
            }
        }

        private async static Task ConfigureOrgs(HttpClient client)
        {
            bool orgFound;
            try
            {
                var response = await client.GetAsync("http://localhost:3000/api/orgs/name/Ficsit");
                orgFound = (response.StatusCode != HttpStatusCode.NotFound);
            }
            catch (HttpRequestException)
            {
                orgFound = false;
            }

            if (!orgFound)
            {
                string bodyJson = @"
                        {
                            ""name"": ""Ficsit""
                        }
                    ";
                HttpContent content = new StringContent(bodyJson);
                content.Headers.ContentType = new System.Net.Http.Headers.MediaTypeHeaderValue("application/json");

                try
                {
                    await client.PostAsync("http://localhost:3000/api/orgs", content);
                }
                catch (HttpRequestException ex)
                {
                    Console.WriteLine("Failed creating Grafana org: {0}", ex.Message);
                }
            }
        }

        private async static Task CreateFolder(HttpClient client)
        {
            string bodyJson = @"
                        {
                            ""title"": ""Ficsit"",
                            ""uuid"": ""ficsit"",
                            ""overwrite"": true
                        }
                    ";
            HttpContent content = new StringContent(bodyJson);
            content.Headers.ContentType = new System.Net.Http.Headers.MediaTypeHeaderValue("application/json");

            try
            {
                await client.PutAsync("http://localhost:3000/api/folders/ficsit", content);
            }
            catch (HttpRequestException ex)
            {
                Console.WriteLine("Failed creating Grafana folder: {0}", ex.Message);
            }
        }

        private async static Task PopulateDashboards(HttpClient client)
        {
            string[] dashboards = new string[]{
                Dashboards.Power,
                Dashboards.Production
            };

            foreach(string dashboard in dashboards)
            {
                HttpContent content = new StringContent(dashboard);
                content.Headers.ContentType = new System.Net.Http.Headers.MediaTypeHeaderValue("application/json");

                try
                {
                    await client.PostAsync("http://localhost:3000/api/dashboards/db ", content);
                }
                catch (HttpRequestException ex)
                {
                    Console.WriteLine("Failed creating Grafana dashboard: {0}", ex.Message);
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
