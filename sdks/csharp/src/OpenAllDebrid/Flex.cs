using System.Globalization;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace OpenAllDebrid;

/// <summary>
/// A value the API sends either as a JSON string or as a JSON number, such as a
/// delayed link ID. It keeps the value as a string.
/// </summary>
[JsonConverter(typeof(FlexStringConverter))]
public readonly record struct FlexString(string Value)
{
    /// <summary>The value as an integer, when it is one.</summary>
    public long? AsLong() => long.TryParse(Value, NumberStyles.Integer, CultureInfo.InvariantCulture, out var n) ? n : null;

    public override string ToString() => Value;
}

/// <summary>Reads a JSON string or number into a <see cref="FlexString"/>.</summary>
public sealed class FlexStringConverter : JsonConverter<FlexString>
{
    public override FlexString Read(ref Utf8JsonReader reader, Type typeToConvert, JsonSerializerOptions options) =>
        reader.TokenType switch
        {
            JsonTokenType.String => new FlexString(reader.GetString() ?? ""),
            JsonTokenType.Number => new FlexString(RawNumber(reader)),
            _ => throw new JsonException($"expected a string or a number, got {reader.TokenType}"),
        };

    private static string RawNumber(Utf8JsonReader reader) =>
        System.Text.Encoding.UTF8.GetString(reader.HasValueSequence ? System.Buffers.BuffersExtensions.ToArray(reader.ValueSequence) : reader.ValueSpan);

    public override void Write(Utf8JsonWriter writer, FlexString value, JsonSerializerOptions options) =>
        writer.WriteStringValue(value.Value);
}

/// <summary>Reads one string or a list of strings into a list.</summary>
public sealed class StringsConverter : JsonConverter<IReadOnlyList<string>>
{
    public override IReadOnlyList<string> Read(ref Utf8JsonReader reader, Type typeToConvert, JsonSerializerOptions options)
    {
        if (reader.TokenType == JsonTokenType.String)
        {
            return new[] { reader.GetString() ?? "" };
        }
        return JsonSerializer.Deserialize<List<string>>(ref reader, options) ?? new List<string>();
    }

    public override void Write(Utf8JsonWriter writer, IReadOnlyList<string> value, JsonSerializerOptions options) =>
        JsonSerializer.Serialize(writer, value, options);
}

/// <summary>Reads the magnets of a status: a list, or a single object when one magnet was asked for.</summary>
public sealed class MagnetListConverter : JsonConverter<IReadOnlyList<Magnet>>
{
    public override IReadOnlyList<Magnet> Read(ref Utf8JsonReader reader, Type typeToConvert, JsonSerializerOptions options)
    {
        if (reader.TokenType == JsonTokenType.StartObject)
        {
            var one = JsonSerializer.Deserialize<Magnet>(ref reader, options) ?? throw new JsonException("null magnet");
            return new[] { one };
        }
        return JsonSerializer.Deserialize<List<Magnet>>(ref reader, options) ?? new List<Magnet>();
    }

    public override void Write(Utf8JsonWriter writer, IReadOnlyList<Magnet> value, JsonSerializerOptions options) =>
        JsonSerializer.Serialize(writer, value, options);
}

/// <summary>Reads quotas per host; the API sends an empty map as an empty list.</summary>
public sealed class QuotasConverter : JsonConverter<IReadOnlyDictionary<string, long>>
{
    public override IReadOnlyDictionary<string, long> Read(ref Utf8JsonReader reader, Type typeToConvert, JsonSerializerOptions options)
    {
        if (reader.TokenType == JsonTokenType.StartArray)
        {
            reader.Read();
            if (reader.TokenType != JsonTokenType.EndArray)
            {
                throw new JsonException("quotas sent as a non-empty list");
            }
            return new Dictionary<string, long>();
        }
        return JsonSerializer.Deserialize<Dictionary<string, long>>(ref reader, options) ?? new Dictionary<string, long>();
    }

    public override void Write(Utf8JsonWriter writer, IReadOnlyDictionary<string, long> value, JsonSerializerOptions options) =>
        JsonSerializer.Serialize(writer, value, options);
}
