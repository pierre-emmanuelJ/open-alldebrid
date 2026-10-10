using System.Globalization;
using System.Security.Cryptography;
using System.Text.Json;
using System.Text.Json.Nodes;

namespace OpenAllDebrid;

/// <summary>
/// Live mode state: the merged list of magnets, updated by <see cref="PollAsync"/>.
/// Poll from one task at a time.
/// </summary>
public sealed class LiveMagnets
{
    private readonly AllDebrid client;
    private long counter;
    private Dictionary<long, JsonObject> state = new();

    internal LiveMagnets(AllDebrid client)
    {
        this.client = client;
        Session = RandomNumberGenerator.GetInt32(int.MaxValue);
    }

    /// <summary>The session ID sent to the API.</summary>
    public long Session { get; }

    /// <summary>Applies the changes since the last call and returns every magnet, newest first.</summary>
    public async Task<IReadOnlyList<Magnet>> PollAsync(CancellationToken ct = default)
    {
        var parameters = new[]
        {
            new KeyValuePair<string, string>("session", Session.ToString(CultureInfo.InvariantCulture)),
            new KeyValuePair<string, string>("counter", counter.ToString(CultureInfo.InvariantCulture)),
        };
        var data = await client.CallAsync(HttpMethod.Post, "/v4.1/magnet/status", parameters, null, ct).ConfigureAwait(false) as JsonObject
            ?? throw new JsonException("AllDebrid: live mode answered without data");
        if (data["fullsync"]?.GetValue<bool>() == true)
        {
            state = new Dictionary<long, JsonObject>();
        }
        var changes = data["magnets"] switch
        {
            JsonArray list => list.OfType<JsonObject>().ToList(),
            JsonObject one => new List<JsonObject> { one },
            _ => new List<JsonObject>(),
        };
        foreach (var change in changes)
        {
            if (change["id"] is not JsonValue idValue || !idValue.TryGetValue<long>(out var id))
            {
                continue;
            }
            if (change["deleted"]?.GetValue<bool>() == true)
            {
                state.Remove(id);
                continue;
            }
            if (!state.TryGetValue(id, out var current))
            {
                current = new JsonObject();
                state[id] = current;
            }
            foreach (var (key, value) in change)
            {
                current[key] = value?.DeepClone();
            }
        }
        counter = data["counter"]?.GetValue<long>() ?? 0;

        return state.Values
            .Select(m => m.Deserialize<Magnet>(AllDebrid.Json) ?? throw new JsonException("null magnet"))
            .OrderByDescending(m => m.UploadDate ?? 0)
            .ThenByDescending(m => m.Id)
            .ToList();
    }
}
