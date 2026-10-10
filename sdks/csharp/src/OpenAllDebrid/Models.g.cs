// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.
#nullable enable

using System.Collections.Generic;
using System.Text.Json.Serialization;

namespace OpenAllDebrid;

public sealed record ApiError
{
    /// <summary>
    /// A machine-readable error code; see `ErrorCode`.
    /// </summary>
    [JsonPropertyName("code")]
    public required string Code { get; init; }

    /// <summary>
    /// A human-readable message, in English.
    /// </summary>
    [JsonPropertyName("message")]
    public required string Message { get; init; }

    /// <summary>
    /// With `AUTH_BLOCKED`, the token to poll on `/v4/user/verif`.
    /// </summary>
    [JsonPropertyName("token")]
    public string? Token { get; init; }
}

public sealed record PingData
{
    [JsonPropertyName("ping")]
    public required string Ping { get; init; }

    /// <summary>
    /// API major version.
    /// </summary>
    [JsonPropertyName("version")]
    public string? Version { get; init; }

    /// <summary>
    /// Only on `/v4.1/ping`.
    /// </summary>
    [JsonPropertyName("subversion")]
    public bool? Subversion { get; init; }
}

public sealed record MessageData
{
    /// <summary>
    /// A confirmation message.
    /// </summary>
    [JsonPropertyName("message")]
    public required string Message { get; init; }
}

public sealed record PinGetData
{
    /// <summary>
    /// The code to show to the user.
    /// </summary>
    [JsonPropertyName("pin")]
    public required string Pin { get; init; }

    /// <summary>
    /// The token to send to `/v4/pin/check`.
    /// </summary>
    [JsonPropertyName("check")]
    public required string Check { get; init; }

    /// <summary>
    /// Seconds before the PIN expires.
    /// </summary>
    [JsonPropertyName("expires_in")]
    public required int ExpiresIn { get; init; }

    /// <summary>
    /// The page where the user enters the PIN, with the PIN filled in.
    /// </summary>
    [JsonPropertyName("user_url")]
    public required string UserUrl { get; init; }

    /// <summary>
    /// The page where the user enters the PIN.
    /// </summary>
    [JsonPropertyName("base_url")]
    public required string BaseUrl { get; init; }
}

public sealed record PinCheckData
{
    /// <summary>
    /// `true` once the user has entered the PIN.
    /// </summary>
    [JsonPropertyName("activated")]
    public required bool Activated { get; init; }

    /// <summary>
    /// Seconds before the PIN expires.
    /// </summary>
    [JsonPropertyName("expires_in")]
    public required int ExpiresIn { get; init; }

    /// <summary>
    /// The new API key, once `activated` is `true`.
    /// </summary>
    [JsonPropertyName("apikey")]
    public string? Apikey { get; init; }
}

/// <summary>
/// A supported file host, streaming site or redirector.
/// </summary>
public sealed record Service
{
    [JsonPropertyName("name")]
    public required string Name { get; init; }

    /// <summary>
    /// `premium` (needs a premium account) or `free`.
    /// </summary>
    [JsonPropertyName("type")]
    public required string Type { get; init; }

    /// <summary>
    /// Every domain of the service.
    /// </summary>
    [JsonPropertyName("domains")]
    public required IReadOnlyList<string> Domains { get; init; }

    /// <summary>
    /// Regular expressions matching the supported links.
    /// </summary>
    [JsonPropertyName("regexps")]
    public IReadOnlyList<string>? Regexps { get; init; }

    /// <summary>
    /// The `regexps` joined into one expression (sometimes a list).
    /// </summary>
    [JsonPropertyName("regexp")]
    [JsonConverter(typeof(StringsConverter))]
    public IReadOnlyList<string>? Regexp { get; init; }

    /// <summary>
    /// Whether the service worked in the last 5 minutes. Only some hosts are
    /// checked, about every 10 minutes.
    /// </summary>
    [JsonPropertyName("status")]
    public bool? Status { get; init; }

    /// <summary>
    /// A display name, when it differs from `name`.
    /// </summary>
    [JsonPropertyName("display")]
    public string? Display { get; init; }

    /// <summary>
    /// Domains that redirect to this host.
    /// </summary>
    [JsonPropertyName("hardRedirect")]
    public IReadOnlyList<string>? HardRedirect { get; init; }

    /// <summary>
    /// Redirectors only. Extraction is slow or rate limited.
    /// </summary>
    [JsonPropertyName("heavy")]
    public bool? Heavy { get; init; }

    /// <summary>
    /// `/v4.1/user/hosts` only. The remaining quota (see `quotaType`).
    /// </summary>
    [JsonPropertyName("quota")]
    public long? Quota { get; init; }

    /// <summary>
    /// `/v4.1/user/hosts` only. The daily quota.
    /// </summary>
    [JsonPropertyName("quotaMax")]
    public long? QuotaMax { get; init; }

    /// <summary>
    /// `/v4.1/user/hosts` only. `traffic` (quota in MB) or `nb_download`
    /// (number of downloads).
    /// </summary>
    [JsonPropertyName("quotaType")]
    public string? QuotaType { get; init; }

    /// <summary>
    /// `/v4.1/user/hosts` only. Free simultaneous download slots.
    /// </summary>
    [JsonPropertyName("limitSimuDl")]
    public int? LimitSimuDl { get; init; }
}

public sealed record HostsData
{
    [JsonPropertyName("hosts")]
    public required IReadOnlyDictionary<string, Service> Hosts { get; init; }

    [JsonPropertyName("streams")]
    public IReadOnlyDictionary<string, Service>? Streams { get; init; }

    [JsonPropertyName("redirectors")]
    public IReadOnlyDictionary<string, Service>? Redirectors { get; init; }
}

public sealed record HostDomainsData
{
    [JsonPropertyName("hosts")]
    public required IReadOnlyList<string> Hosts { get; init; }

    [JsonPropertyName("streams")]
    public IReadOnlyList<string>? Streams { get; init; }

    [JsonPropertyName("redirectors")]
    public IReadOnlyList<string>? Redirectors { get; init; }
}

public sealed record HostPriorityData
{
    /// <summary>
    /// Host names mapped to their rank, 1 being the most open.
    /// </summary>
    [JsonPropertyName("hosts")]
    public required IReadOnlyDictionary<string, int> Hosts { get; init; }
}

public sealed record User
{
    [JsonPropertyName("username")]
    public required string Username { get; init; }

    [JsonPropertyName("email")]
    public required string Email { get; init; }

    [JsonPropertyName("isPremium")]
    public required bool IsPremium { get; init; }

    /// <summary>
    /// Whether the user has an active recurring subscription.
    /// </summary>
    [JsonPropertyName("isSubscribed")]
    public required bool IsSubscribed { get; init; }

    /// <summary>
    /// Whether the account is in its free trial.
    /// </summary>
    [JsonPropertyName("isTrial")]
    public required bool IsTrial { get; init; }

    /// <summary>
    /// Unix time when premium ends, or 0 when not premium.
    /// </summary>
    [JsonPropertyName("premiumUntil")]
    public required long PremiumUntil { get; init; }

    /// <summary>
    /// Interface language, e.g. `en` or `fr`.
    /// </summary>
    [JsonPropertyName("lang")]
    public required string Lang { get; init; }

    /// <summary>
    /// Preferred alldebrid.* top-level domain, e.g. `com` or `fr`.
    /// </summary>
    [JsonPropertyName("preferedDomain")]
    public required string PreferedDomain { get; init; }

    [JsonPropertyName("fidelityPoints")]
    public required int FidelityPoints { get; init; }

    [JsonPropertyName("limitedHostersQuotas")]
    [JsonConverter(typeof(QuotasConverter))]
    public required IReadOnlyDictionary<string, long> LimitedHostersQuotas { get; init; }

    /// <summary>
    /// During the trial, the remaining traffic in MB.
    /// </summary>
    [JsonPropertyName("remainingTrialQuota")]
    public long? RemainingTrialQuota { get; init; }

    /// <summary>
    /// Codes of the current notifications; clear them with `/v4/user/notification/clear`.
    /// </summary>
    [JsonPropertyName("notifications")]
    public required IReadOnlyList<string> Notifications { get; init; }

    /// <summary>
    /// Whether recent links history logging is enabled.
    /// </summary>
    [JsonPropertyName("history")]
    public bool? History { get; init; }
}

public sealed record UserData
{
    [JsonPropertyName("user")]
    public required User User { get; init; }
}

public sealed record UserHostsData
{
    [JsonPropertyName("hosts")]
    public required IReadOnlyDictionary<string, Service> Hosts { get; init; }
}

public sealed record VerificationData
{
    [JsonPropertyName("verif")]
    public required string Verif { get; init; }

    /// <summary>
    /// When `waiting`, whether `/v4/user/verif/resend` may be called.
    /// </summary>
    [JsonPropertyName("resendable")]
    public bool? Resendable { get; init; }

    /// <summary>
    /// When `allowed`, the API key.
    /// </summary>
    [JsonPropertyName("apikey")]
    public string? Apikey { get; init; }
}

public sealed record ResendVerificationData
{
    [JsonPropertyName("sent")]
    public required bool Sent { get; init; }
}

public sealed record LinkInfo
{
    /// <summary>
    /// The requested link.
    /// </summary>
    [JsonPropertyName("link")]
    public required string Link { get; init; }

    [JsonPropertyName("filename")]
    public string? Filename { get; init; }

    /// <summary>
    /// File size in bytes.
    /// </summary>
    [JsonPropertyName("size")]
    public long? Size { get; init; }

    [JsonPropertyName("host")]
    public string? Host { get; init; }

    /// <summary>
    /// The host's main domain.
    /// </summary>
    [JsonPropertyName("hostDomain")]
    public string? HostDomain { get; init; }

    [JsonPropertyName("error")]
    public ApiError? Error { get; init; }
}

public sealed record LinkInfosData
{
    [JsonPropertyName("infos")]
    public required IReadOnlyList<LinkInfo> Infos { get; init; }
}

public sealed record RedirectorData
{
    /// <summary>
    /// Encrypted links, usable with `/v4/link/infos` and `/v4/link/unlock`.
    /// </summary>
    [JsonPropertyName("links")]
    public required IReadOnlyList<string> Links { get; init; }
}

/// <summary>
/// One quality of a streaming link.
/// </summary>
public sealed record Stream
{
    /// <summary>
    /// The `stream` value for `/v4/link/streaming`.
    /// </summary>
    [JsonPropertyName("id")]
    public required string Id { get; init; }

    [JsonPropertyName("ext")]
    public string? Ext { get; init; }

    /// <summary>
    /// Height in pixels (e.g. `1080`), or a label such as `mp3`.
    /// </summary>
    [JsonPropertyName("quality")]
    public FlexString? Quality { get; init; }

    [JsonPropertyName("filesize")]
    public long? Filesize { get; init; }

    [JsonPropertyName("proto")]
    public string? Proto { get; init; }

    [JsonPropertyName("name")]
    public string? Name { get; init; }

    /// <summary>
    /// Audio-only streams. Total bitrate.
    /// </summary>
    [JsonPropertyName("tb")]
    public double? Tb { get; init; }

    /// <summary>
    /// Audio-only streams. Audio bitrate.
    /// </summary>
    [JsonPropertyName("abr")]
    public int? Abr { get; init; }
}

public sealed record UnlockData
{
    /// <summary>
    /// The direct download link. Empty when `delayed` is set.
    /// </summary>
    [JsonPropertyName("link")]
    public string? Link { get; init; }

    [JsonPropertyName("host")]
    public string? Host { get; init; }

    [JsonPropertyName("filename")]
    public string? Filename { get; init; }

    /// <summary>
    /// File size in bytes.
    /// </summary>
    [JsonPropertyName("filesize")]
    public long? Filesize { get; init; }

    /// <summary>
    /// The generation ID, needed by `/v4/link/streaming`.
    /// </summary>
    [JsonPropertyName("id")]
    public string? Id { get; init; }

    [JsonPropertyName("hostDomain")]
    public string? HostDomain { get; init; }

    /// <summary>
    /// Unused.
    /// </summary>
    [JsonPropertyName("paws")]
    public bool? Paws { get; init; }

    /// <summary>
    /// Available qualities for streaming links; empty otherwise.
    /// </summary>
    [JsonPropertyName("streaming")]
    public IReadOnlyList<Stream>? Streaming { get; init; }

    /// <summary>
    /// Set when the link is still being generated; poll `/v4/link/delayed` with it.
    /// </summary>
    [JsonPropertyName("delayed")]
    public FlexString? Delayed { get; init; }
}

public sealed record StreamingData
{
    /// <summary>
    /// The download link; absent when `delayed` is set.
    /// </summary>
    [JsonPropertyName("link")]
    public string? Link { get; init; }

    [JsonPropertyName("filename")]
    public string? Filename { get; init; }

    [JsonPropertyName("filesize")]
    public long? Filesize { get; init; }

    /// <summary>
    /// Set when the link is still being generated; poll `/v4/link/delayed` with it.
    /// </summary>
    [JsonPropertyName("delayed")]
    public FlexString? Delayed { get; init; }
}

public sealed record DelayedData
{
    /// <summary>
    /// `1` still processing, `2` ready (`link` is set), `3` failed.
    /// </summary>
    [JsonPropertyName("status")]
    public required int Status { get; init; }

    /// <summary>
    /// Estimated seconds left.
    /// </summary>
    [JsonPropertyName("time_left")]
    public required int TimeLeft { get; init; }

    /// <summary>
    /// The download link, once `status` is `2`.
    /// </summary>
    [JsonPropertyName("link")]
    public string? Link { get; init; }
}

/// <summary>
/// An uploaded magnet or torrent, or a per-item `error`.
/// </summary>
public sealed record UploadedMagnet
{
    /// <summary>
    /// The magnet as sent (for torrent files, the magnet built from the file).
    /// </summary>
    [JsonPropertyName("magnet")]
    public string? Magnet { get; init; }

    /// <summary>
    /// Torrent file uploads with an `error` only. The file name sent.
    /// </summary>
    [JsonPropertyName("file")]
    public string? File { get; init; }

    /// <summary>
    /// The info hash.
    /// </summary>
    [JsonPropertyName("hash")]
    public string? Hash { get; init; }

    /// <summary>
    /// The torrent name, or `noname` when unknown.
    /// </summary>
    [JsonPropertyName("name")]
    public string? Name { get; init; }

    [JsonPropertyName("filename_original")]
    public string? FilenameOriginal { get; init; }

    /// <summary>
    /// Total size in bytes.
    /// </summary>
    [JsonPropertyName("size")]
    public long? Size { get; init; }

    /// <summary>
    /// Whether the files are already available.
    /// </summary>
    [JsonPropertyName("ready")]
    public bool? Ready { get; init; }

    /// <summary>
    /// The magnet ID, for `/v4.1/magnet/status`.
    /// </summary>
    [JsonPropertyName("id")]
    public long? Id { get; init; }

    [JsonPropertyName("error")]
    public ApiError? Error { get; init; }
}

public sealed record UploadMagnetsData
{
    [JsonPropertyName("magnets")]
    public required IReadOnlyList<UploadedMagnet> Magnets { get; init; }
}

public sealed record UploadTorrentFilesData
{
    [JsonPropertyName("files")]
    public required IReadOnlyList<UploadedMagnet> Files { get; init; }
}

/// <summary>
/// A file (`n`, `s`, `l`) or a folder (`n`, `e`) of a magnet.
/// </summary>
public sealed record FileNode
{
    /// <summary>
    /// File or folder name.
    /// </summary>
    [JsonPropertyName("n")]
    public required string N { get; init; }

    /// <summary>
    /// Files only. Size in bytes.
    /// </summary>
    [JsonPropertyName("s")]
    public long? S { get; init; }

    /// <summary>
    /// Files only. Link to unlock with `/v4/link/unlock`.
    /// </summary>
    [JsonPropertyName("l")]
    public string? L { get; init; }

    /// <summary>
    /// Folders only. The entries of the folder.
    /// </summary>
    [JsonPropertyName("e")]
    public IReadOnlyList<FileNode>? E { get; init; }
}

/// <summary>
/// A magnet of the account. In a live mode delta, only `id` and the
/// fields that changed are present.
/// </summary>
public sealed record Magnet
{
    [JsonPropertyName("id")]
    public required long Id { get; init; }

    [JsonPropertyName("filename")]
    public string? Filename { get; init; }

    /// <summary>
    /// Total size in bytes.
    /// </summary>
    [JsonPropertyName("size")]
    public long? Size { get; init; }

    /// <summary>
    /// The info hash.
    /// </summary>
    [JsonPropertyName("hash")]
    public string? Hash { get; init; }

    /// <summary>
    /// The status in plain English.
    /// </summary>
    [JsonPropertyName("status")]
    public string? Status { get; init; }

    /// <summary>
    /// `0` queued, `1` downloading, `2` compressing or moving, `3` uploading,
    /// `4` ready. `5` to `15` are errors: `5` upload failed, `6` unpacking error,
    /// `7` no peer, `8` too big, `9` internal error, `10` took too long,
    /// `11` files removed from the host, `12` and `13` processing failed,
    /// `14` tracker error, `15` no peer.
    /// </summary>
    [JsonPropertyName("statusCode")]
    public int? StatusCode { get; init; }

    /// <summary>
    /// Bytes downloaded so far (while processing).
    /// </summary>
    [JsonPropertyName("downloaded")]
    public long? Downloaded { get; init; }

    /// <summary>
    /// Bytes uploaded so far (while processing).
    /// </summary>
    [JsonPropertyName("uploaded")]
    public long? Uploaded { get; init; }

    /// <summary>
    /// Seeders (while downloading).
    /// </summary>
    [JsonPropertyName("seeders")]
    public int? Seeders { get; init; }

    /// <summary>
    /// Bytes per second (while downloading).
    /// </summary>
    [JsonPropertyName("downloadSpeed")]
    public long? DownloadSpeed { get; init; }

    /// <summary>
    /// Bytes per second (while uploading).
    /// </summary>
    [JsonPropertyName("uploadSpeed")]
    public long? UploadSpeed { get; init; }

    /// <summary>
    /// Progress percentage while processing.
    /// </summary>
    [JsonPropertyName("processingPerc")]
    public double? ProcessingPerc { get; init; }

    /// <summary>
    /// Unix time the magnet was added.
    /// </summary>
    [JsonPropertyName("uploadDate")]
    public long? UploadDate { get; init; }

    /// <summary>
    /// Unix time the magnet became ready, or 0.
    /// </summary>
    [JsonPropertyName("completionDate")]
    public long? CompletionDate { get; init; }

    /// <summary>
    /// `m` when added as a magnet, `b` when added as a torrent file.
    /// </summary>
    [JsonPropertyName("type")]
    public string? Type { get; init; }

    [JsonPropertyName("notified")]
    public bool? Notified { get; init; }

    [JsonPropertyName("version")]
    public int? Version { get; init; }

    /// <summary>
    /// Number of files, once ready.
    /// </summary>
    [JsonPropertyName("nbLinks")]
    public int? NbLinks { get; init; }

    /// <summary>
    /// Only when a single magnet is requested with `id`.
    /// </summary>
    [JsonPropertyName("files")]
    public IReadOnlyList<FileNode>? Files { get; init; }

    /// <summary>
    /// Live mode deltas only. The magnet was deleted; drop it from the local state.
    /// </summary>
    [JsonPropertyName("deleted")]
    public bool? Deleted { get; init; }
}

public sealed record MagnetStatusData
{
    [JsonPropertyName("magnets")]
    [JsonConverter(typeof(MagnetListConverter))]
    public required IReadOnlyList<Magnet> Magnets { get; init; }

    /// <summary>
    /// Live mode. The counter to send on the next call.
    /// </summary>
    [JsonPropertyName("counter")]
    public long? Counter { get; init; }

    /// <summary>
    /// Live mode. `true` when the response holds the full list.
    /// </summary>
    [JsonPropertyName("fullsync")]
    public bool? Fullsync { get; init; }
}

public sealed record MagnetFiles
{
    /// <summary>
    /// The magnet ID (sent back as a string).
    /// </summary>
    [JsonPropertyName("id")]
    public required FlexString Id { get; init; }

    [JsonPropertyName("files")]
    public IReadOnlyList<FileNode>? Files { get; init; }

    [JsonPropertyName("error")]
    public ApiError? Error { get; init; }
}

public sealed record MagnetFilesData
{
    [JsonPropertyName("magnets")]
    public required IReadOnlyList<MagnetFiles> Magnets { get; init; }
}

public sealed record RestartedMagnet
{
    /// <summary>
    /// The magnet ID.
    /// </summary>
    [JsonPropertyName("magnet")]
    public string? Magnet { get; init; }

    [JsonPropertyName("message")]
    public string? Message { get; init; }

    [JsonPropertyName("error")]
    public ApiError? Error { get; init; }
}

public sealed record RestartMagnetsData
{
    /// <summary>
    /// With `id`. The restarted magnet.
    /// </summary>
    [JsonPropertyName("id")]
    public long? Id { get; init; }

    /// <summary>
    /// With `id`. A confirmation message.
    /// </summary>
    [JsonPropertyName("message")]
    public string? Message { get; init; }

    /// <summary>
    /// With `ids[]`. One entry per magnet.
    /// </summary>
    [JsonPropertyName("magnets")]
    public IReadOnlyList<RestartedMagnet>? Magnets { get; init; }
}

public sealed record SavedLink
{
    [JsonPropertyName("link")]
    public required string Link { get; init; }

    [JsonPropertyName("filename")]
    public string? Filename { get; init; }

    [JsonPropertyName("size")]
    public long? Size { get; init; }

    /// <summary>
    /// Unix time the link was saved or unlocked.
    /// </summary>
    [JsonPropertyName("date")]
    public long? Date { get; init; }

    [JsonPropertyName("host")]
    public string? Host { get; init; }

    /// <summary>
    /// History only. The download link that was generated.
    /// </summary>
    [JsonPropertyName("link_dl")]
    public string? LinkDl { get; init; }

    /// <summary>
    /// History only. The error the unlock gave, if any.
    /// </summary>
    [JsonPropertyName("error")]
    public string? Error { get; init; }
}

public sealed record SavedLinksData
{
    [JsonPropertyName("links")]
    public required IReadOnlyList<SavedLink> Links { get; init; }
}

public sealed record VoucherBalanceData
{
    /// <summary>
    /// The reseller balance.
    /// </summary>
    [JsonPropertyName("balance")]
    public required double Balance { get; init; }
}

public sealed record VouchersData
{
    [JsonPropertyName("codes")]
    public required IReadOnlyList<string> Codes { get; init; }

    /// <summary>
    /// `true` when fewer vouchers than requested were available.
    /// </summary>
    [JsonPropertyName("partialList")]
    public bool? PartialList { get; init; }
}

public sealed record GeneratedVouchersData
{
    [JsonPropertyName("codes")]
    public required IReadOnlyList<string> Codes { get; init; }

    [JsonPropertyName("pricePerVoucher")]
    public double? PricePerVoucher { get; init; }

    /// <summary>
    /// Total price of the request.
    /// </summary>
    [JsonPropertyName("total")]
    public double? Total { get; init; }

    /// <summary>
    /// The balance left after the generation.
    /// </summary>
    [JsonPropertyName("balance")]
    public double? Balance { get; init; }
}

/// <summary>Error codes the API is known to return. The API may return others.</summary>
public static class ErrorCodes
{
    public const string Generic = "GENERIC";
    public const string NotFound = "404";
    public const string Maintenance = "MAINTENANCE";
    public const string Discontinued = "DISCONTINUED";
    public const string AuthMissingApikey = "AUTH_MISSING_APIKEY";
    public const string AuthBadApikey = "AUTH_BAD_APIKEY";
    public const string AuthBlocked = "AUTH_BLOCKED";
    public const string AuthUserBanned = "AUTH_USER_BANNED";
    public const string AlreadySent = "ALREADY_SENT";
    public const string NoServer = "NO_SERVER";
    public const string LinkIsMissing = "LINK_IS_MISSING";
    public const string BadLink = "BAD_LINK";
    public const string LinkHostNotSupported = "LINK_HOST_NOT_SUPPORTED";
    public const string LinkDown = "LINK_DOWN";
    public const string LinkPassProtected = "LINK_PASS_PROTECTED";
    public const string LinkHostUnavailable = "LINK_HOST_UNAVAILABLE";
    public const string LinkTooManyDownloads = "LINK_TOO_MANY_DOWNLOADS";
    public const string LinkHostFull = "LINK_HOST_FULL";
    public const string LinkHostLimitReached = "LINK_HOST_LIMIT_REACHED";
    public const string LinkError = "LINK_ERROR";
    public const string LinkTemporaryUnavailable = "LINK_TEMPORARY_UNAVAILABLE";
    public const string LinkNotSupported = "LINK_NOT_SUPPORTED";
    public const string RedirectorNotSupported = "REDIRECTOR_NOT_SUPPORTED";
    public const string RedirectorError = "REDIRECTOR_ERROR";
    public const string RedirectorRateLimited = "REDIRECTOR_RATE_LIMITED";
    public const string StreamInvalidGenId = "STREAM_INVALID_GEN_ID";
    public const string StreamInvalidStreamId = "STREAM_INVALID_STREAM_ID";
    public const string DelayedInvalidId = "DELAYED_INVALID_ID";
    public const string FreeTrialLimitReached = "FREE_TRIAL_LIMIT_REACHED";
    public const string MustBePremium = "MUST_BE_PREMIUM";
    public const string MagnetInvalidId = "MAGNET_INVALID_ID";
    public const string MagnetInvalidUri = "MAGNET_INVALID_URI";
    public const string MagnetInvalidFile = "MAGNET_INVALID_FILE";
    public const string MagnetFileUploadFailed = "MAGNET_FILE_UPLOAD_FAILED";
    public const string MagnetNoUri = "MAGNET_NO_URI";
    public const string MagnetProcessing = "MAGNET_PROCESSING";
    public const string MagnetTooManyActive = "MAGNET_TOO_MANY_ACTIVE";
    public const string MagnetTooMany = "MAGNET_TOO_MANY";
    public const string MagnetMustBePremium = "MAGNET_MUST_BE_PREMIUM";
    public const string MagnetNoServer = "MAGNET_NO_SERVER";
    public const string MagnetTooLarge = "MAGNET_TOO_LARGE";
    public const string MagnetUploadFailed = "MAGNET_UPLOAD_FAILED";
    public const string MagnetInternalError = "MAGNET_INTERNAL_ERROR";
    public const string MagnetCantBootstrap = "MAGNET_CANT_BOOTSTRAP";
    public const string MagnetMagnetTooBig = "MAGNET_MAGNET_TOO_BIG";
    public const string MagnetTookTooLong = "MAGNET_TOOK_TOO_LONG";
    public const string MagnetLinksRemoved = "MAGNET_LINKS_REMOVED";
    public const string MagnetProcessingFailed = "MAGNET_PROCESSING_FAILED";
    public const string PinAlreadyAuthed = "PIN_ALREADY_AUTHED";
    public const string PinExpired = "PIN_EXPIRED";
    public const string PinInvalid = "PIN_INVALID";
    public const string UserLinkMissing = "USER_LINK_MISSING";
    public const string UserLinkInvalid = "USER_LINK_INVALID";
    public const string MissingNotifEndpoint = "MISSING_NOTIF_ENDPOINT";
    public const string VoucherDurationInvalid = "VOUCHER_DURATION_INVALID";
    public const string VoucherNbInvalid = "VOUCHER_NB_INVALID";
    public const string NoMoreVoucher = "NO_MORE_VOUCHER";
    public const string InsufficientBalance = "INSUFFICIENT_BALANCE";
    public const string DownloadFailed = "DOWNLOAD_FAILED";
    public const string AccountInvalid = "ACCOUNT_INVALID";
    public const string NoJsonParam = "NO_JSON_PARAM";
    public const string JsonInvalid = "JSON_INVALID";
    public const string FreedaysInvalidCountry = "FREEDAYS_INVALID_COUNTRY";
    public const string FreedaysInvalidPhone = "FREEDAYS_INVALID_PHONE";
    public const string FreedaysInvalidProvider = "FREEDAYS_INVALID_PROVIDER";
    public const string FreedaysUsedPhone = "FREEDAYS_USED_PHONE";
    public const string FreedaysAlreadyRequested = "FREEDAYS_ALREADY_REQUESTED";
    public const string FreedaysInvalidStatus = "FREEDAYS_INVALID_STATUS";
    public const string FreedaysTooMuchRetries = "FREEDAYS_TOO_MUCH_RETRIES";
    public const string FreedaysInvalidCode = "FREEDAYS_INVALID_CODE";
    public const string FreedaysDenied = "FREEDAYS_DENIED";
    public const string FreedaysErrorSending = "FREEDAYS_ERROR_SENDING";

    /// <summary>Every code of the spec.</summary>
    public static readonly IReadOnlyList<string> All = new[]
    {
        Generic,
        NotFound,
        Maintenance,
        Discontinued,
        AuthMissingApikey,
        AuthBadApikey,
        AuthBlocked,
        AuthUserBanned,
        AlreadySent,
        NoServer,
        LinkIsMissing,
        BadLink,
        LinkHostNotSupported,
        LinkDown,
        LinkPassProtected,
        LinkHostUnavailable,
        LinkTooManyDownloads,
        LinkHostFull,
        LinkHostLimitReached,
        LinkError,
        LinkTemporaryUnavailable,
        LinkNotSupported,
        RedirectorNotSupported,
        RedirectorError,
        RedirectorRateLimited,
        StreamInvalidGenId,
        StreamInvalidStreamId,
        DelayedInvalidId,
        FreeTrialLimitReached,
        MustBePremium,
        MagnetInvalidId,
        MagnetInvalidUri,
        MagnetInvalidFile,
        MagnetFileUploadFailed,
        MagnetNoUri,
        MagnetProcessing,
        MagnetTooManyActive,
        MagnetTooMany,
        MagnetMustBePremium,
        MagnetNoServer,
        MagnetTooLarge,
        MagnetUploadFailed,
        MagnetInternalError,
        MagnetCantBootstrap,
        MagnetMagnetTooBig,
        MagnetTookTooLong,
        MagnetLinksRemoved,
        MagnetProcessingFailed,
        PinAlreadyAuthed,
        PinExpired,
        PinInvalid,
        UserLinkMissing,
        UserLinkInvalid,
        MissingNotifEndpoint,
        VoucherDurationInvalid,
        VoucherNbInvalid,
        NoMoreVoucher,
        InsufficientBalance,
        DownloadFailed,
        AccountInvalid,
        NoJsonParam,
        JsonInvalid,
        FreedaysInvalidCountry,
        FreedaysInvalidPhone,
        FreedaysInvalidProvider,
        FreedaysUsedPhone,
        FreedaysAlreadyRequested,
        FreedaysInvalidStatus,
        FreedaysTooMuchRetries,
        FreedaysInvalidCode,
        FreedaysDenied,
        FreedaysErrorSending,
    };
}
