using System.Net;
using System.Text.RegularExpressions;
using Xunit;

namespace OpenAllDebrid.Tests;

/// <summary>Records requests and answers canned bodies in order.</summary>
internal sealed class FakeHandler : HttpMessageHandler
{
    private readonly Queue<(int Status, string Body)> answers = new();

    public List<(HttpMethod Method, Uri Uri, Dictionary<string, string> Headers, string? ContentType, string Body)> Requests { get; } = new();

    public void Answer(int status, string body) => answers.Enqueue((status, body));

    public void Ok(string data) => Answer(200, $$"""{"status":"success","data":{{data}}}""");

    protected override async Task<HttpResponseMessage> SendAsync(HttpRequestMessage request, CancellationToken cancellationToken)
    {
        var headers = request.Headers.ToDictionary(h => h.Key, h => string.Join(",", h.Value));
        var body = request.Content is null ? "" : await request.Content.ReadAsStringAsync(cancellationToken);
        Requests.Add((request.Method, request.RequestUri!, headers, request.Content?.Headers.ContentType?.ToString(), body));
        var (status, text) = answers.Count > 0 ? answers.Dequeue() : (200, """{"status":"success","data":{}}""");
        return new HttpResponseMessage((HttpStatusCode)status) { Content = new StringContent(text) };
    }

    public (HttpMethod Method, Uri Uri, Dictionary<string, string> Headers, string? ContentType, string Body) Last => Requests[^1];

    public Dictionary<string, List<string>> LastForm()
    {
        var form = new Dictionary<string, List<string>>();
        foreach (var pair in Last.Body.Split('&', StringSplitOptions.RemoveEmptyEntries))
        {
            var kv = pair.Split('=', 2);
            var key = Uri.UnescapeDataString(kv[0].Replace('+', ' '));
            var value = Uri.UnescapeDataString((kv.Length > 1 ? kv[1] : "").Replace('+', ' '));
            if (!form.TryGetValue(key, out var list))
            {
                form[key] = list = new List<string>();
            }
            list.Add(value);
        }
        return form;
    }
}

public sealed class ClientTests
{
    private static readonly string SpecPath = Path.GetFullPath(Path.Combine(AppContext.BaseDirectory, "../../../../../../../openapi.yaml"));

    private const string UserJson = """{"user":{"username":"u","email":"e","isPremium":true,"isSubscribed":false,"isTrial":false,"premiumUntil":0,"lang":"en","preferedDomain":"com","fidelityPoints":0,"limitedHostersQuotas":[],"notifications":[]}}""";

    private readonly FakeHandler fake = new();
    private readonly AllDebrid client;

    public ClientTests()
    {
        client = new AllDebrid("test-key", "https://api.test", "unit", new HttpClient(fake));
    }

    private static string Spec()
    {
        Assert.True(File.Exists(SpecPath), $"spec not found at {SpecPath}");
        return File.ReadAllText(SpecPath);
    }

    [Fact]
    public async Task ListParametersArePhpStyle()
    {
        fake.Ok("""{"magnets":[{"magnet":"a","id":1},{"magnet":"b","error":{"code":"MAGNET_INVALID_URI","message":"m"}}]}""");
        var result = await client.UploadMagnetsAsync(new[] { "a", "b" });
        Assert.Equal(HttpMethod.Post, fake.Last.Method);
        Assert.Equal("/v4/magnet/upload", fake.Last.Uri.AbsolutePath);
        Assert.Equal(new[] { "a", "b" }, fake.LastForm()["magnets[]"]);
        Assert.Equal("Bearer test-key", fake.Last.Headers["Authorization"]);
        Assert.Equal("unit", fake.Last.Headers["User-Agent"]);
        Assert.Equal(1, result[0].Id);
        Assert.Equal("MAGNET_INVALID_URI", result[1].Error?.Code);
    }

    [Fact]
    public async Task NoKeyNoAuthorization()
    {
        fake.Ok("""{"ping":"pong"}""");
        await client.WithApiKey(null).PingAsync();
        Assert.False(fake.Last.Headers.ContainsKey("Authorization"));
    }

    [Fact]
    public async Task GetParametersInQuery()
    {
        fake.Ok("""{"hosts":{}}""");
        await client.HostsAsync(hostsOnly: true);
        Assert.Equal(HttpMethod.Get, fake.Last.Method);
        Assert.Equal("?hostsOnly=true", fake.Last.Uri.Query);
    }

    [Fact]
    public async Task TorrentFilesAreMultipart()
    {
        fake.Ok("""{"files":[{"magnet":"m","id":9}]}""");
        await client.UploadTorrentFilesAsync(new[] { new TorrentFile("a.torrent", new byte[] { 100 }), new TorrentFile("b.torrent", new byte[] { 120 }) });
        Assert.StartsWith("multipart/form-data", fake.Last.ContentType);
        Assert.Equal(2, Regex.Matches(fake.Last.Body, "name=\"?files\\[\\]\"?").Count);
        Assert.Contains("a.torrent", fake.Last.Body);
        Assert.Contains("b.torrent", fake.Last.Body);
    }

    public static TheoryData<string, Func<AllDebrid, Task>, string, string> Routes() => new()
    {
        { "GET /v4/ping", c => c.PingAsync(), "", "{}" },
        { "GET /v4.1/pin/get", c => c.GetPinAsync(), "", """{"pin":"P","check":"C","expires_in":600,"user_url":"u","base_url":"b"}""" },
        { "POST /v4/pin/check", c => c.CheckPinAsync("P", "C"), "pin=P&check=C", """{"activated":false,"expires_in":600}""" },
        { "GET /v4/hosts", c => c.HostsAsync(), "", """{"hosts":{}}""" },
        { "GET /v4/hosts/domains", c => c.HostDomainsAsync(), "", """{"hosts":[]}""" },
        { "GET /v4/hosts/priority", c => c.HostPriorityAsync(), "", """{"hosts":{}}""" },
        { "GET /v4/user", c => c.UserAsync(), "", UserJson },
        { "GET /v4.1/user/hosts", c => c.UserHostsAsync(), "", """{"hosts":{}}""" },
        { "POST /v4/user/verif", c => c.VerificationStatusAsync("T"), "token=T", """{"verif":"waiting"}""" },
        { "POST /v4/user/verif/resend", c => c.ResendVerificationAsync("T"), "token=T", "{}" },
        { "POST /v4/user/notification/clear", c => c.ClearNotificationAsync("N"), "code=N", "{}" },
        { "POST /v4/link/infos", c => c.LinkInfosAsync(new[] { "a", "b" }, "p"), "link%5B%5D=a&link%5B%5D=b&password=p", """{"infos":[]}""" },
        { "POST /v4/link/redirector", c => c.ExtractRedirectorAsync("L"), "link=L", """{"links":[]}""" },
        { "POST /v4/link/unlock", c => c.UnlockAsync("L"), "link=L", "{}" },
        { "POST /v4/link/streaming", c => c.StreamingLinkAsync("G", "S"), "id=G&stream=S", "{}" },
        { "POST /v4/link/delayed", c => c.DelayedLinkAsync("42"), "id=42", """{"status":1,"time_left":5}""" },
        { "POST /v4/magnet/upload", c => c.UploadMagnetsAsync(new[] { "M" }), "magnets%5B%5D=M", """{"magnets":[]}""" },
        { "POST /v4/magnet/upload/file", c => c.UploadTorrentFilesAsync(new[] { new TorrentFile("a", new byte[] { 1 }) }), "", """{"files":[]}""" },
        { "POST /v4.1/magnet/status", c => c.MagnetsAsync(), "", """{"magnets":[]}""" },
        { "POST /v4/magnet/files", c => c.MagnetFilesAsync(new long[] { 1, 2 }), "id%5B%5D=1&id%5B%5D=2", """{"magnets":[]}""" },
        { "POST /v4/magnet/delete", c => c.DeleteMagnetAsync(5), "id=5", "{}" },
        { "POST /v4/magnet/restart", c => c.RestartMagnetsAsync(new long[] { 5, 6 }), "ids%5B%5D=5&ids%5B%5D=6", "{}" },
        { "GET /v4/user/links", c => c.SavedLinksAsync(), "", """{"links":[]}""" },
        { "POST /v4/user/links/save", c => c.SaveLinksAsync(new[] { "L" }), "links%5B%5D=L", "{}" },
        { "POST /v4/user/links/delete", c => c.DeleteSavedLinksAsync(new[] { "L" }), "links%5B%5D=L", "{}" },
        { "GET /v4/user/history", c => c.HistoryAsync(), "", """{"links":[]}""" },
        { "POST /v4/user/history/delete", c => c.PurgeHistoryAsync(), "", "{}" },
        { "GET /v4/voucher/balance", c => c.VoucherBalanceAsync(), "", """{"balance":1.5}""" },
        { "POST /v4/voucher/get", c => c.VouchersAsync(30, 2), "duration=30&nb=2", """{"codes":[]}""" },
        { "POST /v4/voucher/generate", c => c.GenerateVouchersAsync(90, 1), "duration=90&nb=1", """{"codes":[]}""" },
    };

    [Theory]
    [MemberData(nameof(Routes))]
    public async Task Routing(string want, Func<AllDebrid, Task> call, string form, string data)
    {
        ArgumentNullException.ThrowIfNull(call);
        fake.Ok(data);
        await call(client);
        var got = $"{fake.Last.Method} {fake.Last.Uri.AbsolutePath}";
        Assert.Equal(want, got);
        Assert.Contains(got, SpecOperations());
        if (form.Length > 0)
        {
            Assert.Equal(form, fake.Last.Body);
        }
    }

    [Fact]
    public void EveryOperationIsRouted()
    {
        var routed = Routes().Select(r => (string)r.Data.Item1).OrderBy(s => s, StringComparer.Ordinal).ToList();
        Assert.Equal(SpecOperations().OrderBy(s => s, StringComparer.Ordinal).ToList(), routed);
    }

    [Fact]
    public void ErrorCodesMatchTheSpec()
    {
        var block = Spec().Split("\n    ErrorCode:\n")[1];
        var codes = block.Split('\n')
            .TakeWhile(l => !l.StartsWith("    ", StringComparison.Ordinal) || l.StartsWith("     ", StringComparison.Ordinal))
            .Where(l => l.StartsWith("        - ", StringComparison.Ordinal))
            .Select(l => l["        - ".Length..].Trim('"'))
            .OrderBy(s => s, StringComparer.Ordinal);
        Assert.Equal(codes, ErrorCodes.All.OrderBy(s => s, StringComparer.Ordinal));
    }

    private static List<string> SpecOperations()
    {
        var ops = new List<string>();
        string? path = null;
        foreach (var line in Spec().Split('\n'))
        {
            if (line == "components:")
            {
                break;
            }
            var p = Regex.Match(line, @"^  (/\S+):$");
            var m = Regex.Match(line, "^    (get|post):$");
            if (p.Success)
            {
                path = p.Groups[1].Value;
            }
            else if (m.Success && path is not null)
            {
                ops.Add($"{m.Groups[1].Value.ToUpperInvariant()} {path}");
            }
        }
        return ops;
    }

    [Fact]
    public async Task ApiErrorsCarryTheirCode()
    {
        fake.Answer(200, """{"status":"error","error":{"code":"AUTH_BLOCKED","message":"m","token":"tok"}}""");
        var e = await Assert.ThrowsAsync<AllDebridException>(() => client.UserAsync());
        Assert.Equal(ErrorCodes.AuthBlocked, e.Code);
        Assert.Equal("tok", e.Token);
    }

    [Fact]
    public async Task HttpErrors()
    {
        fake.Answer(429, "Too Many Requests");
        fake.Answer(404, """{"status":"error","error":{"code":"404","message":"Endpoint doesn't exist"}}""");
        fake.Answer(200, "<html>");
        Assert.True((await Assert.ThrowsAsync<AllDebridHttpException>(() => client.PingAsync())).IsRateLimited);
        Assert.Equal(ErrorCodes.NotFound, (await Assert.ThrowsAsync<AllDebridHttpException>(() => client.PingAsync())).ApiError?.Code);
        await Assert.ThrowsAsync<System.Text.Json.JsonException>(() => client.PingAsync());
    }

    [Fact]
    public async Task FlexibleShapes()
    {
        fake.Ok("""{"magnets":{"id":7,"statusCode":4,"files":[{"n":"a","s":1,"l":"L"}]}}""");
        var magnet = await client.MagnetAsync(7);
        Assert.Equal(7, magnet.Id);
        Assert.True(magnet.IsReady());
        Assert.Equal("L", magnet.Files![0].L);

        fake.Ok(UserJson);
        Assert.Empty((await client.UserAsync()).LimitedHostersQuotas);

        fake.Ok("""{"link":"","delayed":123,"streaming":[{"id":"mp3","quality":"mp3"},{"id":"hd","quality":1080}]}""");
        var unlocked = await client.UnlockAsync("L");
        Assert.Equal("123", unlocked.Delayed?.Value);
        Assert.Equal(new[] { "mp3", "1080" }, unlocked.Streaming!.Select(s => s.Quality?.Value));

        fake.Ok("""{"hosts":{"a":{"name":"a","type":"free","domains":[],"regexp":"x"}}}""");
        Assert.Equal(new[] { "x" }, (await client.HostsAsync()).Hosts["a"].Regexp);
    }

    [Fact]
    public async Task MissingRequiredFieldsAreRejected()
    {
        fake.Ok("""{"user":{"username":"u"}}""");
        await Assert.ThrowsAsync<System.Text.Json.JsonException>(() => client.UserAsync());
    }

    [Fact]
    public async Task WaitForPin()
    {
        fake.Ok("""{"activated":false,"expires_in":600}""");
        fake.Ok("""{"activated":true,"expires_in":590,"apikey":"K"}""");
        var pin = new PinGetData { Pin = "ABCD", Check = "chk", ExpiresIn = 600, UserUrl = "u", BaseUrl = "b" };
        Assert.Equal("K", await client.WaitForPinAsync(pin, TimeSpan.FromMilliseconds(1)));
        Assert.Equal("pin=ABCD&check=chk", fake.Last.Body);
    }

    [Fact]
    public async Task WaitDelayedLink()
    {
        fake.Ok("""{"status":1,"time_left":5}""");
        fake.Ok("""{"status":2,"time_left":0,"link":"L"}""");
        fake.Ok("""{"status":3,"time_left":0}""");
        Assert.Equal("L", await client.WaitDelayedLinkAsync("1", TimeSpan.FromMilliseconds(1)));
        await Assert.ThrowsAsync<InvalidOperationException>(() => client.WaitDelayedLinkAsync("1", TimeSpan.FromMilliseconds(1)));
    }

    [Fact]
    public async Task LiveMagnets()
    {
        fake.Ok("""{"fullsync":true,"counter":1,"magnets":[{"id":1,"filename":"a","uploadDate":10,"downloaded":5},{"id":2,"filename":"b","uploadDate":20}]}""");
        fake.Ok("""{"counter":2,"magnets":{"id":1,"downloaded":9}}""");
        fake.Ok("""{"counter":3,"magnets":[{"id":2,"deleted":true},{"id":3,"filename":"c","uploadDate":30}]}""");
        fake.Ok("""{"counter":4,"magnets":[]}""");
        fake.Ok("""{"fullsync":true,"counter":1,"magnets":[{"id":9,"uploadDate":1}]}""");
        var live = client.LiveMagnets();

        Assert.Equal(new long[] { 2, 1 }, (await live.PollAsync()).Select(m => m.Id));
        var first = fake.LastForm();
        Assert.Equal("0", first["counter"][0]);

        var merged = await live.PollAsync();
        var second = fake.LastForm();
        Assert.Equal("1", second["counter"][0]);
        Assert.Equal(first["session"], second["session"]);
        var one = merged.Single(m => m.Id == 1);
        Assert.Equal(("a", (long?)9), (one.Filename, one.Downloaded));

        Assert.Equal(new long[] { 3, 1 }, (await live.PollAsync()).Select(m => m.Id));
        Assert.Equal(new long[] { 3, 1 }, (await live.PollAsync()).Select(m => m.Id));
        Assert.Equal(new long[] { 9 }, (await live.PollAsync()).Select(m => m.Id));
    }
}
