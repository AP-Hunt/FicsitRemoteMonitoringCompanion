using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Text;

namespace Companion
{
    static class Paths
    {
        public static string RootDirectory
        {
            get
            {
                return Path.GetDirectoryName(Process.GetCurrentProcess().MainModule.FileName);
            }
        }

        public static string ConfigFile
        {
            get
            {
                return Path.Combine(Paths.RootDirectory, "companion_config.json");
            }
        }
    }
}
