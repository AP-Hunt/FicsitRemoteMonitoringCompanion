using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Text;

namespace Companion
{
    static class PrometheusHost
    {
        static Process _prometheusProcess;

        internal static void Start()
        {
            string currentExeLocation = System.Reflection.Assembly.GetExecutingAssembly().Location;
            string currentExeDir = Path.GetDirectoryName(currentExeLocation);

            string prometheusWorkingDir = currentExeDir;
            string prometheusExePath = Path.Combine(prometheusWorkingDir, "prometheus.exe");

            string prometheusConfigPath = WriteGrafanaConfig(prometheusWorkingDir);

            ProcessStartInfo promProcessStartInfo = new ProcessStartInfo()
            {
                FileName = prometheusExePath,
                WorkingDirectory = prometheusWorkingDir,
                UseShellExecute = false,
                CreateNoWindow = true,

                Arguments = $"--config.file=\"{prometheusConfigPath}\""
            };
            _prometheusProcess = Process.Start(promProcessStartInfo);
        }

        private static string WriteGrafanaConfig(string prometheusWorkingDir)
        {
            string configPath = Path.Combine(prometheusWorkingDir, "config.yml");
            File.WriteAllText(configPath, ConfigFileResources.PrometheusConfiguration);
            return configPath;
        }

        internal static void Stop()
        {
            if (_prometheusProcess != null && !_prometheusProcess.HasExited)
            {
                _prometheusProcess.Kill();
            }
        }
    }
}
