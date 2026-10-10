// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.
@file:Suppress("ktlint")

package io.github.pierreemmanuelj.alldebrid

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.JsonElement

@Serializable
public data class ApiError(
    /** A machine-readable error code; see `ErrorCode`. */
    val code: String,
    /** A human-readable message, in English. */
    val message: String,
    /** With `AUTH_BLOCKED`, the token to poll on `/v4/user/verif`. */
    val token: String? = null,
)

@Serializable
public data class PingData(
    val ping: String,
    /** API major version. */
    val version: String? = null,
    /** Only on `/v4.1/ping`. */
    val subversion: Boolean? = null,
)

@Serializable
public data class MessageData(
    /** A confirmation message. */
    val message: String,
)

@Serializable
public data class PinGetData(
    /** The code to show to the user. */
    val pin: String,
    /** The token to send to `/v4/pin/check`. */
    val check: String,
    /** Seconds before the PIN expires. */
    @SerialName("expires_in") val expiresIn: Int,
    /** The page where the user enters the PIN, with the PIN filled in. */
    @SerialName("user_url") val userUrl: String,
    /** The page where the user enters the PIN. */
    @SerialName("base_url") val baseUrl: String,
)

@Serializable
public data class PinCheckData(
    /** `true` once the user has entered the PIN. */
    val activated: Boolean,
    /** Seconds before the PIN expires. */
    @SerialName("expires_in") val expiresIn: Int,
    /** The new API key, once `activated` is `true`. */
    val apikey: String? = null,
)

/** A supported file host, streaming site or redirector. */
@Serializable
public data class Service(
    val name: String,
    /** `premium` (needs a premium account) or `free`. */
    val type: String,
    /** Every domain of the service. */
    val domains: List<String>,
    /** Regular expressions matching the supported links. */
    val regexps: List<String>? = null,
    /** The `regexps` joined into one expression (sometimes a list). */
    @Serializable(with = StringsSerializer::class) val regexp: List<String>? = null,
    /**
     * Whether the service worked in the last 5 minutes. Only some hosts are
     * checked, about every 10 minutes.
     */
    val status: Boolean? = null,
    /** A display name, when it differs from `name`. */
    val display: String? = null,
    /** Domains that redirect to this host. */
    val hardRedirect: List<String>? = null,
    /** Redirectors only. Extraction is slow or rate limited. */
    val heavy: Boolean? = null,
    /** `/v4.1/user/hosts` only. The remaining quota (see `quotaType`). */
    val quota: Long? = null,
    /** `/v4.1/user/hosts` only. The daily quota. */
    val quotaMax: Long? = null,
    /**
     * `/v4.1/user/hosts` only. `traffic` (quota in MB) or `nb_download`
     * (number of downloads).
     */
    val quotaType: String? = null,
    /** `/v4.1/user/hosts` only. Free simultaneous download slots. */
    val limitSimuDl: Int? = null,
)

/** Services keyed by name. */
public typealias ServiceMap = Map<String, Service>

@Serializable
public data class HostsData(
    val hosts: ServiceMap,
    val streams: ServiceMap? = null,
    val redirectors: ServiceMap? = null,
)

@Serializable
public data class HostDomainsData(
    val hosts: List<String>,
    val streams: List<String>? = null,
    val redirectors: List<String>? = null,
)

@Serializable
public data class HostPriorityData(
    /** Host names mapped to their rank, 1 being the most open. */
    val hosts: Map<String, Int>,
)

@Serializable
public data class User(
    val username: String,
    val email: String,
    val isPremium: Boolean,
    /** Whether the user has an active recurring subscription. */
    val isSubscribed: Boolean,
    /** Whether the account is in its free trial. */
    val isTrial: Boolean,
    /** Unix time when premium ends, or 0 when not premium. */
    val premiumUntil: Long,
    /** Interface language, e.g. `en` or `fr`. */
    val lang: String,
    /** Preferred alldebrid.* top-level domain, e.g. `com` or `fr`. */
    val preferedDomain: String,
    val fidelityPoints: Int,
    @Serializable(with = QuotasSerializer::class) val limitedHostersQuotas: Map<String, Long>,
    /** During the trial, the remaining traffic in MB. */
    val remainingTrialQuota: Long? = null,
    /** Codes of the current notifications; clear them with `/v4/user/notification/clear`. */
    val notifications: List<String>,
    /** Whether recent links history logging is enabled. */
    val history: Boolean? = null,
)

@Serializable
public data class UserData(
    val user: User,
)

@Serializable
public data class UserHostsData(
    val hosts: ServiceMap,
)

@Serializable
public data class VerificationData(
    val verif: String,
    /** When `waiting`, whether `/v4/user/verif/resend` may be called. */
    val resendable: Boolean? = null,
    /** When `allowed`, the API key. */
    val apikey: String? = null,
)

@Serializable
public data class ResendVerificationData(
    val sent: Boolean,
)

@Serializable
public data class LinkInfo(
    /** The requested link. */
    val link: String,
    val filename: String? = null,
    /** File size in bytes. */
    val size: Long? = null,
    val host: String? = null,
    /** The host's main domain. */
    val hostDomain: String? = null,
    val error: ApiError? = null,
)

@Serializable
public data class LinkInfosData(
    val infos: List<LinkInfo>,
)

@Serializable
public data class RedirectorData(
    /** Encrypted links, usable with `/v4/link/infos` and `/v4/link/unlock`. */
    val links: List<String>,
)

/** One quality of a streaming link. */
@Serializable
public data class Stream(
    /** The `stream` value for `/v4/link/streaming`. */
    val id: String,
    val ext: String? = null,
    /** Height in pixels (e.g. `1080`), or a label such as `mp3`. */
    @Serializable(with = FlexStringSerializer::class) val quality: FlexString? = null,
    val filesize: Long? = null,
    val proto: String? = null,
    val name: String? = null,
    /** Audio-only streams. Total bitrate. */
    val tb: Double? = null,
    /** Audio-only streams. Audio bitrate. */
    val abr: Int? = null,
)

@Serializable
public data class UnlockData(
    /** The direct download link. Empty when `delayed` is set. */
    val link: String? = null,
    val host: String? = null,
    val filename: String? = null,
    /** File size in bytes. */
    val filesize: Long? = null,
    /** The generation ID, needed by `/v4/link/streaming`. */
    val id: String? = null,
    val hostDomain: String? = null,
    /** Unused. */
    val paws: Boolean? = null,
    /** Available qualities for streaming links; empty otherwise. */
    val streaming: List<Stream>? = null,
    /** Set when the link is still being generated; poll `/v4/link/delayed` with it. */
    @Serializable(with = FlexStringSerializer::class) val delayed: FlexString? = null,
)

@Serializable
public data class StreamingData(
    /** The download link; absent when `delayed` is set. */
    val link: String? = null,
    val filename: String? = null,
    val filesize: Long? = null,
    /** Set when the link is still being generated; poll `/v4/link/delayed` with it. */
    @Serializable(with = FlexStringSerializer::class) val delayed: FlexString? = null,
)

@Serializable
public data class DelayedData(
    /** `1` still processing, `2` ready (`link` is set), `3` failed. */
    val status: Int,
    /** Estimated seconds left. */
    @SerialName("time_left") val timeLeft: Int,
    /** The download link, once `status` is `2`. */
    val link: String? = null,
)

/** An uploaded magnet or torrent, or a per-item `error`. */
@Serializable
public data class UploadedMagnet(
    /** The magnet as sent (for torrent files, the magnet built from the file). */
    val magnet: String? = null,
    /** Torrent file uploads with an `error` only. The file name sent. */
    val file: String? = null,
    /** The info hash. */
    val hash: String? = null,
    /** The torrent name, or `noname` when unknown. */
    val name: String? = null,
    @SerialName("filename_original") val filenameOriginal: String? = null,
    /** Total size in bytes. */
    val size: Long? = null,
    /** Whether the files are already available. */
    val ready: Boolean? = null,
    /** The magnet ID, for `/v4.1/magnet/status`. */
    val id: Long? = null,
    val error: ApiError? = null,
)

@Serializable
public data class UploadMagnetsData(
    val magnets: List<UploadedMagnet>,
)

@Serializable
public data class UploadTorrentFilesData(
    val files: List<UploadedMagnet>,
)

/** A file (`n`, `s`, `l`) or a folder (`n`, `e`) of a magnet. */
@Serializable
public data class FileNode(
    /** File or folder name. */
    val n: String,
    /** Files only. Size in bytes. */
    val s: Long? = null,
    /** Files only. Link to unlock with `/v4/link/unlock`. */
    val l: String? = null,
    /** Folders only. The entries of the folder. */
    val e: List<FileNode>? = null,
)

/**
 * A magnet of the account. In a live mode delta, only `id` and the
 * fields that changed are present.
 */
@Serializable
public data class Magnet(
    val id: Long,
    val filename: String? = null,
    /** Total size in bytes. */
    val size: Long? = null,
    /** The info hash. */
    val hash: String? = null,
    /** The status in plain English. */
    val status: String? = null,
    /**
     * `0` queued, `1` downloading, `2` compressing or moving, `3` uploading,
     * `4` ready. `5` to `15` are errors: `5` upload failed, `6` unpacking error,
     * `7` no peer, `8` too big, `9` internal error, `10` took too long,
     * `11` files removed from the host, `12` and `13` processing failed,
     * `14` tracker error, `15` no peer.
     */
    val statusCode: Int? = null,
    /** Bytes downloaded so far (while processing). */
    val downloaded: Long? = null,
    /** Bytes uploaded so far (while processing). */
    val uploaded: Long? = null,
    /** Seeders (while downloading). */
    val seeders: Int? = null,
    /** Bytes per second (while downloading). */
    val downloadSpeed: Long? = null,
    /** Bytes per second (while uploading). */
    val uploadSpeed: Long? = null,
    /** Progress percentage while processing. */
    val processingPerc: Double? = null,
    /** Unix time the magnet was added. */
    val uploadDate: Long? = null,
    /** Unix time the magnet became ready, or 0. */
    val completionDate: Long? = null,
    /** `m` when added as a magnet, `b` when added as a torrent file. */
    val type: String? = null,
    val notified: Boolean? = null,
    val version: Int? = null,
    /** Number of files, once ready. */
    val nbLinks: Int? = null,
    /** Only when a single magnet is requested with `id`. */
    val files: List<FileNode>? = null,
    /** Live mode deltas only. The magnet was deleted; drop it from the local state. */
    val deleted: Boolean? = null,
)

@Serializable
public data class MagnetStatusData(
    @Serializable(with = MagnetsSerializer::class) val magnets: List<Magnet>,
    /** Live mode. The counter to send on the next call. */
    val counter: Long? = null,
    /** Live mode. `true` when the response holds the full list. */
    val fullsync: Boolean? = null,
)

@Serializable
public data class MagnetFiles(
    /** The magnet ID (sent back as a string). */
    @Serializable(with = FlexStringSerializer::class) val id: FlexString,
    val files: List<FileNode>? = null,
    val error: ApiError? = null,
)

@Serializable
public data class MagnetFilesData(
    val magnets: List<MagnetFiles>,
)

@Serializable
public data class RestartedMagnet(
    /** The magnet ID. */
    val magnet: String? = null,
    val message: String? = null,
    val error: ApiError? = null,
)

@Serializable
public data class RestartMagnetsData(
    /** With `id`. The restarted magnet. */
    val id: Long? = null,
    /** With `id`. A confirmation message. */
    val message: String? = null,
    /** With `ids[]`. One entry per magnet. */
    val magnets: List<RestartedMagnet>? = null,
)

@Serializable
public data class SavedLink(
    val link: String,
    val filename: String? = null,
    val size: Long? = null,
    /** Unix time the link was saved or unlocked. */
    val date: Long? = null,
    val host: String? = null,
    /** History only. The download link that was generated. */
    @SerialName("link_dl") val linkDl: String? = null,
    /** History only. The error the unlock gave, if any. */
    val error: String? = null,
)

@Serializable
public data class SavedLinksData(
    val links: List<SavedLink>,
)

@Serializable
public data class VoucherBalanceData(
    /** The reseller balance. */
    val balance: Double,
)

@Serializable
public data class VouchersData(
    val codes: List<String>,
    /** `true` when fewer vouchers than requested were available. */
    val partialList: Boolean? = null,
)

@Serializable
public data class GeneratedVouchersData(
    val codes: List<String>,
    val pricePerVoucher: Double? = null,
    /** Total price of the request. */
    val total: Double? = null,
    /** The balance left after the generation. */
    val balance: Double? = null,
)

/** Error codes the API is known to return. The API may return others. */
public object ErrorCodes {
    public const val GENERIC: String = "GENERIC"
    public const val NOT_FOUND: String = "404"
    public const val MAINTENANCE: String = "MAINTENANCE"
    public const val DISCONTINUED: String = "DISCONTINUED"
    public const val AUTH_MISSING_APIKEY: String = "AUTH_MISSING_APIKEY"
    public const val AUTH_BAD_APIKEY: String = "AUTH_BAD_APIKEY"
    public const val AUTH_BLOCKED: String = "AUTH_BLOCKED"
    public const val AUTH_USER_BANNED: String = "AUTH_USER_BANNED"
    public const val ALREADY_SENT: String = "ALREADY_SENT"
    public const val NO_SERVER: String = "NO_SERVER"
    public const val LINK_IS_MISSING: String = "LINK_IS_MISSING"
    public const val BAD_LINK: String = "BAD_LINK"
    public const val LINK_HOST_NOT_SUPPORTED: String = "LINK_HOST_NOT_SUPPORTED"
    public const val LINK_DOWN: String = "LINK_DOWN"
    public const val LINK_PASS_PROTECTED: String = "LINK_PASS_PROTECTED"
    public const val LINK_HOST_UNAVAILABLE: String = "LINK_HOST_UNAVAILABLE"
    public const val LINK_TOO_MANY_DOWNLOADS: String = "LINK_TOO_MANY_DOWNLOADS"
    public const val LINK_HOST_FULL: String = "LINK_HOST_FULL"
    public const val LINK_HOST_LIMIT_REACHED: String = "LINK_HOST_LIMIT_REACHED"
    public const val LINK_ERROR: String = "LINK_ERROR"
    public const val LINK_TEMPORARY_UNAVAILABLE: String = "LINK_TEMPORARY_UNAVAILABLE"
    public const val LINK_NOT_SUPPORTED: String = "LINK_NOT_SUPPORTED"
    public const val REDIRECTOR_NOT_SUPPORTED: String = "REDIRECTOR_NOT_SUPPORTED"
    public const val REDIRECTOR_ERROR: String = "REDIRECTOR_ERROR"
    public const val REDIRECTOR_RATE_LIMITED: String = "REDIRECTOR_RATE_LIMITED"
    public const val STREAM_INVALID_GEN_ID: String = "STREAM_INVALID_GEN_ID"
    public const val STREAM_INVALID_STREAM_ID: String = "STREAM_INVALID_STREAM_ID"
    public const val DELAYED_INVALID_ID: String = "DELAYED_INVALID_ID"
    public const val FREE_TRIAL_LIMIT_REACHED: String = "FREE_TRIAL_LIMIT_REACHED"
    public const val MUST_BE_PREMIUM: String = "MUST_BE_PREMIUM"
    public const val MAGNET_INVALID_ID: String = "MAGNET_INVALID_ID"
    public const val MAGNET_INVALID_URI: String = "MAGNET_INVALID_URI"
    public const val MAGNET_INVALID_FILE: String = "MAGNET_INVALID_FILE"
    public const val MAGNET_FILE_UPLOAD_FAILED: String = "MAGNET_FILE_UPLOAD_FAILED"
    public const val MAGNET_NO_URI: String = "MAGNET_NO_URI"
    public const val MAGNET_PROCESSING: String = "MAGNET_PROCESSING"
    public const val MAGNET_TOO_MANY_ACTIVE: String = "MAGNET_TOO_MANY_ACTIVE"
    public const val MAGNET_TOO_MANY: String = "MAGNET_TOO_MANY"
    public const val MAGNET_MUST_BE_PREMIUM: String = "MAGNET_MUST_BE_PREMIUM"
    public const val MAGNET_NO_SERVER: String = "MAGNET_NO_SERVER"
    public const val MAGNET_TOO_LARGE: String = "MAGNET_TOO_LARGE"
    public const val MAGNET_UPLOAD_FAILED: String = "MAGNET_UPLOAD_FAILED"
    public const val MAGNET_INTERNAL_ERROR: String = "MAGNET_INTERNAL_ERROR"
    public const val MAGNET_CANT_BOOTSTRAP: String = "MAGNET_CANT_BOOTSTRAP"
    public const val MAGNET_MAGNET_TOO_BIG: String = "MAGNET_MAGNET_TOO_BIG"
    public const val MAGNET_TOOK_TOO_LONG: String = "MAGNET_TOOK_TOO_LONG"
    public const val MAGNET_LINKS_REMOVED: String = "MAGNET_LINKS_REMOVED"
    public const val MAGNET_PROCESSING_FAILED: String = "MAGNET_PROCESSING_FAILED"
    public const val PIN_ALREADY_AUTHED: String = "PIN_ALREADY_AUTHED"
    public const val PIN_EXPIRED: String = "PIN_EXPIRED"
    public const val PIN_INVALID: String = "PIN_INVALID"
    public const val USER_LINK_MISSING: String = "USER_LINK_MISSING"
    public const val USER_LINK_INVALID: String = "USER_LINK_INVALID"
    public const val MISSING_NOTIF_ENDPOINT: String = "MISSING_NOTIF_ENDPOINT"
    public const val VOUCHER_DURATION_INVALID: String = "VOUCHER_DURATION_INVALID"
    public const val VOUCHER_NB_INVALID: String = "VOUCHER_NB_INVALID"
    public const val NO_MORE_VOUCHER: String = "NO_MORE_VOUCHER"
    public const val INSUFFICIENT_BALANCE: String = "INSUFFICIENT_BALANCE"
    public const val DOWNLOAD_FAILED: String = "DOWNLOAD_FAILED"
    public const val ACCOUNT_INVALID: String = "ACCOUNT_INVALID"
    public const val NO_JSON_PARAM: String = "NO_JSON_PARAM"
    public const val JSON_INVALID: String = "JSON_INVALID"
    public const val FREEDAYS_INVALID_COUNTRY: String = "FREEDAYS_INVALID_COUNTRY"
    public const val FREEDAYS_INVALID_PHONE: String = "FREEDAYS_INVALID_PHONE"
    public const val FREEDAYS_INVALID_PROVIDER: String = "FREEDAYS_INVALID_PROVIDER"
    public const val FREEDAYS_USED_PHONE: String = "FREEDAYS_USED_PHONE"
    public const val FREEDAYS_ALREADY_REQUESTED: String = "FREEDAYS_ALREADY_REQUESTED"
    public const val FREEDAYS_INVALID_STATUS: String = "FREEDAYS_INVALID_STATUS"
    public const val FREEDAYS_TOO_MUCH_RETRIES: String = "FREEDAYS_TOO_MUCH_RETRIES"
    public const val FREEDAYS_INVALID_CODE: String = "FREEDAYS_INVALID_CODE"
    public const val FREEDAYS_DENIED: String = "FREEDAYS_DENIED"
    public const val FREEDAYS_ERROR_SENDING: String = "FREEDAYS_ERROR_SENDING"

    /** Every code of the spec. */
    public val ALL: List<String> = listOf(
        GENERIC,
        NOT_FOUND,
        MAINTENANCE,
        DISCONTINUED,
        AUTH_MISSING_APIKEY,
        AUTH_BAD_APIKEY,
        AUTH_BLOCKED,
        AUTH_USER_BANNED,
        ALREADY_SENT,
        NO_SERVER,
        LINK_IS_MISSING,
        BAD_LINK,
        LINK_HOST_NOT_SUPPORTED,
        LINK_DOWN,
        LINK_PASS_PROTECTED,
        LINK_HOST_UNAVAILABLE,
        LINK_TOO_MANY_DOWNLOADS,
        LINK_HOST_FULL,
        LINK_HOST_LIMIT_REACHED,
        LINK_ERROR,
        LINK_TEMPORARY_UNAVAILABLE,
        LINK_NOT_SUPPORTED,
        REDIRECTOR_NOT_SUPPORTED,
        REDIRECTOR_ERROR,
        REDIRECTOR_RATE_LIMITED,
        STREAM_INVALID_GEN_ID,
        STREAM_INVALID_STREAM_ID,
        DELAYED_INVALID_ID,
        FREE_TRIAL_LIMIT_REACHED,
        MUST_BE_PREMIUM,
        MAGNET_INVALID_ID,
        MAGNET_INVALID_URI,
        MAGNET_INVALID_FILE,
        MAGNET_FILE_UPLOAD_FAILED,
        MAGNET_NO_URI,
        MAGNET_PROCESSING,
        MAGNET_TOO_MANY_ACTIVE,
        MAGNET_TOO_MANY,
        MAGNET_MUST_BE_PREMIUM,
        MAGNET_NO_SERVER,
        MAGNET_TOO_LARGE,
        MAGNET_UPLOAD_FAILED,
        MAGNET_INTERNAL_ERROR,
        MAGNET_CANT_BOOTSTRAP,
        MAGNET_MAGNET_TOO_BIG,
        MAGNET_TOOK_TOO_LONG,
        MAGNET_LINKS_REMOVED,
        MAGNET_PROCESSING_FAILED,
        PIN_ALREADY_AUTHED,
        PIN_EXPIRED,
        PIN_INVALID,
        USER_LINK_MISSING,
        USER_LINK_INVALID,
        MISSING_NOTIF_ENDPOINT,
        VOUCHER_DURATION_INVALID,
        VOUCHER_NB_INVALID,
        NO_MORE_VOUCHER,
        INSUFFICIENT_BALANCE,
        DOWNLOAD_FAILED,
        ACCOUNT_INVALID,
        NO_JSON_PARAM,
        JSON_INVALID,
        FREEDAYS_INVALID_COUNTRY,
        FREEDAYS_INVALID_PHONE,
        FREEDAYS_INVALID_PROVIDER,
        FREEDAYS_USED_PHONE,
        FREEDAYS_ALREADY_REQUESTED,
        FREEDAYS_INVALID_STATUS,
        FREEDAYS_TOO_MUCH_RETRIES,
        FREEDAYS_INVALID_CODE,
        FREEDAYS_DENIED,
        FREEDAYS_ERROR_SENDING,
    )
}
