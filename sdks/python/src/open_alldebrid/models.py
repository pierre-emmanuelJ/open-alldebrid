# Generated from openapi.yaml by datamodel-code-generator. DO NOT EDIT.

from __future__ import annotations

from typing import Any, Literal, NotRequired, TypeAlias, TypedDict


ResponseStatus: TypeAlias = Literal["success", "error"]
"""
`success` with `data`, or `error` with `error`.
"""


class Error(TypedDict):
    code: str
    """
    A machine-readable error code; see `ErrorCode`.
    """
    message: str
    """
    A human-readable message, in English.
    """
    token: NotRequired[str]
    """
    With `AUTH_BLOCKED`, the token to poll on `/v4/user/verif`.
    """


ErrorCode: TypeAlias = Literal[
    "GENERIC",
    "404",
    "MAINTENANCE",
    "DISCONTINUED",
    "AUTH_MISSING_APIKEY",
    "AUTH_BAD_APIKEY",
    "AUTH_BLOCKED",
    "AUTH_USER_BANNED",
    "ALREADY_SENT",
    "NO_SERVER",
    "LINK_IS_MISSING",
    "BAD_LINK",
    "LINK_HOST_NOT_SUPPORTED",
    "LINK_DOWN",
    "LINK_PASS_PROTECTED",
    "LINK_HOST_UNAVAILABLE",
    "LINK_TOO_MANY_DOWNLOADS",
    "LINK_HOST_FULL",
    "LINK_HOST_LIMIT_REACHED",
    "LINK_ERROR",
    "LINK_TEMPORARY_UNAVAILABLE",
    "LINK_NOT_SUPPORTED",
    "REDIRECTOR_NOT_SUPPORTED",
    "REDIRECTOR_ERROR",
    "REDIRECTOR_RATE_LIMITED",
    "STREAM_INVALID_GEN_ID",
    "STREAM_INVALID_STREAM_ID",
    "DELAYED_INVALID_ID",
    "FREE_TRIAL_LIMIT_REACHED",
    "MUST_BE_PREMIUM",
    "MAGNET_INVALID_ID",
    "MAGNET_INVALID_URI",
    "MAGNET_INVALID_FILE",
    "MAGNET_FILE_UPLOAD_FAILED",
    "MAGNET_NO_URI",
    "MAGNET_PROCESSING",
    "MAGNET_TOO_MANY_ACTIVE",
    "MAGNET_TOO_MANY",
    "MAGNET_MUST_BE_PREMIUM",
    "MAGNET_NO_SERVER",
    "MAGNET_TOO_LARGE",
    "MAGNET_UPLOAD_FAILED",
    "MAGNET_INTERNAL_ERROR",
    "MAGNET_CANT_BOOTSTRAP",
    "MAGNET_MAGNET_TOO_BIG",
    "MAGNET_TOOK_TOO_LONG",
    "MAGNET_LINKS_REMOVED",
    "MAGNET_PROCESSING_FAILED",
    "PIN_ALREADY_AUTHED",
    "PIN_EXPIRED",
    "PIN_INVALID",
    "USER_LINK_MISSING",
    "USER_LINK_INVALID",
    "MISSING_NOTIF_ENDPOINT",
    "VOUCHER_DURATION_INVALID",
    "VOUCHER_NB_INVALID",
    "NO_MORE_VOUCHER",
    "INSUFFICIENT_BALANCE",
    "DOWNLOAD_FAILED",
    "ACCOUNT_INVALID",
    "NO_JSON_PARAM",
    "JSON_INVALID",
    "FREEDAYS_INVALID_COUNTRY",
    "FREEDAYS_INVALID_PHONE",
    "FREEDAYS_INVALID_PROVIDER",
    "FREEDAYS_USED_PHONE",
    "FREEDAYS_ALREADY_REQUESTED",
    "FREEDAYS_INVALID_STATUS",
    "FREEDAYS_TOO_MUCH_RETRIES",
    "FREEDAYS_INVALID_CODE",
    "FREEDAYS_DENIED",
    "FREEDAYS_ERROR_SENDING",
]
"""
Error codes the API is known to return. New codes may appear, so
`Error.code` is a plain string and this list is informational.

"""


StringOrInteger: TypeAlias = str | int
"""
A value the API sends either as a string or as an integer.
"""


StringOrStrings: TypeAlias = str | list[str]
"""
A value the API sends either as one string or as a list of strings.
"""


class PingData(TypedDict):
    ping: Literal["pong"]
    version: NotRequired[str]
    """
    API major version.
    """
    subversion: NotRequired[bool]
    """
    Only on `/v4.1/ping`.
    """


class PingResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[PingData]
    error: NotRequired[Error]


class MessageData(TypedDict):
    message: str
    """
    A confirmation message.
    """


class MessageResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[MessageData]
    error: NotRequired[Error]


class PinGetData(TypedDict):
    pin: str
    """
    The code to show to the user.
    """
    check: str
    """
    The token to send to `/v4/pin/check`.
    """
    expires_in: int
    """
    Seconds before the PIN expires.
    """
    user_url: str
    """
    The page where the user enters the PIN, with the PIN filled in.
    """
    base_url: str
    """
    The page where the user enters the PIN.
    """


class PinGetResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[PinGetData]
    error: NotRequired[Error]


class PinCheckRequest(TypedDict):
    pin: str
    """
    The `pin` from `/v4.1/pin/get`.
    """
    check: str
    """
    The `check` token from `/v4.1/pin/get`.
    """


class PinCheckData(TypedDict):
    activated: bool
    """
    `true` once the user has entered the PIN.
    """
    expires_in: int
    """
    Seconds before the PIN expires.
    """
    apikey: NotRequired[str]
    """
    The new API key, once `activated` is `true`.
    """


class PinCheckResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[PinCheckData]
    error: NotRequired[Error]


class Service(TypedDict):
    """
    A supported file host, streaming site or redirector.
    """

    name: str
    type: str
    """
    `premium` (needs a premium account) or `free`.
    """
    domains: list[str]
    """
    Every domain of the service.
    """
    regexps: NotRequired[list[str]]
    """
    Regular expressions matching the supported links.
    """
    regexp: NotRequired[StringOrStrings]
    """
    The `regexps` joined into one expression (sometimes a list).
    """
    status: NotRequired[bool]
    """
    Whether the service worked in the last 5 minutes. Only some hosts are
    checked, about every 10 minutes.

    """
    display: NotRequired[str]
    """
    A display name, when it differs from `name`.
    """
    hardRedirect: NotRequired[list[str]]
    """
    Domains that redirect to this host.
    """
    heavy: NotRequired[bool]
    """
    Redirectors only. Extraction is slow or rate limited.
    """
    quota: NotRequired[int]
    """
    `/v4.1/user/hosts` only. The remaining quota (see `quotaType`).
    """
    quotaMax: NotRequired[int]
    """
    `/v4.1/user/hosts` only. The daily quota.
    """
    quotaType: NotRequired[str]
    """
    `/v4.1/user/hosts` only. `traffic` (quota in MB) or `nb_download`
    (number of downloads).

    """
    limitSimuDl: NotRequired[int]
    """
    `/v4.1/user/hosts` only. Free simultaneous download slots.
    """


ServiceMap: TypeAlias = dict[str, Service]
"""
Services keyed by name.
"""


class HostsData(TypedDict):
    hosts: ServiceMap
    streams: NotRequired[ServiceMap]
    redirectors: NotRequired[ServiceMap]


class HostsResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[HostsData]
    error: NotRequired[Error]


class HostDomainsData(TypedDict):
    hosts: list[str]
    streams: NotRequired[list[str]]
    redirectors: NotRequired[list[str]]


class HostDomainsResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[HostDomainsData]
    error: NotRequired[Error]


class HostPriorityData(TypedDict):
    hosts: dict[str, int]
    """
    Host names mapped to their rank, 1 being the most open.
    """


class HostPriorityResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[HostPriorityData]
    error: NotRequired[Error]


QuotaMap: TypeAlias = dict[str, int] | list[Any]
"""
Remaining quota per limited host, in MB. An empty map may be sent as
an empty JSON array.

"""


class User(TypedDict):
    username: str
    email: str
    isPremium: bool
    isSubscribed: bool
    """
    Whether the user has an active recurring subscription.
    """
    isTrial: bool
    """
    Whether the account is in its free trial.
    """
    premiumUntil: int
    """
    Unix time when premium ends, or 0 when not premium.
    """
    lang: str
    """
    Interface language, e.g. `en` or `fr`.
    """
    preferedDomain: str
    """
    Preferred alldebrid.* top-level domain, e.g. `com` or `fr`.
    """
    fidelityPoints: int
    limitedHostersQuotas: QuotaMap
    remainingTrialQuota: NotRequired[int]
    """
    During the trial, the remaining traffic in MB.
    """
    notifications: list[str]
    """
    Codes of the current notifications; clear them with `/v4/user/notification/clear`.
    """
    history: NotRequired[bool]
    """
    Whether recent links history logging is enabled.
    """


class UserData(TypedDict):
    user: User


class UserResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[UserData]
    error: NotRequired[Error]


class UserHostsData(TypedDict):
    hosts: ServiceMap


class UserHostsResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[UserHostsData]
    error: NotRequired[Error]


class VerificationRequest(TypedDict):
    token: str
    """
    The `token` returned with the `AUTH_BLOCKED` error.
    """


class VerificationData(TypedDict):
    verif: Literal["waiting", "allowed", "denied"]
    resendable: NotRequired[bool]
    """
    When `waiting`, whether `/v4/user/verif/resend` may be called.
    """
    apikey: NotRequired[str]
    """
    When `allowed`, the API key.
    """


class VerificationResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[VerificationData]
    error: NotRequired[Error]


class ResendVerificationData(TypedDict):
    sent: bool


class ResendVerificationResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[ResendVerificationData]
    error: NotRequired[Error]


class ClearNotificationRequest(TypedDict):
    code: str
    """
    A notification code from `/v4/user`.
    """


LinkInfosRequest = TypedDict(
    "LinkInfosRequest",
    {
        "link[]": list[str],
        "password": NotRequired[str],
    },
)


class LinkInfo(TypedDict):
    link: str
    """
    The requested link.
    """
    filename: NotRequired[str]
    size: NotRequired[int]
    """
    File size in bytes.
    """
    host: NotRequired[str]
    hostDomain: NotRequired[str]
    """
    The host's main domain.
    """
    error: NotRequired[Error]


class LinkInfosData(TypedDict):
    infos: list[LinkInfo]


class LinkInfosResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[LinkInfosData]
    error: NotRequired[Error]


class RedirectorRequest(TypedDict):
    link: str
    """
    The redirector or link protector URL.
    """


class RedirectorData(TypedDict):
    links: list[str]
    """
    Encrypted links, usable with `/v4/link/infos` and `/v4/link/unlock`.
    """


class RedirectorResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[RedirectorData]
    error: NotRequired[Error]


class UnlockRequest(TypedDict):
    link: str
    """
    The link to unlock.
    """
    password: NotRequired[str]
    """
    The link's password, if any.
    """


class Stream(TypedDict):
    """
    One quality of a streaming link.
    """

    id: str
    """
    The `stream` value for `/v4/link/streaming`.
    """
    ext: NotRequired[str]
    quality: NotRequired[StringOrInteger]
    """
    Height in pixels (e.g. `1080`), or a label such as `mp3`.
    """
    filesize: NotRequired[int]
    proto: NotRequired[str]
    name: NotRequired[str]
    tb: NotRequired[float]
    """
    Audio-only streams. Total bitrate.
    """
    abr: NotRequired[int]
    """
    Audio-only streams. Audio bitrate.
    """


class UnlockData(TypedDict):
    link: NotRequired[str]
    """
    The direct download link. Empty when `delayed` is set.
    """
    host: NotRequired[str]
    filename: NotRequired[str]
    filesize: NotRequired[int]
    """
    File size in bytes.
    """
    id: NotRequired[str]
    """
    The generation ID, needed by `/v4/link/streaming`.
    """
    hostDomain: NotRequired[str]
    paws: NotRequired[bool]
    """
    Unused.
    """
    streaming: NotRequired[list[Stream]]
    """
    Available qualities for streaming links; empty otherwise.
    """
    delayed: NotRequired[StringOrInteger]
    """
    Set when the link is still being generated; poll `/v4/link/delayed` with it.
    """


class UnlockResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[UnlockData]
    error: NotRequired[Error]


class StreamingRequest(TypedDict):
    id: str
    """
    The `id` returned by `/v4/link/unlock`.
    """
    stream: str
    """
    The `id` of the chosen entry of `streaming`.
    """


class StreamingData(TypedDict):
    link: NotRequired[str]
    """
    The download link; absent when `delayed` is set.
    """
    filename: NotRequired[str]
    filesize: NotRequired[int]
    delayed: NotRequired[StringOrInteger]
    """
    Set when the link is still being generated; poll `/v4/link/delayed` with it.
    """


class StreamingResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[StreamingData]
    error: NotRequired[Error]


class DelayedRequest(TypedDict):
    id: str
    """
    The `delayed` value from `/v4/link/unlock` or `/v4/link/streaming`.
    """


class DelayedData(TypedDict):
    status: Literal[1, 2, 3]
    """
    `1` still processing, `2` ready (`link` is set), `3` failed.
    """
    time_left: int
    """
    Estimated seconds left.
    """
    link: NotRequired[str]
    """
    The download link, once `status` is `2`.
    """


class DelayedResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[DelayedData]
    error: NotRequired[Error]


UploadMagnetsRequest = TypedDict(
    "UploadMagnetsRequest",
    {
        "magnets[]": list[str],
    },
)


class UploadedMagnet(TypedDict):
    """
    An uploaded magnet or torrent, or a per-item `error`.
    """

    magnet: NotRequired[str]
    """
    The magnet as sent (for torrent files, the magnet built from the file).
    """
    file: NotRequired[str]
    """
    Torrent file uploads with an `error` only. The file name sent.
    """
    hash: NotRequired[str]
    """
    The info hash.
    """
    name: NotRequired[str]
    """
    The torrent name, or `noname` when unknown.
    """
    filename_original: NotRequired[str]
    size: NotRequired[int]
    """
    Total size in bytes.
    """
    ready: NotRequired[bool]
    """
    Whether the files are already available.
    """
    id: NotRequired[int]
    """
    The magnet ID, for `/v4.1/magnet/status`.
    """
    error: NotRequired[Error]


class UploadMagnetsData(TypedDict):
    magnets: list[UploadedMagnet]


class UploadMagnetsResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[UploadMagnetsData]
    error: NotRequired[Error]


UploadTorrentFilesRequest = TypedDict(
    "UploadTorrentFilesRequest",
    {
        "files[]": list[str],
    },
)


class UploadTorrentFilesData(TypedDict):
    files: list[UploadedMagnet]


class UploadTorrentFilesResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[UploadTorrentFilesData]
    error: NotRequired[Error]


class MagnetStatusRequest(TypedDict):
    id: NotRequired[int]
    """
    Only this magnet; the response then holds a single object.
    """
    status: NotRequired[Literal["active", "ready", "expired", "error"]]
    """
    Documented filter. The live `/v4.1/` endpoint currently ignores it
    and returns every magnet.

    """
    session: NotRequired[int]
    """
    Live mode. A random session ID chosen by the client.
    """
    counter: NotRequired[int]
    """
    Live mode. `0` first, then the `counter` of the previous response.
    """


class FileNode(TypedDict):
    """
    A file (`n`, `s`, `l`) or a folder (`n`, `e`) of a magnet.

    """

    n: str
    """
    File or folder name.
    """
    s: NotRequired[int]
    """
    Files only. Size in bytes.
    """
    l: NotRequired[str]
    """
    Files only. Link to unlock with `/v4/link/unlock`.
    """
    e: NotRequired[list[FileNode]]
    """
    Folders only. The entries of the folder.
    """


class Magnet(TypedDict):
    """
    A magnet of the account. In a live mode delta, only `id` and the
    fields that changed are present.

    """

    id: int
    filename: NotRequired[str]
    size: NotRequired[int]
    """
    Total size in bytes.
    """
    hash: NotRequired[str]
    """
    The info hash.
    """
    status: NotRequired[str]
    """
    The status in plain English.
    """
    statusCode: NotRequired[int]
    """
    `0` queued, `1` downloading, `2` compressing or moving, `3` uploading,
    `4` ready. `5` to `15` are errors: `5` upload failed, `6` unpacking error,
    `7` no peer, `8` too big, `9` internal error, `10` took too long,
    `11` files removed from the host, `12` and `13` processing failed,
    `14` tracker error, `15` no peer.

    """
    downloaded: NotRequired[int]
    """
    Bytes downloaded so far (while processing).
    """
    uploaded: NotRequired[int]
    """
    Bytes uploaded so far (while processing).
    """
    seeders: NotRequired[int]
    """
    Seeders (while downloading).
    """
    downloadSpeed: NotRequired[int]
    """
    Bytes per second (while downloading).
    """
    uploadSpeed: NotRequired[int]
    """
    Bytes per second (while uploading).
    """
    processingPerc: NotRequired[float]
    """
    Progress percentage while processing.
    """
    uploadDate: NotRequired[int]
    """
    Unix time the magnet was added.
    """
    completionDate: NotRequired[int]
    """
    Unix time the magnet became ready, or 0.
    """
    type: NotRequired[str]
    """
    `m` when added as a magnet, `b` when added as a torrent file.
    """
    notified: NotRequired[bool]
    version: NotRequired[int]
    nbLinks: NotRequired[int]
    """
    Number of files, once ready.
    """
    files: NotRequired[list[FileNode]]
    """
    Only when a single magnet is requested with `id`.
    """
    deleted: NotRequired[bool]
    """
    Live mode deltas only. The magnet was deleted; drop it from the local state.
    """


MagnetList: TypeAlias = list[Magnet] | Magnet
"""
A list of magnets, or a single magnet object when one `id` was requested
(and sometimes in live mode deltas).

"""


class MagnetStatusData(TypedDict):
    magnets: MagnetList
    counter: NotRequired[int]
    """
    Live mode. The counter to send on the next call.
    """
    fullsync: NotRequired[bool]
    """
    Live mode. `true` when the response holds the full list.
    """


class MagnetStatusResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[MagnetStatusData]
    error: NotRequired[Error]


MagnetFilesRequest = TypedDict(
    "MagnetFilesRequest",
    {
        "id[]": list[int],
    },
)


class MagnetFiles(TypedDict):
    id: StringOrInteger
    """
    The magnet ID (sent back as a string).
    """
    files: NotRequired[list[FileNode]]
    error: NotRequired[Error]


class MagnetFilesData(TypedDict):
    magnets: list[MagnetFiles]


class MagnetFilesResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[MagnetFilesData]
    error: NotRequired[Error]


class MagnetIdRequest(TypedDict):
    id: int
    """
    The magnet ID.
    """


RestartMagnetsRequest = TypedDict(
    "RestartMagnetsRequest",
    {
        "id": NotRequired[int],
        "ids[]": NotRequired[list[int]],
    },
)


class RestartedMagnet(TypedDict):
    magnet: NotRequired[str]
    """
    The magnet ID.
    """
    message: NotRequired[str]
    error: NotRequired[Error]


class RestartMagnetsData(TypedDict):
    id: NotRequired[int]
    """
    With `id`. The restarted magnet.
    """
    message: NotRequired[str]
    """
    With `id`. A confirmation message.
    """
    magnets: NotRequired[list[RestartedMagnet]]
    """
    With `ids[]`. One entry per magnet.
    """


class RestartMagnetsResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[RestartMagnetsData]
    error: NotRequired[Error]


class SavedLink(TypedDict):
    link: str
    filename: NotRequired[str]
    size: NotRequired[int]
    date: NotRequired[int]
    """
    Unix time the link was saved or unlocked.
    """
    host: NotRequired[str]
    link_dl: NotRequired[str]
    """
    History only. The download link that was generated.
    """
    error: NotRequired[str]
    """
    History only. The error the unlock gave, if any.
    """


class SavedLinksData(TypedDict):
    links: list[SavedLink]


class SavedLinksResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[SavedLinksData]
    error: NotRequired[Error]


class HistoryResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[SavedLinksData]
    error: NotRequired[Error]


LinksRequest = TypedDict(
    "LinksRequest",
    {
        "links[]": list[str],
    },
)


class VoucherBalanceData(TypedDict):
    balance: float
    """
    The reseller balance.
    """


class VoucherBalanceResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[VoucherBalanceData]
    error: NotRequired[Error]


class VoucherRequest(TypedDict):
    duration: Literal[15, 30, 90, 180, 365]
    """
    Voucher duration in days.
    """
    nb: int
    """
    Number of vouchers.
    """


class VouchersData(TypedDict):
    codes: list[str]
    partialList: NotRequired[bool]
    """
    `true` when fewer vouchers than requested were available.
    """


class VouchersResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[VouchersData]
    error: NotRequired[Error]


class GeneratedVouchersData(TypedDict):
    codes: list[str]
    pricePerVoucher: NotRequired[float]
    total: NotRequired[float]
    """
    Total price of the request.
    """
    balance: NotRequired[float]
    """
    The balance left after the generation.
    """


class GeneratedVouchersResponse(TypedDict):
    status: ResponseStatus
    demo: NotRequired[bool]
    deprecated: NotRequired[bool]
    data: NotRequired[GeneratedVouchersData]
    error: NotRequired[Error]
