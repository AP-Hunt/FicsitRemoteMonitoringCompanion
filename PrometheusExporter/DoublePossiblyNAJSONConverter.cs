using System;
using System.Collections.Generic;
using System.Text;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace PrometheusExporter
{
    public class DoublePossiblyNAJSONConverter : JsonConverter<double>
    {
        public override double Read(ref Utf8JsonReader reader, Type typeToConvert, JsonSerializerOptions options)
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
                            return -1;
                        }
                        else
                        {
                            return double.Parse(sValue);
                        }

                    default: 
                        return -1;
                }
            }
            catch(InvalidOperationException ex)
            {
                return -1;
            }
        }

        public override void Write(Utf8JsonWriter writer, double value, JsonSerializerOptions options)
        {
            writer.WriteStringValue(value.ToString());
        }
    }
}
