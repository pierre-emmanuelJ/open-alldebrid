using Xunit;

namespace OpenAllDebrid.Tests;

/// <summary>Calls the real API with the public demo keys. Set OFFLINE=1 to skip.</summary>
public sealed class DemoTests
{
    private readonly AllDebrid client = new(AllDebrid.DemoPremiumKey, userAgent: "open-alldebrid-csharp-tests");
    private readonly AllDebrid anonymous = new(userAgent: "open-alldebrid-csharp-tests");

    private static bool Offline => Environment.GetEnvironmentVariable("OFFLINE") == "1";

    [Fact]
    public async Task PublicEndpoints()
    {
        Assert.SkipWhen(Offline, "OFFLINE=1");
        await anonymous.PingAsync();
        Assert.NotEmpty((await anonymous.HostsAsync()).Hosts);
        Assert.NotEmpty((await anonymous.HostDomainsAsync()).Hosts);
        Assert.NotEmpty(await anonymous.HostPriorityAsync());
    }

    [Fact]
    public async Task Pin()
    {
        Assert.SkipWhen(Offline, "OFFLINE=1");
        Assert.StartsWith("https://", (await anonymous.GetPinAsync()).UserUrl, StringComparison.Ordinal);
        Assert.True((await anonymous.CheckPinAsync("123B", "664c3ca2635c99f291d28e11ea18e154750bd21a")).Activated);
    }

    [Fact]
    public async Task User()
    {
        Assert.SkipWhen(Offline, "OFFLINE=1");
        Assert.Equal("demoUserPremium", (await client.UserAsync()).Username);
        Assert.NotEmpty(await client.UserHostsAsync());
        await client.ClearNotificationAsync("NOTIF_CODE");
        Assert.Equal("allowed", (await anonymous.VerificationStatusAsync("verificationTokenAllowed3f9eb0db")).Verif);
        var e = await Assert.ThrowsAsync<AllDebridException>(() => anonymous.ResendVerificationAsync("verificationTokenAlreadySentb0db"));
        Assert.Equal(ErrorCodes.AlreadySent, e.Code);
    }

    [Fact]
    public async Task Links()
    {
        Assert.SkipWhen(Offline, "OFFLINE=1");
        var infos = await client.LinkInfosAsync(new[] { "https://example.com/file2MB.txt", "https://example.com/down" });
        Assert.False(string.IsNullOrEmpty(infos[0].Filename));
        Assert.NotEmpty(await client.ExtractRedirectorAsync("https://example.com/redir/folder"));
        Assert.NotEmpty((await client.UnlockAsync("https://example.com/streaming")).Streaming!);
        Assert.NotNull((await client.UnlockAsync("https://example.com/delayed")).Delayed);
        var down = await Assert.ThrowsAsync<AllDebridException>(() => client.UnlockAsync("https://example.com/down"));
        Assert.Equal(ErrorCodes.LinkDown, down.Code);
        Assert.False(string.IsNullOrEmpty((await client.StreamingLinkAsync("demov1bc797", "hd-720")).Link));
        Assert.Equal(DelayedStatus.Ready, (await client.DelayedLinkAsync("1234567")).Status);
    }

    [Fact]
    public async Task Magnets()
    {
        Assert.SkipWhen(Offline, "OFFLINE=1");
        var uploaded = await client.UploadMagnetsAsync(new[] { "magnet:?xt=urn:btih:3A842783e3005495d5d1637f5364b59343c7844707" });
        Assert.NotNull(uploaded[0].Id);
        await client.UploadTorrentFilesAsync(new[] { new TorrentFile("demo.torrent", "d4:infod4:name4:demoee"u8.ToArray()) });
        Assert.NotEmpty(await client.MagnetsAsync());
        Assert.NotEmpty((await client.MagnetAsync(56789)).Files!);
        Assert.Single(await client.MagnetFilesAsync(new long[] { 56789 }));
        await client.DeleteMagnetAsync(56789);
        await client.RestartMagnetAsync(99999);
        Assert.NotEmpty(await client.LiveMagnets().PollAsync());
    }

    [Fact]
    public async Task SavedLinks()
    {
        Assert.SkipWhen(Offline, "OFFLINE=1");
        Assert.NotEmpty(await client.SavedLinksAsync());
        await client.SaveLinksAsync(new[] { "https://example.com/file2MB.txt" });
        await client.DeleteSavedLinksAsync(new[] { "https://example.com/file2MB.txt" });
        await client.HistoryAsync();
        await client.PurgeHistoryAsync();
    }

    [Fact]
    public async Task BadKey()
    {
        Assert.SkipWhen(Offline, "OFFLINE=1");
        var e = await Assert.ThrowsAsync<AllDebridException>(() => client.WithApiKey("not-a-real-key").UserAsync());
        Assert.Equal(ErrorCodes.AuthBadApikey, e.Code);
    }
}
