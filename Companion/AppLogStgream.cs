using System;
using System.Collections.Generic;
using System.Text;

namespace Companion
{
    internal class AppLogStgream : ILogStream
    {
        public static AppLogStgream Instance { get; private set; }

        private StringBuilder _sb;
        public string FullLogOutput => _sb.ToString();

        public event LogLineArrived OnLogLine;

        static AppLogStgream()
        {
            Instance = new AppLogStgream();
        }

        private AppLogStgream()
        {
            _sb = new StringBuilder();
        }

        public void WriteLine(string format, params object[] args)
        {
            string line = string.Format("[{0}] {1}", DateTime.Now.ToString("O"), string.Format(format, args));

            _sb.AppendLine(line);
            if(OnLogLine != null)
            {
                OnLogLine(line);
            }
        }
    }
}
