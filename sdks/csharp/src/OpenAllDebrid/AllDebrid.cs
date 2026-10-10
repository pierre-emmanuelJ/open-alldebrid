using System.Globalization;
using System.Net.Http.Headers;
using System.Text.Json;
using System.Text.Json.Nodes;

namespace OpenAllDebrid;

/// <summary>A .torrent file to upload.</summary>
public sealed record TorrentFile(string Name, byte[] Content);

/// <summary>Magnet status codes. 5 and above are errors.</summary>
public static class MagnetStatus
{
    public const int Queued = 0;
    public const int Downloading = 1;
    public const int Compressing = 2;
    public const int Uploading = 3;
    public const int Ready = 4;

    /// <summary>Whether the magnet's files are available.</summary>
    public static bool IsReady(this Magnet magnet) => magnet.StatusCode == Ready;

    /// <summary>Whether the magnet failed.</summary>
    public static bool IsError(this Magnet magnet) => (magnet.StatusCode ?? 0) > Ready;
}

/// <summary>Delayed link statuses.</summary>
public static class DelayedStatus
{
    public const int Processing = 1;
    public const int Ready = 2;
    public const int Failed = 3;
}

/// <summary>
/// A client for the AllDebrid API. API errors throw <see cref="AllDebridException"/>,
/// other HTTP statuses <see cref="AllDebridHttpException"/>.
/// </summary>
/// <example>
/// <code>
/// var client = new AllDebrid(apiKey, userAgent: "my-app");
/// var user = await client.UserAsync();
/// </code>
/// </example>
public sealed class AllDebrid
{
    public const string DefaultBaseUrl = "https://api.alldebrid.com";
    public const string DefaultUserAgent = "open-alldebrid-csharp";

    /// <summary>Public demo keys: the API answers them with static sample data.</summary>
    public const string DemoPremiumKey = "staticDemoApikeyPrem";
    public const string DemoTrialKey = "staticDemoApikeyTria";
    public const string DemoFreeKey = "staticDemoApikeyFree";

    private static readonly HttpClient SharedHttp = new();

    internal static readonly JsonSerializerOptions Json = new()
    {
        NumberHandling = System.Text.Json.Serialization.JsonNumberHandling.AllowReadingFromString,
    };

    private readonly string? apiKey;
    private readonly string baseUrl;
    private readonly string userAgent;
    private readonly HttpClient http;

    /// <param name="apiKey">An API key; null for public endpoints and the PIN flow.</param>
    /// <param name="baseUrl">The API base URL, for tests or a proxy.</param>
    /// <param name="userAgent">Identifies your application to AllDebrid.</param>
    /// <param name="httpClient">The HTTP client; a shared one by default.</param>
    public AllDebrid(string? apiKey = null, string baseUrl = DefaultBaseUrl, string userAgent = DefaultUserAgent, HttpClient? httpClient = null)
    {
        this.apiKey = string.IsNullOrEmpty(apiKey) ? null : apiKey;
        this.baseUrl = baseUrl.TrimEnd('/');
        this.userAgent = userAgent;
        http = httpClient ?? SharedHttp;
    }

    /// <summary>A client with another API key, sharing this one's HTTP client.</summary>
    public AllDebrid WithApiKey(string? key) => new(key, baseUrl, userAgent, http);

    /// <summary>Calls the API and returns the "data" of a successful answer.</summary>
    internal async Task<JsonNode?> CallAsync(HttpMethod method, string path, IEnumerable<KeyValuePair<string, string>>? parameters, HttpContent? content, CancellationToken ct)
    {
        var pairs = parameters?.ToList() ?? new List<KeyValuePair<string, string>>();
        var url = baseUrl + path;
        if (method == HttpMethod.Get && pairs.Count > 0)
        {
            url += "?" + string.Join("&", pairs.Select(p => Uri.EscapeDataString(p.Key) + "=" + Uri.EscapeDataString(p.Value)));
        }
        using var request = new HttpRequestMessage(method, url);
        request.Headers.Accept.Add(new MediaTypeWithQualityHeaderValue("application/json"));
        request.Headers.TryAddWithoutValidation("User-Agent", userAgent);
        if (apiKey is not null)
        {
            request.Headers.Authorization = new AuthenticationHeaderValue("Bearer", apiKey);
        }
        if (content is not null)
        {
            request.Content = content;
        }
        else if (method != HttpMethod.Get)
        {
            // List parameters keep their PHP-style names and repeat: magnets[]=a&magnets[]=b.
            request.Content = new FormUrlEncodedContent(pairs);
        }

        using var response = await http.SendAsync(request, ct).ConfigureAwait(false);
        var text = await response.Content.ReadAsStringAsync(ct).ConfigureAwait(false);
        JsonObject? envelope = null;
        try
        {
            envelope = JsonNode.Parse(text) as JsonObject;
        }
        catch (JsonException)
        {
        }
        var apiError = ToException(envelope?["error"]);
        var status = (int)response.StatusCode;
        if (status != 200)
        {
            throw new AllDebridHttpException(status, apiError);
        }
        if (envelope is null)
        {
            throw new JsonException($"AllDebrid {method} {path}: the response is not a JSON object");
        }
        if ((string?)envelope["status"] != "success")
        {
            throw (Exception?)apiError ?? new JsonException($"AllDebrid {method} {path}: error without details");
        }
        return envelope["data"];
    }

    private static AllDebridException? ToException(JsonNode? error) => error is JsonObject e
        ? new AllDebridException(Text(e["code"]), Text(e["message"]), e["token"] is null ? null : Text(e["token"]))
        : null;

    private static string Text(JsonNode? node) => node is JsonValue v ? v.ToString() : "";

    private async Task<T> GetAsync<T>(string path, IEnumerable<KeyValuePair<string, string>>? parameters, CancellationToken ct) =>
        Decode<T>(await CallAsync(HttpMethod.Get, path, parameters, null, ct).ConfigureAwait(false), path);

    private async Task<T> PostAsync<T>(string path, IEnumerable<KeyValuePair<string, string>>? parameters, CancellationToken ct) =>
        Decode<T>(await CallAsync(HttpMethod.Post, path, parameters, null, ct).ConfigureAwait(false), path);

    private async Task PostAsync(string path, IEnumerable<KeyValuePair<string, string>>? parameters, CancellationToken ct) =>
        await CallAsync(HttpMethod.Post, path, parameters, null, ct).ConfigureAwait(false);

    private static T Decode<T>(JsonNode? data, string path) =>
        data.Deserialize<T>(Json) ?? throw new JsonException($"AllDebrid {path}: empty data");

    private static KeyValuePair<string, string> P(string name, string value) => new(name, value);

    private static KeyValuePair<string, string> P(string name, long value) => new(name, value.ToString(CultureInfo.InvariantCulture));

    private static IEnumerable<KeyValuePair<string, string>> Repeat(string name, IEnumerable<string> values) => values.Select(v => P(name, v));

    private static IEnumerable<KeyValuePair<string, string>> Repeat(string name, IEnumerable<long> values) => values.Select(v => P(name, v));

    private static KeyValuePair<string, string>[] HostsOnly(bool hostsOnly) => hostsOnly ? new[] { P("hostsOnly", "true") } : Array.Empty<KeyValuePair<string, string>>();

    // ------------------------------------------------------------------ general

    /// <summary>Checks that the API is up.</summary>
    public async Task PingAsync(CancellationToken ct = default) =>
        await CallAsync(HttpMethod.Get, "/v4/ping", null, null, ct).ConfigureAwait(false);

    // ---------------------------------------------------------------------- pin

    /// <summary>Starts a PIN authentication. Show Pin or UserUrl, then <see cref="WaitForPinAsync"/>.</summary>
    public Task<PinGetData> GetPinAsync(CancellationToken ct = default) => GetAsync<PinGetData>("/v4.1/pin/get", null, ct);

    /// <summary>The state of a PIN authentication; Apikey is set once activated.</summary>
    public Task<PinCheckData> CheckPinAsync(string pin, string check, CancellationToken ct = default) =>
        PostAsync<PinCheckData>("/v4/pin/check", new[] { P("pin", pin), P("check", check) }, ct);

    /// <summary>Polls until the user confirms the PIN and returns the new API key.</summary>
    public async Task<string> WaitForPinAsync(PinGetData pin, TimeSpan? interval = null, CancellationToken ct = default)
    {
        ArgumentNullException.ThrowIfNull(pin);
        while (true)
        {
            var state = await CheckPinAsync(pin.Pin, pin.Check, ct).ConfigureAwait(false);
            if (state.Activated && !string.IsNullOrEmpty(state.Apikey))
            {
                return state.Apikey;
            }
            await Task.Delay(interval ?? TimeSpan.FromSeconds(5), ct).ConfigureAwait(false);
        }
    }

    // -------------------------------------------------------------------- hosts

    /// <summary>Supported hosts, streaming sites and redirectors. Cache the result.</summary>
    public Task<HostsData> HostsAsync(bool hostsOnly = false, CancellationToken ct = default) =>
        GetAsync<HostsData>("/v4/hosts", HostsOnly(hostsOnly), ct);

    /// <summary>Every supported domain.</summary>
    public Task<HostDomainsData> HostDomainsAsync(CancellationToken ct = default) => GetAsync<HostDomainsData>("/v4/hosts/domains", null, ct);

    /// <summary>Host names ranked from the most open (1) to the most restricted.</summary>
    public async Task<IReadOnlyDictionary<string, int>> HostPriorityAsync(CancellationToken ct = default) =>
        (await GetAsync<HostPriorityData>("/v4/hosts/priority", null, ct).ConfigureAwait(false)).Hosts;

    // --------------------------------------------------------------------- user

    /// <summary>The account information.</summary>
    public async Task<User> UserAsync(CancellationToken ct = default) =>
        (await GetAsync<UserData>("/v4/user", null, ct).ConfigureAwait(false)).User;

    /// <summary>The hosts this user can use now, with their quotas.</summary>
    public async Task<IReadOnlyDictionary<string, Service>> UserHostsAsync(bool hostsOnly = false, CancellationToken ct = default) =>
        (await GetAsync<UserHostsData>("/v4.1/user/hosts", HostsOnly(hostsOnly), ct).ConfigureAwait(false)).Hosts;

    /// <summary>Follows a new-location login verification (token from an AUTH_BLOCKED error).</summary>
    public Task<VerificationData> VerificationStatusAsync(string token, CancellationToken ct = default) =>
        PostAsync<VerificationData>("/v4/user/verif", new[] { P("token", token) }, ct);

    /// <summary>Sends the verification email again (once per token).</summary>
    public Task ResendVerificationAsync(string token, CancellationToken ct = default) =>
        PostAsync("/v4/user/verif/resend", new[] { P("token", token) }, ct);

    /// <summary>Clears one of the codes in User.Notifications.</summary>
    public Task ClearNotificationAsync(string code, CancellationToken ct = default) =>
        PostAsync("/v4/user/notification/clear", new[] { P("code", code) }, ct);

    // -------------------------------------------------------------------- links

    /// <summary>File name and size of each link, in order; unreadable links have an Error.</summary>
    public async Task<IReadOnlyList<LinkInfo>> LinkInfosAsync(IEnumerable<string> links, string? password = null, CancellationToken ct = default)
    {
        var parameters = Repeat("link[]", links).ToList();
        if (password is not null)
        {
            parameters.Add(P("password", password));
        }
        return (await PostAsync<LinkInfosData>("/v4/link/infos", parameters, ct).ConfigureAwait(false)).Infos;
    }

    /// <summary>Links behind a redirector; they only work with <see cref="LinkInfosAsync"/> and <see cref="UnlockAsync"/>.</summary>
    public async Task<IReadOnlyList<string>> ExtractRedirectorAsync(string link, CancellationToken ct = default) =>
        (await PostAsync<RedirectorData>("/v4/link/redirector", new[] { P("link", link) }, ct).ConfigureAwait(false)).Links;

    /// <summary>
    /// Unlocks a hoster link. When Delayed is set, wait with <see cref="WaitDelayedLinkAsync"/>;
    /// when Streaming lists qualities, choose one with <see cref="StreamingLinkAsync"/>.
    /// </summary>
    public Task<UnlockData> UnlockAsync(string link, string? password = null, CancellationToken ct = default)
    {
        var parameters = new List<KeyValuePair<string, string>> { P("link", link) };
        if (password is not null)
        {
            parameters.Add(P("password", password));
        }
        return PostAsync<UnlockData>("/v4/link/unlock", parameters, ct);
    }

    /// <summary>The link of one quality: id from <see cref="UnlockAsync"/>, stream from one of its Streaming entries.</summary>
    public Task<StreamingData> StreamingLinkAsync(string id, string stream, CancellationToken ct = default) =>
        PostAsync<StreamingData>("/v4/link/streaming", new[] { P("id", id), P("stream", stream) }, ct);

    /// <summary>The status of a delayed link.</summary>
    public Task<DelayedData> DelayedLinkAsync(string id, CancellationToken ct = default) =>
        PostAsync<DelayedData>("/v4/link/delayed", new[] { P("id", id) }, ct);

    /// <summary>Polls a delayed link until it is ready and returns the download link.</summary>
    public async Task<string> WaitDelayedLinkAsync(string id, TimeSpan? interval = null, CancellationToken ct = default)
    {
        while (true)
        {
            var state = await DelayedLinkAsync(id, ct).ConfigureAwait(false);
            if (state.Status == DelayedStatus.Ready && !string.IsNullOrEmpty(state.Link))
            {
                return state.Link;
            }
            if (state.Status == DelayedStatus.Failed)
            {
                throw new InvalidOperationException($"AllDebrid: delayed link {id} failed");
            }
            await Task.Delay(interval ?? TimeSpan.FromSeconds(5), ct).ConfigureAwait(false);
        }
    }

    // ------------------------------------------------------------------ magnets

    /// <summary>Adds magnet URIs or info hashes. Results are in order; invalid ones have an Error.</summary>
    public async Task<IReadOnlyList<UploadedMagnet>> UploadMagnetsAsync(IEnumerable<string> magnets, CancellationToken ct = default) =>
        (await PostAsync<UploadMagnetsData>("/v4/magnet/upload", Repeat("magnets[]", magnets), ct).ConfigureAwait(false)).Magnets;

    /// <summary>Adds .torrent files. Results are in order; invalid files have an Error.</summary>
    public async Task<IReadOnlyList<UploadedMagnet>> UploadTorrentFilesAsync(IEnumerable<TorrentFile> files, CancellationToken ct = default)
    {
        ArgumentNullException.ThrowIfNull(files);
        using var form = new MultipartFormDataContent();
        foreach (var file in files)
        {
            var part = new ByteArrayContent(file.Content);
            part.Headers.ContentType = new MediaTypeHeaderValue("application/x-bittorrent");
            form.Add(part, "files[]", file.Name);
        }
        var data = await CallAsync(HttpMethod.Post, "/v4/magnet/upload/file", null, form, ct).ConfigureAwait(false);
        return Decode<UploadTorrentFilesData>(data, "/v4/magnet/upload/file").Files;
    }

    /// <summary>Every magnet of the account.</summary>
    public async Task<IReadOnlyList<Magnet>> MagnetsAsync(CancellationToken ct = default) =>
        (await PostAsync<MagnetStatusData>("/v4.1/magnet/status", null, ct).ConfigureAwait(false)).Magnets;

    /// <summary>One magnet, with its Files.</summary>
    public async Task<Magnet> MagnetAsync(long id, CancellationToken ct = default)
    {
        var found = (await PostAsync<MagnetStatusData>("/v4.1/magnet/status", new[] { P("id", id) }, ct).ConfigureAwait(false)).Magnets;
        return found.Count == 1 ? found[0] : throw new JsonException($"AllDebrid: magnet {id}: got {found.Count} magnets");
    }

    /// <summary>The file tree of each magnet, in order; unknown IDs have an Error.</summary>
    public async Task<IReadOnlyList<MagnetFiles>> MagnetFilesAsync(IEnumerable<long> ids, CancellationToken ct = default) =>
        (await PostAsync<MagnetFilesData>("/v4/magnet/files", Repeat("id[]", ids), ct).ConfigureAwait(false)).Magnets;

    /// <summary>Deletes a magnet.</summary>
    public Task DeleteMagnetAsync(long id, CancellationToken ct = default) => PostAsync("/v4/magnet/delete", new[] { P("id", id) }, ct);

    /// <summary>Restarts a failed magnet.</summary>
    public Task RestartMagnetAsync(long id, CancellationToken ct = default) => PostAsync("/v4/magnet/restart", new[] { P("id", id) }, ct);

    /// <summary>Restarts failed magnets; each result has its own Error.</summary>
    public async Task<IReadOnlyList<RestartedMagnet>> RestartMagnetsAsync(IEnumerable<long> ids, CancellationToken ct = default) =>
        (await PostAsync<RestartMagnetsData>("/v4/magnet/restart", Repeat("ids[]", ids), ct).ConfigureAwait(false)).Magnets
        ?? Array.Empty<RestartedMagnet>();

    /// <summary>Follows the magnets with live mode: each <see cref="LiveMagnets.PollAsync"/> only transfers changes.</summary>
    public LiveMagnets LiveMagnets() => new(this);

    // -------------------------------------------------------------- saved links

    /// <summary>The links the user saved.</summary>
    public async Task<IReadOnlyList<SavedLink>> SavedLinksAsync(CancellationToken ct = default) =>
        (await GetAsync<SavedLinksData>("/v4/user/links", null, ct).ConfigureAwait(false)).Links;

    /// <summary>Saves links for later.</summary>
    public Task SaveLinksAsync(IEnumerable<string> links, CancellationToken ct = default) => PostAsync("/v4/user/links/save", Repeat("links[]", links), ct);

    /// <summary>Deletes saved links.</summary>
    public Task DeleteSavedLinksAsync(IEnumerable<string> links, CancellationToken ct = default) => PostAsync("/v4/user/links/delete", Repeat("links[]", links), ct);

    /// <summary>Links unlocked in the last 3 days, when history logging is enabled.</summary>
    public async Task<IReadOnlyList<SavedLink>> HistoryAsync(CancellationToken ct = default) =>
        (await GetAsync<SavedLinksData>("/v4/user/history", null, ct).ConfigureAwait(false)).Links;

    /// <summary>Deletes the whole recent links history.</summary>
    public Task PurgeHistoryAsync(CancellationToken ct = default) => PostAsync("/v4/user/history/delete", null, ct);

    // ----------------------------------------------------------------- reseller

    /// <summary>The reseller balance. Reseller accounts only.</summary>
    public async Task<double> VoucherBalanceAsync(CancellationToken ct = default) =>
        (await GetAsync<VoucherBalanceData>("/v4/voucher/balance", null, ct).ConfigureAwait(false)).Balance;

    /// <summary>Takes vouchers (15, 30, 90, 180 or 365 days) from the reseller stock.</summary>
    public Task<VouchersData> VouchersAsync(int durationDays, int count, CancellationToken ct = default) =>
        PostAsync<VouchersData>("/v4/voucher/get", new[] { P("duration", durationDays), P("nb", count) }, ct);

    /// <summary>Generates new vouchers from the reseller balance.</summary>
    public Task<GeneratedVouchersData> GenerateVouchersAsync(int durationDays, int count, CancellationToken ct = default) =>
        PostAsync<GeneratedVouchersData>("/v4/voucher/generate", new[] { P("duration", durationDays), P("nb", count) }, ct);
}
