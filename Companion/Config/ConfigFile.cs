using System;
using System.Collections.Generic;
using System.IO;
using System.Text;
using System.Text.Json.Serialization;

namespace Companion.Config
{
    public class ConfigFile
    {
        [JsonPropertyName("satisfactory_game_directory")]
        public string SatisfactoryGameDirectory { get; set; }

        public bool IsValidGamePath()
        {
            string expectedExePath = Path.Combine(this.SatisfactoryGameDirectory, "FactoryGame.exe");
            return File.Exists(expectedExePath);
        }
    }
}
