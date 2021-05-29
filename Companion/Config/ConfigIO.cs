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
        public static async Task<ConfigFile> ReadCompanionConfigFile(string path)
        {
            if(!File.Exists(path)){
                return new ConfigFile();
            }

            string fileText = await File.ReadAllTextAsync(path);
            return JsonSerializer.Deserialize<ConfigFile>(fileText);
        }

        public static async Task WriteCompanionConfigFile(string path, ConfigFile config)
        {
            string jsonText = JsonSerializer.Serialize<ConfigFile>(config);
            await File.WriteAllTextAsync(path, jsonText);
        }

        public static async Task<FicsitRemoteMonitoringConfig> ReadFRMConfigFile(string satisfactoryGameDir)
        {
            string configFilePath = Path.Combine(satisfactoryGameDir, "FactoryGame", "Configs", "FicsitRemoteMonitoring.cfg");
            if(!File.Exists(configFilePath))
            {
                throw new FileNotFoundException($"Could not find the config file at {configFilePath}");
            }

            string fileText = await File.ReadAllTextAsync(configFilePath);
            return JsonSerializer.Deserialize<FicsitRemoteMonitoringConfig>(fileText);
        }
    }
}
