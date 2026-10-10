// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

#![allow(clippy::doc_markdown, clippy::doc_lazy_continuation)]

use std::collections::HashMap;

use serde::{Deserialize, Serialize};

use crate::flex::{self, FlexString};

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct ApiError {
    /// A machine-readable error code; see `ErrorCode`.
    pub code: String,
    /// A human-readable message, in English.
    pub message: String,
    /// With `AUTH_BLOCKED`, the token to poll on `/v4/user/verif`.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub token: Option<String>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct PingData {
    pub ping: String,
    /// API major version.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub version: Option<String>,
    /// Only on `/v4.1/ping`.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub subversion: Option<bool>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct MessageData {
    /// A confirmation message.
    pub message: String,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct PinGetData {
    /// The code to show to the user.
    pub pin: String,
    /// The token to send to `/v4/pin/check`.
    pub check: String,
    /// Seconds before the PIN expires.
    pub expires_in: i64,
    /// The page where the user enters the PIN, with the PIN filled in.
    pub user_url: String,
    /// The page where the user enters the PIN.
    pub base_url: String,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct PinCheckData {
    /// `true` once the user has entered the PIN.
    pub activated: bool,
    /// Seconds before the PIN expires.
    pub expires_in: i64,
    /// The new API key, once `activated` is `true`.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub apikey: Option<String>,
}

/// A supported file host, streaming site or redirector.
#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct Service {
    pub name: String,
    /// `premium` (needs a premium account) or `free`.
    pub r#type: String,
    /// Every domain of the service.
    pub domains: Vec<String>,
    /// Regular expressions matching the supported links.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub regexps: Option<Vec<String>>,
    /// The `regexps` joined into one expression (sometimes a list).
    #[serde(default, skip_serializing_if = "Option::is_none", deserialize_with = "flex::opt_strings")]
    pub regexp: Option<Vec<String>>,
    /// Whether the service worked in the last 5 minutes. Only some hosts are
    /// checked, about every 10 minutes.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub status: Option<bool>,
    /// A display name, when it differs from `name`.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub display: Option<String>,
    /// Domains that redirect to this host.
    #[serde(rename = "hardRedirect", default, skip_serializing_if = "Option::is_none")]
    pub hard_redirect: Option<Vec<String>>,
    /// Redirectors only. Extraction is slow or rate limited.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub heavy: Option<bool>,
    /// `/v4.1/user/hosts` only. The remaining quota (see `quotaType`).
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub quota: Option<i64>,
    /// `/v4.1/user/hosts` only. The daily quota.
    #[serde(rename = "quotaMax", default, skip_serializing_if = "Option::is_none")]
    pub quota_max: Option<i64>,
    /// `/v4.1/user/hosts` only. `traffic` (quota in MB) or `nb_download`
    /// (number of downloads).
    #[serde(rename = "quotaType", default, skip_serializing_if = "Option::is_none")]
    pub quota_type: Option<String>,
    /// `/v4.1/user/hosts` only. Free simultaneous download slots.
    #[serde(rename = "limitSimuDl", default, skip_serializing_if = "Option::is_none")]
    pub limit_simu_dl: Option<i64>,
}

/// Services keyed by name.
pub type ServiceMap = HashMap<String, Service>;

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct HostsData {
    pub hosts: ServiceMap,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub streams: Option<ServiceMap>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub redirectors: Option<ServiceMap>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct HostDomainsData {
    pub hosts: Vec<String>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub streams: Option<Vec<String>>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub redirectors: Option<Vec<String>>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct HostPriorityData {
    /// Host names mapped to their rank, 1 being the most open.
    pub hosts: HashMap<String, i64>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct User {
    pub username: String,
    pub email: String,
    #[serde(rename = "isPremium")]
    pub is_premium: bool,
    /// Whether the user has an active recurring subscription.
    #[serde(rename = "isSubscribed")]
    pub is_subscribed: bool,
    /// Whether the account is in its free trial.
    #[serde(rename = "isTrial")]
    pub is_trial: bool,
    /// Unix time when premium ends, or 0 when not premium.
    #[serde(rename = "premiumUntil")]
    pub premium_until: i64,
    /// Interface language, e.g. `en` or `fr`.
    pub lang: String,
    /// Preferred alldebrid.* top-level domain, e.g. `com` or `fr`.
    #[serde(rename = "preferedDomain")]
    pub prefered_domain: String,
    #[serde(rename = "fidelityPoints")]
    pub fidelity_points: i64,
    #[serde(rename = "limitedHostersQuotas", deserialize_with = "flex::quotas")]
    pub limited_hosters_quotas: HashMap<String, i64>,
    /// During the trial, the remaining traffic in MB.
    #[serde(rename = "remainingTrialQuota", default, skip_serializing_if = "Option::is_none")]
    pub remaining_trial_quota: Option<i64>,
    /// Codes of the current notifications; clear them with `/v4/user/notification/clear`.
    pub notifications: Vec<String>,
    /// Whether recent links history logging is enabled.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub history: Option<bool>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct UserData {
    pub user: User,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct UserHostsData {
    pub hosts: ServiceMap,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct VerificationData {
    pub verif: String,
    /// When `waiting`, whether `/v4/user/verif/resend` may be called.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub resendable: Option<bool>,
    /// When `allowed`, the API key.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub apikey: Option<String>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct ResendVerificationData {
    pub sent: bool,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct LinkInfo {
    /// The requested link.
    pub link: String,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub filename: Option<String>,
    /// File size in bytes.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub size: Option<i64>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub host: Option<String>,
    /// The host's main domain.
    #[serde(rename = "hostDomain", default, skip_serializing_if = "Option::is_none")]
    pub host_domain: Option<String>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub error: Option<ApiError>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct LinkInfosData {
    pub infos: Vec<LinkInfo>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct RedirectorData {
    /// Encrypted links, usable with `/v4/link/infos` and `/v4/link/unlock`.
    pub links: Vec<String>,
}

/// One quality of a streaming link.
#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct Stream {
    /// The `stream` value for `/v4/link/streaming`.
    pub id: String,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub ext: Option<String>,
    /// Height in pixels (e.g. `1080`), or a label such as `mp3`.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub quality: Option<FlexString>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub filesize: Option<i64>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub proto: Option<String>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub name: Option<String>,
    /// Audio-only streams. Total bitrate.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub tb: Option<f64>,
    /// Audio-only streams. Audio bitrate.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub abr: Option<i64>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct UnlockData {
    /// The direct download link. Empty when `delayed` is set.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub link: Option<String>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub host: Option<String>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub filename: Option<String>,
    /// File size in bytes.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub filesize: Option<i64>,
    /// The generation ID, needed by `/v4/link/streaming`.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub id: Option<String>,
    #[serde(rename = "hostDomain", default, skip_serializing_if = "Option::is_none")]
    pub host_domain: Option<String>,
    /// Unused.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub paws: Option<bool>,
    /// Available qualities for streaming links; empty otherwise.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub streaming: Option<Vec<Stream>>,
    /// Set when the link is still being generated; poll `/v4/link/delayed` with it.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub delayed: Option<FlexString>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct StreamingData {
    /// The download link; absent when `delayed` is set.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub link: Option<String>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub filename: Option<String>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub filesize: Option<i64>,
    /// Set when the link is still being generated; poll `/v4/link/delayed` with it.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub delayed: Option<FlexString>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct DelayedData {
    /// `1` still processing, `2` ready (`link` is set), `3` failed.
    pub status: i64,
    /// Estimated seconds left.
    pub time_left: i64,
    /// The download link, once `status` is `2`.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub link: Option<String>,
}

/// An uploaded magnet or torrent, or a per-item `error`.
#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct UploadedMagnet {
    /// The magnet as sent (for torrent files, the magnet built from the file).
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub magnet: Option<String>,
    /// Torrent file uploads with an `error` only. The file name sent.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub file: Option<String>,
    /// The info hash.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub hash: Option<String>,
    /// The torrent name, or `noname` when unknown.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub name: Option<String>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub filename_original: Option<String>,
    /// Total size in bytes.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub size: Option<i64>,
    /// Whether the files are already available.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub ready: Option<bool>,
    /// The magnet ID, for `/v4.1/magnet/status`.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub id: Option<i64>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub error: Option<ApiError>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct UploadMagnetsData {
    pub magnets: Vec<UploadedMagnet>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct UploadTorrentFilesData {
    pub files: Vec<UploadedMagnet>,
}

/// A file (`n`, `s`, `l`) or a folder (`n`, `e`) of a magnet.
#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct FileNode {
    /// File or folder name.
    pub n: String,
    /// Files only. Size in bytes.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub s: Option<i64>,
    /// Files only. Link to unlock with `/v4/link/unlock`.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub l: Option<String>,
    /// Folders only. The entries of the folder.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub e: Option<Vec<FileNode>>,
}

/// A magnet of the account. In a live mode delta, only `id` and the
/// fields that changed are present.
#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct Magnet {
    pub id: i64,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub filename: Option<String>,
    /// Total size in bytes.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub size: Option<i64>,
    /// The info hash.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub hash: Option<String>,
    /// The status in plain English.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub status: Option<String>,
    /// `0` queued, `1` downloading, `2` compressing or moving, `3` uploading,
    /// `4` ready. `5` to `15` are errors: `5` upload failed, `6` unpacking error,
    /// `7` no peer, `8` too big, `9` internal error, `10` took too long,
    /// `11` files removed from the host, `12` and `13` processing failed,
    /// `14` tracker error, `15` no peer.
    #[serde(rename = "statusCode", default, skip_serializing_if = "Option::is_none")]
    pub status_code: Option<i64>,
    /// Bytes downloaded so far (while processing).
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub downloaded: Option<i64>,
    /// Bytes uploaded so far (while processing).
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub uploaded: Option<i64>,
    /// Seeders (while downloading).
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub seeders: Option<i64>,
    /// Bytes per second (while downloading).
    #[serde(rename = "downloadSpeed", default, skip_serializing_if = "Option::is_none")]
    pub download_speed: Option<i64>,
    /// Bytes per second (while uploading).
    #[serde(rename = "uploadSpeed", default, skip_serializing_if = "Option::is_none")]
    pub upload_speed: Option<i64>,
    /// Progress percentage while processing.
    #[serde(rename = "processingPerc", default, skip_serializing_if = "Option::is_none")]
    pub processing_perc: Option<f64>,
    /// Unix time the magnet was added.
    #[serde(rename = "uploadDate", default, skip_serializing_if = "Option::is_none")]
    pub upload_date: Option<i64>,
    /// Unix time the magnet became ready, or 0.
    #[serde(rename = "completionDate", default, skip_serializing_if = "Option::is_none")]
    pub completion_date: Option<i64>,
    /// `m` when added as a magnet, `b` when added as a torrent file.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub r#type: Option<String>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub notified: Option<bool>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub version: Option<i64>,
    /// Number of files, once ready.
    #[serde(rename = "nbLinks", default, skip_serializing_if = "Option::is_none")]
    pub nb_links: Option<i64>,
    /// Only when a single magnet is requested with `id`.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub files: Option<Vec<FileNode>>,
    /// Live mode deltas only. The magnet was deleted; drop it from the local state.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub deleted: Option<bool>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct MagnetStatusData {
    #[serde(deserialize_with = "flex::magnets")]
    pub magnets: Vec<Magnet>,
    /// Live mode. The counter to send on the next call.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub counter: Option<i64>,
    /// Live mode. `true` when the response holds the full list.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub fullsync: Option<bool>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct MagnetFiles {
    /// The magnet ID (sent back as a string).
    pub id: FlexString,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub files: Option<Vec<FileNode>>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub error: Option<ApiError>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct MagnetFilesData {
    pub magnets: Vec<MagnetFiles>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct RestartedMagnet {
    /// The magnet ID.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub magnet: Option<String>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub message: Option<String>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub error: Option<ApiError>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct RestartMagnetsData {
    /// With `id`. The restarted magnet.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub id: Option<i64>,
    /// With `id`. A confirmation message.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub message: Option<String>,
    /// With `ids[]`. One entry per magnet.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub magnets: Option<Vec<RestartedMagnet>>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct SavedLink {
    pub link: String,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub filename: Option<String>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub size: Option<i64>,
    /// Unix time the link was saved or unlocked.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub date: Option<i64>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub host: Option<String>,
    /// History only. The download link that was generated.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub link_dl: Option<String>,
    /// History only. The error the unlock gave, if any.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub error: Option<String>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct SavedLinksData {
    pub links: Vec<SavedLink>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct VoucherBalanceData {
    /// The reseller balance.
    pub balance: f64,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct VouchersData {
    pub codes: Vec<String>,
    /// `true` when fewer vouchers than requested were available.
    #[serde(rename = "partialList", default, skip_serializing_if = "Option::is_none")]
    pub partial_list: Option<bool>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct GeneratedVouchersData {
    pub codes: Vec<String>,
    #[serde(rename = "pricePerVoucher", default, skip_serializing_if = "Option::is_none")]
    pub price_per_voucher: Option<f64>,
    /// Total price of the request.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub total: Option<f64>,
    /// The balance left after the generation.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub balance: Option<f64>,
}

/// Error codes the API is known to return. The API may return others.
pub mod error_codes {
    pub const GENERIC: &str = "GENERIC";
    pub const NOT_FOUND: &str = "404";
    pub const MAINTENANCE: &str = "MAINTENANCE";
    pub const DISCONTINUED: &str = "DISCONTINUED";
    pub const AUTH_MISSING_APIKEY: &str = "AUTH_MISSING_APIKEY";
    pub const AUTH_BAD_APIKEY: &str = "AUTH_BAD_APIKEY";
    pub const AUTH_BLOCKED: &str = "AUTH_BLOCKED";
    pub const AUTH_USER_BANNED: &str = "AUTH_USER_BANNED";
    pub const ALREADY_SENT: &str = "ALREADY_SENT";
    pub const NO_SERVER: &str = "NO_SERVER";
    pub const LINK_IS_MISSING: &str = "LINK_IS_MISSING";
    pub const BAD_LINK: &str = "BAD_LINK";
    pub const LINK_HOST_NOT_SUPPORTED: &str = "LINK_HOST_NOT_SUPPORTED";
    pub const LINK_DOWN: &str = "LINK_DOWN";
    pub const LINK_PASS_PROTECTED: &str = "LINK_PASS_PROTECTED";
    pub const LINK_HOST_UNAVAILABLE: &str = "LINK_HOST_UNAVAILABLE";
    pub const LINK_TOO_MANY_DOWNLOADS: &str = "LINK_TOO_MANY_DOWNLOADS";
    pub const LINK_HOST_FULL: &str = "LINK_HOST_FULL";
    pub const LINK_HOST_LIMIT_REACHED: &str = "LINK_HOST_LIMIT_REACHED";
    pub const LINK_ERROR: &str = "LINK_ERROR";
    pub const LINK_TEMPORARY_UNAVAILABLE: &str = "LINK_TEMPORARY_UNAVAILABLE";
    pub const LINK_NOT_SUPPORTED: &str = "LINK_NOT_SUPPORTED";
    pub const REDIRECTOR_NOT_SUPPORTED: &str = "REDIRECTOR_NOT_SUPPORTED";
    pub const REDIRECTOR_ERROR: &str = "REDIRECTOR_ERROR";
    pub const REDIRECTOR_RATE_LIMITED: &str = "REDIRECTOR_RATE_LIMITED";
    pub const STREAM_INVALID_GEN_ID: &str = "STREAM_INVALID_GEN_ID";
    pub const STREAM_INVALID_STREAM_ID: &str = "STREAM_INVALID_STREAM_ID";
    pub const DELAYED_INVALID_ID: &str = "DELAYED_INVALID_ID";
    pub const FREE_TRIAL_LIMIT_REACHED: &str = "FREE_TRIAL_LIMIT_REACHED";
    pub const MUST_BE_PREMIUM: &str = "MUST_BE_PREMIUM";
    pub const MAGNET_INVALID_ID: &str = "MAGNET_INVALID_ID";
    pub const MAGNET_INVALID_URI: &str = "MAGNET_INVALID_URI";
    pub const MAGNET_INVALID_FILE: &str = "MAGNET_INVALID_FILE";
    pub const MAGNET_FILE_UPLOAD_FAILED: &str = "MAGNET_FILE_UPLOAD_FAILED";
    pub const MAGNET_NO_URI: &str = "MAGNET_NO_URI";
    pub const MAGNET_PROCESSING: &str = "MAGNET_PROCESSING";
    pub const MAGNET_TOO_MANY_ACTIVE: &str = "MAGNET_TOO_MANY_ACTIVE";
    pub const MAGNET_TOO_MANY: &str = "MAGNET_TOO_MANY";
    pub const MAGNET_MUST_BE_PREMIUM: &str = "MAGNET_MUST_BE_PREMIUM";
    pub const MAGNET_NO_SERVER: &str = "MAGNET_NO_SERVER";
    pub const MAGNET_TOO_LARGE: &str = "MAGNET_TOO_LARGE";
    pub const MAGNET_UPLOAD_FAILED: &str = "MAGNET_UPLOAD_FAILED";
    pub const MAGNET_INTERNAL_ERROR: &str = "MAGNET_INTERNAL_ERROR";
    pub const MAGNET_CANT_BOOTSTRAP: &str = "MAGNET_CANT_BOOTSTRAP";
    pub const MAGNET_MAGNET_TOO_BIG: &str = "MAGNET_MAGNET_TOO_BIG";
    pub const MAGNET_TOOK_TOO_LONG: &str = "MAGNET_TOOK_TOO_LONG";
    pub const MAGNET_LINKS_REMOVED: &str = "MAGNET_LINKS_REMOVED";
    pub const MAGNET_PROCESSING_FAILED: &str = "MAGNET_PROCESSING_FAILED";
    pub const PIN_ALREADY_AUTHED: &str = "PIN_ALREADY_AUTHED";
    pub const PIN_EXPIRED: &str = "PIN_EXPIRED";
    pub const PIN_INVALID: &str = "PIN_INVALID";
    pub const USER_LINK_MISSING: &str = "USER_LINK_MISSING";
    pub const USER_LINK_INVALID: &str = "USER_LINK_INVALID";
    pub const MISSING_NOTIF_ENDPOINT: &str = "MISSING_NOTIF_ENDPOINT";
    pub const VOUCHER_DURATION_INVALID: &str = "VOUCHER_DURATION_INVALID";
    pub const VOUCHER_NB_INVALID: &str = "VOUCHER_NB_INVALID";
    pub const NO_MORE_VOUCHER: &str = "NO_MORE_VOUCHER";
    pub const INSUFFICIENT_BALANCE: &str = "INSUFFICIENT_BALANCE";
    pub const DOWNLOAD_FAILED: &str = "DOWNLOAD_FAILED";
    pub const ACCOUNT_INVALID: &str = "ACCOUNT_INVALID";
    pub const NO_JSON_PARAM: &str = "NO_JSON_PARAM";
    pub const JSON_INVALID: &str = "JSON_INVALID";
    pub const FREEDAYS_INVALID_COUNTRY: &str = "FREEDAYS_INVALID_COUNTRY";
    pub const FREEDAYS_INVALID_PHONE: &str = "FREEDAYS_INVALID_PHONE";
    pub const FREEDAYS_INVALID_PROVIDER: &str = "FREEDAYS_INVALID_PROVIDER";
    pub const FREEDAYS_USED_PHONE: &str = "FREEDAYS_USED_PHONE";
    pub const FREEDAYS_ALREADY_REQUESTED: &str = "FREEDAYS_ALREADY_REQUESTED";
    pub const FREEDAYS_INVALID_STATUS: &str = "FREEDAYS_INVALID_STATUS";
    pub const FREEDAYS_TOO_MUCH_RETRIES: &str = "FREEDAYS_TOO_MUCH_RETRIES";
    pub const FREEDAYS_INVALID_CODE: &str = "FREEDAYS_INVALID_CODE";
    pub const FREEDAYS_DENIED: &str = "FREEDAYS_DENIED";
    pub const FREEDAYS_ERROR_SENDING: &str = "FREEDAYS_ERROR_SENDING";

    /// Every code of the spec.
    pub const ALL: &[&str] = &[
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
    ];
}
