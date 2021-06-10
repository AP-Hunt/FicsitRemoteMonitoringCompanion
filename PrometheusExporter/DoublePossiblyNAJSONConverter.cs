using System;
using System.Collections.Generic;
using System.Text;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace PrometheusExporter
{
    public class DoublePossiblyNAJSONConverter : JsonConverter<Nullable<double>>
    {
        public override double? Read(ref Utf8JsonReader reader, Type typeToConvert, JsonSerializerOptions options)
        {
            try
            {
                switch(reader.TokenType)
                {
                    case JsonTokenType.Number:
                        return reader.GetDouble();

                    case JsonTokenType.String:
                        string sValue = reader.GetString();
                        if (sValue.ToUpperInvariant() == "N/A")
                        {
                            return null;
                        }
                        else
                        {
                            return double.Parse(sValue);
                        }

                    default: 
                        return null;
                }
            }
            catch(InvalidOperationException)
            {
                return null;
            }
        }

        public override void Write(Utf8JsonWriter writer, double? value, JsonSerializerOptions options)
        {
            if (value.HasValue)
            {
                writer.WriteStringValue(value.ToString());
            }
        }
    }
}
