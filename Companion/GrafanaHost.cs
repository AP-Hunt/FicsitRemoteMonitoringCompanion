using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Net;
using System.Text;

namespace Companion
{
    static class GrafanaHost
    {
        static Process _grafanaProcess;

        internal static void Start()
        {
            string currentExeLocation = System.Reflection.Assembly.GetExecutingAssembly().Location;
            string currentExeDir = Path.GetDirectoryName(currentExeLocation);

            string grafanaWorkingDir = Path.Combine(currentExeDir, "grafana", "bin");
            string grafanaExePath = Path.Combine(grafanaWorkingDir, "grafana-server.exe");

            string grafanaConfigPath = WriteGrafanaConfig(grafanaWorkingDir);

            _grafanaProcess = Process.Start(new ProcessStartInfo()
            {
                FileName = grafanaExePath,
                RedirectStandardOutput = false,
                RedirectStandardError = false,
                WorkingDirectory = grafanaWorkingDir,
                UseShellExecute = true,

                Arguments = $"-config \"{grafanaConfigPath}\""
            });


        }

        internal static void WaitForReadiness()
        {
            WaitAPIHealthy();
            ConfigureGrafana();
        }

        private static void WaitAPIHealthy()
        {
            bool statusOK = false;
            while (!statusOK)
            {
                try
                {
                    HttpWebRequest req = (HttpWebRequest)WebRequest.Create("http://localhost:3000/api/health");
                    req.Timeout = 3000;
                    HttpWebResponse resp = (HttpWebResponse)req.GetResponse();
                    statusOK = (resp.StatusCode == HttpStatusCode.OK);
                }
                catch(WebException)
                {
                    continue;
                }
            }

        }

        private static void ConfigureGrafana()
        {
            HttpWebRequest req = (HttpWebRequest)WebRequest.Create("http://localhost:3000/api/datasources");
            req.Method = "POST";
            req.Accept = "application/json";
            req.ContentType = "application/json";
            req.Headers["Authorization"] = string.Format(
                "Basic {0}",
                Convert.ToBase64String(Encoding.Default.GetBytes("ficsit:pioneer"))
            );
           
            using(Stream reqStream = req.GetRequestStream())
            {
                using(StreamWriter writer = new StreamWriter(reqStream))
                {
                    string bodyJson = @"
                        {
                            ""name"": ""prometheus"",
                            ""type"": ""prometheus"",
                            ""url"": ""http://localhost:9090"",
                            ""access"": ""proxy""
                        }
                    ";
                    writer.Write(bodyJson);
                }
            }

            try
            {
                HttpWebResponse resp = (HttpWebResponse)req.GetResponse();

            }
            catch(WebException ex)
            {
                Console.WriteLine("Grafana: Create datasource returned {0}", ex.Status);

                using (Stream respStream = ex.Response.GetResponseStream())
                {
                    using (StreamReader rdr = new StreamReader(respStream))
                    {
                        string respText = rdr.ReadToEnd();
                        Console.WriteLine(respText);
                    }
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
