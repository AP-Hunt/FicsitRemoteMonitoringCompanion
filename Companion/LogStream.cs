using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Text;

namespace Companion
{
    internal delegate void LogLineArrived(string data);

    internal class LogStream
    {
        private readonly Process _process;
        private StringBuilder _builder;

        public event LogLineArrived OnLogLine;

        public LogStream(Process process)
        {
            this._process= process;

            _builder = new StringBuilder();
            _process.OutputDataReceived += _process_OutputDataReceived;
        }

        public string FullLogOutput
        {
            get 
            {
                return _builder.ToString();
            }
        }

        private void _process_OutputDataReceived(object sender, DataReceivedEventArgs e)
        {
            _builder.AppendLine(e.Data);

            if (OnLogLine != null)
            {
                OnLogLine(e.Data);
            }
        }
    }
}
