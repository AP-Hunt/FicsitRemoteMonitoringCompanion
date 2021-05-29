using System;
using System.Collections.Generic;
using System.IO;
using System.Text;
using System.Text.Json;
using System.Threading.Tasks;

namespace Companion.Config
{
    static class ConfigIO
    {
        public static async Task<ConfigFile> ReadConfigFile(string path)
        {
            if(!File.Exists(path)){
                return new ConfigFile();
            }

            string fileText = await File.ReadAllTextAsync(path);
            return JsonSerializer.Deserialize<ConfigFile>(fileText);
        }

        public static async Task WriteConfigFile(string path, ConfigFile config)
        {
            string jsonText = JsonSerializer.Serialize<ConfigFile>(config);
            await File.WriteAllTextAsync(path, jsonText);
        }
    }
}
