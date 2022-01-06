using System;
using System.Collections.Generic;
using System.Text;
using System.Text.Json.Serialization;
using System.Text.RegularExpressions;

namespace PrometheusExporter
{
    public class FactoryBuildingDetail
    {
        [JsonPropertyName("building")]
        public string Building { get; set; }

        [JsonPropertyName("Recipe")]
        public string Recipe { get; set; }

        [JsonPropertyName("location")]
        public MachineLocation Location { get;set; }

        [JsonPropertyName("ingredients")]
        public List<RecipeIngredientDetail> Ingredients { get; set; }

        [JsonPropertyName("production")]
        public List<RecipeProductionDetail> Production { get; set ;}
    }

    public class RecipeIngredientDetail 
    {
        [JsonPropertyName("Name")]
        public string Name { get; set; }

        [JsonPropertyName("CurrentConsumed")]
        public double CurrentConsumed { get; set; }

        [JsonPropertyName("MaxConsumed")]
        public double MaxConsumed { get; set; }

        [JsonPropertyName("ConsPercent")]
        public string ConsumptionPercent { get; set; }
    }

    public class RecipeProductionDetail
    {
        [JsonPropertyName("Name")]
        public string Name { get; set; }

        [JsonPropertyName("CurrentProd")]
        public double CurrentProduced { get; set; }

        [JsonPropertyName("MaxProd")]
        public double MaxProduced { get; set; }

        [JsonPropertyName("ProdPercent")]
        public string ProductionPercent { get; set; }
    }

    public class MachineLocation
    {
        [JsonPropertyName("x")]
        public double X { get; set; }

        [JsonPropertyName("y")]
        public double Y { get; set; }

        [JsonPropertyName("z")]
        public double Z { get; set; }

        [JsonPropertyName("rotation")]
        public int Rotation { get; set; }
    }
}
