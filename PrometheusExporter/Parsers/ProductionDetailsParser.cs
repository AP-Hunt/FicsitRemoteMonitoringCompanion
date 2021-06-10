using System;
using System.Collections.Generic;
using System.Text;

namespace PrometheusExporter.Parsers
{
    public class ProductionDetailsParser
    {
        public static List<ProductionDetail> ParseJSON(string json)
        {
            var options = new System.Text.Json.JsonSerializerOptions
            {
                AllowTrailingCommas = true,
                PropertyNameCaseInsensitive = true
            };

            List<ProductionDetail> productionDetails = System.Text.Json.JsonSerializer.Deserialize<List<ProductionDetail>>(json, options);

            return productionDetails;
        }
    }
}
