using System;
using System.Collections.Generic;
using System.Text;
using System.Text.Json.Serialization;

namespace Companion.Config
{
    public class FicsitRemoteMonitoringConfig
    {
        [JsonPropertyName("Listen_IP")]
        public string ListenIp { get; set; }
        [JsonPropertyName("HTTP_Port")]
        public int HTTPPort { get; set; }

        public Uri ListenAddress 
        {
            get
            {
                var builder = new UriBuilder();
                builder.Host = ListenIp;
                builder.Port = HTTPPort;
                builder.Scheme = "http";
                return builder.Uri;
            }
        }
    }
}
