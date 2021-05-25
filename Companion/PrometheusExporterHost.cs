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
            string currentExeDir = Path.GetDirectoryName(Process.GetCurrentProcess().MainModule.FileName);

            string prometheusExporterWorkingDir = Path.Combine(currentExeDir);
            string prometheusExporterExePath = Path.Combine(prometheusExporterWorkingDir, "PrometheusExporter.exe");

            try
            {
                ProcessStartInfo promExporterProcessStartInfo = new ProcessStartInfo()
                {
                    FileName = prometheusExporterExePath,
                    WorkingDirectory = prometheusExporterWorkingDir,
                    UseShellExecute = false,
                    CreateNoWindow = true
                };
                _prometheusExporterProcess = Process.Start(promExporterProcessStartInfo);
            }
            catch(Exception ex)
            {
                Console.WriteLine(ex.Message);
                Console.WriteLine("Working dir: {0}\nExe path: {1}", prometheusExporterWorkingDir, prometheusExporterExePath);
            }
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
