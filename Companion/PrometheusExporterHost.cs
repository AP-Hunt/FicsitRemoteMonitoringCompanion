using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Text;

namespace Companion
{
    static class PrometheusExporterHost
    {
        static Process _prometheusExporterProcess;

        internal static void Start()
        {
            string currentExeLocation = System.Reflection.Assembly.GetExecutingAssembly().Location;
            string currentExeDir = Path.GetDirectoryName(currentExeLocation);

            string prometheusExporterWorkingDir = Path.Combine(currentExeDir, "prometheus-exporter");
            string prometheusExporterExePath = Path.Combine(prometheusExporterWorkingDir, "PrometheusExporter.exe");

            ProcessStartInfo promExporterProcessStartInfo = new ProcessStartInfo()
            {
                FileName = prometheusExporterExePath,
                RedirectStandardOutput = false,
                RedirectStandardError = false,
                WorkingDirectory = prometheusExporterWorkingDir,
                UseShellExecute = true,
            };
            _prometheusExporterProcess = Process.Start(promExporterProcessStartInfo);
        }

        internal static void Stop()
        {
            if (_prometheusExporterProcess != null && !_prometheusExporterProcess.HasExited)
            {
                _prometheusExporterProcess.Kill();
            }
        }
    }
}
