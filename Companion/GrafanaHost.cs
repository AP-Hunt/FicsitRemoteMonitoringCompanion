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
            bool statusOK = false;
            while (!statusOK)
            {
                HttpWebRequest req = (HttpWebRequest)WebRequest.Create("http://localhost:3000/api/health");
                req.Timeout = 3000;
                HttpWebResponse resp = (HttpWebResponse)req.GetResponse();
                statusOK = (resp.StatusCode == HttpStatusCode.OK);
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
