using System;
using System.Collections.Generic;
using System.Text;
using System.Text.Json.Serialization;
using System.Text.RegularExpressions;

namespace PrometheusExporter
{
    public class ProductionDetail
    {
        private readonly Regex _productionWithCapacityRegex = new Regex(
            "P: (?<prod_current>[\\d.]+)/(?<prod_capacity>[\\d.]+)/min - C: (?<cons_current>[\\d.]+)/(?<cons_capacity>[\\d.]+)/min"
        );

        private readonly Regex _productionWithoutCapacityRegex = new Regex(
            "P:(?<prod_current>[\\d.]+)/min - C: (?<cons_current>[\\d.]+)/min"
        );

        private MatchCollection _perMinStatsMatches;

        public string ItemName { get; set; }
        public string ProdPerMin { get; set; }
        [JsonPropertyName("ProdPercent")]
        public double ProductionPercent { get; set; }
        [JsonPropertyName("ConsPercent")]
        public double ConsumptionPercent { get; set; }
        [JsonPropertyName("CurrentProd")]
        public double CurrentProduction { get; set; }
        [JsonPropertyName("MaxProd")]
        [JsonConverter(typeof(DoublePossiblyNAJSONConverter))]
        public double MaxProduction { get; set; }
        [JsonPropertyName("CurrentConsumed")]
        public double CurrentConsumption { get; set; }
        [JsonPropertyName("MaxConsumed")]
        [JsonConverter(typeof(DoublePossiblyNAJSONConverter))]
        public double MaxConsumption { get; set; }

        [JsonIgnore]
        public double ProductionCapacity
        {
            get
            { 

                return PerMinStat("prod_capacity");
            }
        }

        [JsonIgnore]
        public double ConsumptionCapacity
        {
            get
            {

                return PerMinStat("cons_capacity");
            }
        }

        private double PerMinStat(string statName)
        {
            if (_perMinStatsMatches == null)
            {
                _perMinStatsMatches = _productionWithCapacityRegex.Matches(this.ProdPerMin);
                if (_perMinStatsMatches.Count == 0)
                {
                    _perMinStatsMatches = _productionWithoutCapacityRegex.Matches(this.ProdPerMin);
                }

                if (_perMinStatsMatches.Count == 0)
                {
                    return -1;
                }
            }

            Match m = _perMinStatsMatches[0];
            string statFigure = m.Groups[statName].Value;
            return double.Parse(statFigure);
        }
    }
}
