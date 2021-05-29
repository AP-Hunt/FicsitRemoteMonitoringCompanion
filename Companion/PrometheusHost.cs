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
            string prometheusWorkingDir = Paths.RootDirectory;
            string prometheusExePath = Path.Combine(prometheusWorkingDir, "prometheus.exe");

            string prometheusConfigPath = WritePrometheusConfig(prometheusWorkingDir);

            try
            {
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
            catch (Exception ex)
            {
                Console.WriteLine(ex.Message);
                Console.WriteLine("Working dir: {0}\nExe path: {1}", prometheusWorkingDir, prometheusExePath);
            }
        }

        private static string WritePrometheusConfig(string prometheusWorkingDir)
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
