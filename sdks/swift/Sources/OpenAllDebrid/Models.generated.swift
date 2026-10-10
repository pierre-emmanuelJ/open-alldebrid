// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.
// swift-format-ignore-file

import Foundation

public struct ApiError: Decodable, Sendable, Equatable {
    /// A machine-readable error code; see `ErrorCode`.
    public let code: String
    /// A human-readable message, in English.
    public let message: String
    /// With `AUTH_BLOCKED`, the token to poll on `/v4/user/verif`.
    public let token: String?

    public init(code: String, message: String, token: String? = nil) {
        self.code = code
        self.message = message
        self.token = token
    }

    private enum CodingKeys: String, CodingKey {
        case code
        case message
        case token
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        code = try c.decode(String.self, forKey: .code)
        message = try c.decode(String.self, forKey: .message)
        token = try c.decodeIfPresent(String.self, forKey: .token)
    }
}

public struct PingData: Decodable, Sendable, Equatable {
    public let ping: String
    /// API major version.
    public let version: String?
    /// Only on `/v4.1/ping`.
    public let subversion: Bool?

    public init(ping: String, version: String? = nil, subversion: Bool? = nil) {
        self.ping = ping
        self.version = version
        self.subversion = subversion
    }

    private enum CodingKeys: String, CodingKey {
        case ping
        case version
        case subversion
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        ping = try c.decode(String.self, forKey: .ping)
        version = try c.decodeIfPresent(String.self, forKey: .version)
        subversion = try c.decodeIfPresent(Bool.self, forKey: .subversion)
    }
}

public struct MessageData: Decodable, Sendable, Equatable {
    /// A confirmation message.
    public let message: String

    public init(message: String) {
        self.message = message
    }

    private enum CodingKeys: String, CodingKey {
        case message
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        message = try c.decode(String.self, forKey: .message)
    }
}

public struct PinGetData: Decodable, Sendable, Equatable {
    /// The code to show to the user.
    public let pin: String
    /// The token to send to `/v4/pin/check`.
    public let check: String
    /// Seconds before the PIN expires.
    public let expiresIn: Int
    /// The page where the user enters the PIN, with the PIN filled in.
    public let userUrl: String
    /// The page where the user enters the PIN.
    public let baseUrl: String

    public init(pin: String, check: String, expiresIn: Int, userUrl: String, baseUrl: String) {
        self.pin = pin
        self.check = check
        self.expiresIn = expiresIn
        self.userUrl = userUrl
        self.baseUrl = baseUrl
    }

    private enum CodingKeys: String, CodingKey {
        case pin
        case check
        case expiresIn = "expires_in"
        case userUrl = "user_url"
        case baseUrl = "base_url"
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        pin = try c.decode(String.self, forKey: .pin)
        check = try c.decode(String.self, forKey: .check)
        expiresIn = try c.decode(Int.self, forKey: .expiresIn)
        userUrl = try c.decode(String.self, forKey: .userUrl)
        baseUrl = try c.decode(String.self, forKey: .baseUrl)
    }
}

public struct PinCheckData: Decodable, Sendable, Equatable {
    /// `true` once the user has entered the PIN.
    public let activated: Bool
    /// Seconds before the PIN expires.
    public let expiresIn: Int
    /// The new API key, once `activated` is `true`.
    public let apikey: String?

    public init(activated: Bool, expiresIn: Int, apikey: String? = nil) {
        self.activated = activated
        self.expiresIn = expiresIn
        self.apikey = apikey
    }

    private enum CodingKeys: String, CodingKey {
        case activated
        case expiresIn = "expires_in"
        case apikey
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        activated = try c.decode(Bool.self, forKey: .activated)
        expiresIn = try c.decode(Int.self, forKey: .expiresIn)
        apikey = try c.decodeIfPresent(String.self, forKey: .apikey)
    }
}

/// A supported file host, streaming site or redirector.
public struct Service: Decodable, Sendable, Equatable {
    public let name: String
    /// `premium` (needs a premium account) or `free`.
    public let type: String
    /// Every domain of the service.
    public let domains: [String]
    /// Regular expressions matching the supported links.
    public let regexps: [String]?
    /// The `regexps` joined into one expression (sometimes a list).
    public let regexp: [String]?
    /// Whether the service worked in the last 5 minutes. Only some hosts are
    /// checked, about every 10 minutes.
    public let status: Bool?
    /// A display name, when it differs from `name`.
    public let display: String?
    /// Domains that redirect to this host.
    public let hardRedirect: [String]?
    /// Redirectors only. Extraction is slow or rate limited.
    public let heavy: Bool?
    /// `/v4.1/user/hosts` only. The remaining quota (see `quotaType`).
    public let quota: Int64?
    /// `/v4.1/user/hosts` only. The daily quota.
    public let quotaMax: Int64?
    /// `/v4.1/user/hosts` only. `traffic` (quota in MB) or `nb_download`
    /// (number of downloads).
    public let quotaType: String?
    /// `/v4.1/user/hosts` only. Free simultaneous download slots.
    public let limitSimuDl: Int?

    public init(name: String, type: String, domains: [String], regexps: [String]? = nil, regexp: [String]? = nil, status: Bool? = nil, display: String? = nil, hardRedirect: [String]? = nil, heavy: Bool? = nil, quota: Int64? = nil, quotaMax: Int64? = nil, quotaType: String? = nil, limitSimuDl: Int? = nil) {
        self.name = name
        self.type = type
        self.domains = domains
        self.regexps = regexps
        self.regexp = regexp
        self.status = status
        self.display = display
        self.hardRedirect = hardRedirect
        self.heavy = heavy
        self.quota = quota
        self.quotaMax = quotaMax
        self.quotaType = quotaType
        self.limitSimuDl = limitSimuDl
    }

    private enum CodingKeys: String, CodingKey {
        case name
        case type
        case domains
        case regexps
        case regexp
        case status
        case display
        case hardRedirect
        case heavy
        case quota
        case quotaMax
        case quotaType
        case limitSimuDl
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        name = try c.decode(String.self, forKey: .name)
        type = try c.decode(String.self, forKey: .type)
        domains = try c.decode([String].self, forKey: .domains)
        regexps = try c.decodeIfPresent([String].self, forKey: .regexps)
        regexp = try c.decodeStringsIfPresent(forKey: .regexp)
        status = try c.decodeIfPresent(Bool.self, forKey: .status)
        display = try c.decodeIfPresent(String.self, forKey: .display)
        hardRedirect = try c.decodeIfPresent([String].self, forKey: .hardRedirect)
        heavy = try c.decodeIfPresent(Bool.self, forKey: .heavy)
        quota = try c.decodeIfPresent(Int64.self, forKey: .quota)
        quotaMax = try c.decodeIfPresent(Int64.self, forKey: .quotaMax)
        quotaType = try c.decodeIfPresent(String.self, forKey: .quotaType)
        limitSimuDl = try c.decodeIfPresent(Int.self, forKey: .limitSimuDl)
    }
}

/// Services keyed by name.
public typealias ServiceMap = [String: Service]

public struct HostsData: Decodable, Sendable, Equatable {
    public let hosts: ServiceMap
    public let streams: ServiceMap?
    public let redirectors: ServiceMap?

    public init(hosts: ServiceMap, streams: ServiceMap? = nil, redirectors: ServiceMap? = nil) {
        self.hosts = hosts
        self.streams = streams
        self.redirectors = redirectors
    }

    private enum CodingKeys: String, CodingKey {
        case hosts
        case streams
        case redirectors
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        hosts = try c.decode(ServiceMap.self, forKey: .hosts)
        streams = try c.decodeIfPresent(ServiceMap.self, forKey: .streams)
        redirectors = try c.decodeIfPresent(ServiceMap.self, forKey: .redirectors)
    }
}

public struct HostDomainsData: Decodable, Sendable, Equatable {
    public let hosts: [String]
    public let streams: [String]?
    public let redirectors: [String]?

    public init(hosts: [String], streams: [String]? = nil, redirectors: [String]? = nil) {
        self.hosts = hosts
        self.streams = streams
        self.redirectors = redirectors
    }

    private enum CodingKeys: String, CodingKey {
        case hosts
        case streams
        case redirectors
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        hosts = try c.decode([String].self, forKey: .hosts)
        streams = try c.decodeIfPresent([String].self, forKey: .streams)
        redirectors = try c.decodeIfPresent([String].self, forKey: .redirectors)
    }
}

public struct HostPriorityData: Decodable, Sendable, Equatable {
    /// Host names mapped to their rank, 1 being the most open.
    public let hosts: [String: Int]

    public init(hosts: [String: Int]) {
        self.hosts = hosts
    }

    private enum CodingKeys: String, CodingKey {
        case hosts
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        hosts = try c.decode([String: Int].self, forKey: .hosts)
    }
}

public struct User: Decodable, Sendable, Equatable {
    public let username: String
    public let email: String
    public let isPremium: Bool
    /// Whether the user has an active recurring subscription.
    public let isSubscribed: Bool
    /// Whether the account is in its free trial.
    public let isTrial: Bool
    /// Unix time when premium ends, or 0 when not premium.
    public let premiumUntil: Int64
    /// Interface language, e.g. `en` or `fr`.
    public let lang: String
    /// Preferred alldebrid.* top-level domain, e.g. `com` or `fr`.
    public let preferedDomain: String
    public let fidelityPoints: Int
    public let limitedHostersQuotas: [String: Int64]
    /// During the trial, the remaining traffic in MB.
    public let remainingTrialQuota: Int64?
    /// Codes of the current notifications; clear them with `/v4/user/notification/clear`.
    public let notifications: [String]
    /// Whether recent links history logging is enabled.
    public let history: Bool?

    public init(username: String, email: String, isPremium: Bool, isSubscribed: Bool, isTrial: Bool, premiumUntil: Int64, lang: String, preferedDomain: String, fidelityPoints: Int, limitedHostersQuotas: [String: Int64], remainingTrialQuota: Int64? = nil, notifications: [String], history: Bool? = nil) {
        self.username = username
        self.email = email
        self.isPremium = isPremium
        self.isSubscribed = isSubscribed
        self.isTrial = isTrial
        self.premiumUntil = premiumUntil
        self.lang = lang
        self.preferedDomain = preferedDomain
        self.fidelityPoints = fidelityPoints
        self.limitedHostersQuotas = limitedHostersQuotas
        self.remainingTrialQuota = remainingTrialQuota
        self.notifications = notifications
        self.history = history
    }

    private enum CodingKeys: String, CodingKey {
        case username
        case email
        case isPremium
        case isSubscribed
        case isTrial
        case premiumUntil
        case lang
        case preferedDomain
        case fidelityPoints
        case limitedHostersQuotas
        case remainingTrialQuota
        case notifications
        case history
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        username = try c.decode(String.self, forKey: .username)
        email = try c.decode(String.self, forKey: .email)
        isPremium = try c.decode(Bool.self, forKey: .isPremium)
        isSubscribed = try c.decode(Bool.self, forKey: .isSubscribed)
        isTrial = try c.decode(Bool.self, forKey: .isTrial)
        premiumUntil = try c.decode(Int64.self, forKey: .premiumUntil)
        lang = try c.decode(String.self, forKey: .lang)
        preferedDomain = try c.decode(String.self, forKey: .preferedDomain)
        fidelityPoints = try c.decode(Int.self, forKey: .fidelityPoints)
        limitedHostersQuotas = try c.decodeQuotas(forKey: .limitedHostersQuotas)
        remainingTrialQuota = try c.decodeIfPresent(Int64.self, forKey: .remainingTrialQuota)
        notifications = try c.decode([String].self, forKey: .notifications)
        history = try c.decodeIfPresent(Bool.self, forKey: .history)
    }
}

public struct UserData: Decodable, Sendable, Equatable {
    public let user: User

    public init(user: User) {
        self.user = user
    }

    private enum CodingKeys: String, CodingKey {
        case user
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        user = try c.decode(User.self, forKey: .user)
    }
}

public struct UserHostsData: Decodable, Sendable, Equatable {
    public let hosts: ServiceMap

    public init(hosts: ServiceMap) {
        self.hosts = hosts
    }

    private enum CodingKeys: String, CodingKey {
        case hosts
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        hosts = try c.decode(ServiceMap.self, forKey: .hosts)
    }
}

public struct VerificationData: Decodable, Sendable, Equatable {
    public let verif: String
    /// When `waiting`, whether `/v4/user/verif/resend` may be called.
    public let resendable: Bool?
    /// When `allowed`, the API key.
    public let apikey: String?

    public init(verif: String, resendable: Bool? = nil, apikey: String? = nil) {
        self.verif = verif
        self.resendable = resendable
        self.apikey = apikey
    }

    private enum CodingKeys: String, CodingKey {
        case verif
        case resendable
        case apikey
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        verif = try c.decode(String.self, forKey: .verif)
        resendable = try c.decodeIfPresent(Bool.self, forKey: .resendable)
        apikey = try c.decodeIfPresent(String.self, forKey: .apikey)
    }
}

public struct ResendVerificationData: Decodable, Sendable, Equatable {
    public let sent: Bool

    public init(sent: Bool) {
        self.sent = sent
    }

    private enum CodingKeys: String, CodingKey {
        case sent
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        sent = try c.decode(Bool.self, forKey: .sent)
    }
}

public struct LinkInfo: Decodable, Sendable, Equatable {
    /// The requested link.
    public let link: String
    public let filename: String?
    /// File size in bytes.
    public let size: Int64?
    public let host: String?
    /// The host's main domain.
    public let hostDomain: String?
    public let error: ApiError?

    public init(link: String, filename: String? = nil, size: Int64? = nil, host: String? = nil, hostDomain: String? = nil, error: ApiError? = nil) {
        self.link = link
        self.filename = filename
        self.size = size
        self.host = host
        self.hostDomain = hostDomain
        self.error = error
    }

    private enum CodingKeys: String, CodingKey {
        case link
        case filename
        case size
        case host
        case hostDomain
        case error
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        link = try c.decode(String.self, forKey: .link)
        filename = try c.decodeIfPresent(String.self, forKey: .filename)
        size = try c.decodeIfPresent(Int64.self, forKey: .size)
        host = try c.decodeIfPresent(String.self, forKey: .host)
        hostDomain = try c.decodeIfPresent(String.self, forKey: .hostDomain)
        error = try c.decodeIfPresent(ApiError.self, forKey: .error)
    }
}

public struct LinkInfosData: Decodable, Sendable, Equatable {
    public let infos: [LinkInfo]

    public init(infos: [LinkInfo]) {
        self.infos = infos
    }

    private enum CodingKeys: String, CodingKey {
        case infos
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        infos = try c.decode([LinkInfo].self, forKey: .infos)
    }
}

public struct RedirectorData: Decodable, Sendable, Equatable {
    /// Encrypted links, usable with `/v4/link/infos` and `/v4/link/unlock`.
    public let links: [String]

    public init(links: [String]) {
        self.links = links
    }

    private enum CodingKeys: String, CodingKey {
        case links
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        links = try c.decode([String].self, forKey: .links)
    }
}

/// One quality of a streaming link.
public struct Stream: Decodable, Sendable, Equatable {
    /// The `stream` value for `/v4/link/streaming`.
    public let id: String
    public let ext: String?
    /// Height in pixels (e.g. `1080`), or a label such as `mp3`.
    public let quality: FlexString?
    public let filesize: Int64?
    public let proto: String?
    public let name: String?
    /// Audio-only streams. Total bitrate.
    public let tb: Double?
    /// Audio-only streams. Audio bitrate.
    public let abr: Int?

    public init(id: String, ext: String? = nil, quality: FlexString? = nil, filesize: Int64? = nil, proto: String? = nil, name: String? = nil, tb: Double? = nil, abr: Int? = nil) {
        self.id = id
        self.ext = ext
        self.quality = quality
        self.filesize = filesize
        self.proto = proto
        self.name = name
        self.tb = tb
        self.abr = abr
    }

    private enum CodingKeys: String, CodingKey {
        case id
        case ext
        case quality
        case filesize
        case proto
        case name
        case tb
        case abr
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        id = try c.decode(String.self, forKey: .id)
        ext = try c.decodeIfPresent(String.self, forKey: .ext)
        quality = try c.decodeIfPresent(FlexString.self, forKey: .quality)
        filesize = try c.decodeIfPresent(Int64.self, forKey: .filesize)
        proto = try c.decodeIfPresent(String.self, forKey: .proto)
        name = try c.decodeIfPresent(String.self, forKey: .name)
        tb = try c.decodeIfPresent(Double.self, forKey: .tb)
        abr = try c.decodeIfPresent(Int.self, forKey: .abr)
    }
}

public struct UnlockData: Decodable, Sendable, Equatable {
    /// The direct download link. Empty when `delayed` is set.
    public let link: String?
    public let host: String?
    public let filename: String?
    /// File size in bytes.
    public let filesize: Int64?
    /// The generation ID, needed by `/v4/link/streaming`.
    public let id: String?
    public let hostDomain: String?
    /// Unused.
    public let paws: Bool?
    /// Available qualities for streaming links; empty otherwise.
    public let streaming: [Stream]?
    /// Set when the link is still being generated; poll `/v4/link/delayed` with it.
    public let delayed: FlexString?

    public init(link: String? = nil, host: String? = nil, filename: String? = nil, filesize: Int64? = nil, id: String? = nil, hostDomain: String? = nil, paws: Bool? = nil, streaming: [Stream]? = nil, delayed: FlexString? = nil) {
        self.link = link
        self.host = host
        self.filename = filename
        self.filesize = filesize
        self.id = id
        self.hostDomain = hostDomain
        self.paws = paws
        self.streaming = streaming
        self.delayed = delayed
    }

    private enum CodingKeys: String, CodingKey {
        case link
        case host
        case filename
        case filesize
        case id
        case hostDomain
        case paws
        case streaming
        case delayed
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        link = try c.decodeIfPresent(String.self, forKey: .link)
        host = try c.decodeIfPresent(String.self, forKey: .host)
        filename = try c.decodeIfPresent(String.self, forKey: .filename)
        filesize = try c.decodeIfPresent(Int64.self, forKey: .filesize)
        id = try c.decodeIfPresent(String.self, forKey: .id)
        hostDomain = try c.decodeIfPresent(String.self, forKey: .hostDomain)
        paws = try c.decodeIfPresent(Bool.self, forKey: .paws)
        streaming = try c.decodeIfPresent([Stream].self, forKey: .streaming)
        delayed = try c.decodeIfPresent(FlexString.self, forKey: .delayed)
    }
}

public struct StreamingData: Decodable, Sendable, Equatable {
    /// The download link; absent when `delayed` is set.
    public let link: String?
    public let filename: String?
    public let filesize: Int64?
    /// Set when the link is still being generated; poll `/v4/link/delayed` with it.
    public let delayed: FlexString?

    public init(link: String? = nil, filename: String? = nil, filesize: Int64? = nil, delayed: FlexString? = nil) {
        self.link = link
        self.filename = filename
        self.filesize = filesize
        self.delayed = delayed
    }

    private enum CodingKeys: String, CodingKey {
        case link
        case filename
        case filesize
        case delayed
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        link = try c.decodeIfPresent(String.self, forKey: .link)
        filename = try c.decodeIfPresent(String.self, forKey: .filename)
        filesize = try c.decodeIfPresent(Int64.self, forKey: .filesize)
        delayed = try c.decodeIfPresent(FlexString.self, forKey: .delayed)
    }
}

public struct DelayedData: Decodable, Sendable, Equatable {
    /// `1` still processing, `2` ready (`link` is set), `3` failed.
    public let status: Int
    /// Estimated seconds left.
    public let timeLeft: Int
    /// The download link, once `status` is `2`.
    public let link: String?

    public init(status: Int, timeLeft: Int, link: String? = nil) {
        self.status = status
        self.timeLeft = timeLeft
        self.link = link
    }

    private enum CodingKeys: String, CodingKey {
        case status
        case timeLeft = "time_left"
        case link
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        status = try c.decode(Int.self, forKey: .status)
        timeLeft = try c.decode(Int.self, forKey: .timeLeft)
        link = try c.decodeIfPresent(String.self, forKey: .link)
    }
}

/// An uploaded magnet or torrent, or a per-item `error`.
public struct UploadedMagnet: Decodable, Sendable, Equatable {
    /// The magnet as sent (for torrent files, the magnet built from the file).
    public let magnet: String?
    /// Torrent file uploads with an `error` only. The file name sent.
    public let file: String?
    /// The info hash.
    public let hash: String?
    /// The torrent name, or `noname` when unknown.
    public let name: String?
    public let filenameOriginal: String?
    /// Total size in bytes.
    public let size: Int64?
    /// Whether the files are already available.
    public let ready: Bool?
    /// The magnet ID, for `/v4.1/magnet/status`.
    public let id: Int64?
    public let error: ApiError?

    public init(magnet: String? = nil, file: String? = nil, hash: String? = nil, name: String? = nil, filenameOriginal: String? = nil, size: Int64? = nil, ready: Bool? = nil, id: Int64? = nil, error: ApiError? = nil) {
        self.magnet = magnet
        self.file = file
        self.hash = hash
        self.name = name
        self.filenameOriginal = filenameOriginal
        self.size = size
        self.ready = ready
        self.id = id
        self.error = error
    }

    private enum CodingKeys: String, CodingKey {
        case magnet
        case file
        case hash
        case name
        case filenameOriginal = "filename_original"
        case size
        case ready
        case id
        case error
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        magnet = try c.decodeIfPresent(String.self, forKey: .magnet)
        file = try c.decodeIfPresent(String.self, forKey: .file)
        hash = try c.decodeIfPresent(String.self, forKey: .hash)
        name = try c.decodeIfPresent(String.self, forKey: .name)
        filenameOriginal = try c.decodeIfPresent(String.self, forKey: .filenameOriginal)
        size = try c.decodeIfPresent(Int64.self, forKey: .size)
        ready = try c.decodeIfPresent(Bool.self, forKey: .ready)
        id = try c.decodeIfPresent(Int64.self, forKey: .id)
        error = try c.decodeIfPresent(ApiError.self, forKey: .error)
    }
}

public struct UploadMagnetsData: Decodable, Sendable, Equatable {
    public let magnets: [UploadedMagnet]

    public init(magnets: [UploadedMagnet]) {
        self.magnets = magnets
    }

    private enum CodingKeys: String, CodingKey {
        case magnets
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        magnets = try c.decode([UploadedMagnet].self, forKey: .magnets)
    }
}

public struct UploadTorrentFilesData: Decodable, Sendable, Equatable {
    public let files: [UploadedMagnet]

    public init(files: [UploadedMagnet]) {
        self.files = files
    }

    private enum CodingKeys: String, CodingKey {
        case files
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        files = try c.decode([UploadedMagnet].self, forKey: .files)
    }
}

/// A file (`n`, `s`, `l`) or a folder (`n`, `e`) of a magnet.
public struct FileNode: Decodable, Sendable, Equatable {
    /// File or folder name.
    public let n: String
    /// Files only. Size in bytes.
    public let s: Int64?
    /// Files only. Link to unlock with `/v4/link/unlock`.
    public let l: String?
    /// Folders only. The entries of the folder.
    public let e: [FileNode]?

    public init(n: String, s: Int64? = nil, l: String? = nil, e: [FileNode]? = nil) {
        self.n = n
        self.s = s
        self.l = l
        self.e = e
    }

    private enum CodingKeys: String, CodingKey {
        case n
        case s
        case l
        case e
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        n = try c.decode(String.self, forKey: .n)
        s = try c.decodeIfPresent(Int64.self, forKey: .s)
        l = try c.decodeIfPresent(String.self, forKey: .l)
        e = try c.decodeIfPresent([FileNode].self, forKey: .e)
    }
}

/// A magnet of the account. In a live mode delta, only `id` and the
/// fields that changed are present.
public struct Magnet: Decodable, Sendable, Equatable {
    public let id: Int64
    public let filename: String?
    /// Total size in bytes.
    public let size: Int64?
    /// The info hash.
    public let hash: String?
    /// The status in plain English.
    public let status: String?
    /// `0` queued, `1` downloading, `2` compressing or moving, `3` uploading,
    /// `4` ready. `5` to `15` are errors: `5` upload failed, `6` unpacking error,
    /// `7` no peer, `8` too big, `9` internal error, `10` took too long,
    /// `11` files removed from the host, `12` and `13` processing failed,
    /// `14` tracker error, `15` no peer.
    public let statusCode: Int?
    /// Bytes downloaded so far (while processing).
    public let downloaded: Int64?
    /// Bytes uploaded so far (while processing).
    public let uploaded: Int64?
    /// Seeders (while downloading).
    public let seeders: Int?
    /// Bytes per second (while downloading).
    public let downloadSpeed: Int64?
    /// Bytes per second (while uploading).
    public let uploadSpeed: Int64?
    /// Progress percentage while processing.
    public let processingPerc: Double?
    /// Unix time the magnet was added.
    public let uploadDate: Int64?
    /// Unix time the magnet became ready, or 0.
    public let completionDate: Int64?
    /// `m` when added as a magnet, `b` when added as a torrent file.
    public let type: String?
    public let notified: Bool?
    public let version: Int?
    /// Number of files, once ready.
    public let nbLinks: Int?
    /// Only when a single magnet is requested with `id`.
    public let files: [FileNode]?
    /// Live mode deltas only. The magnet was deleted; drop it from the local state.
    public let deleted: Bool?

    public init(id: Int64, filename: String? = nil, size: Int64? = nil, hash: String? = nil, status: String? = nil, statusCode: Int? = nil, downloaded: Int64? = nil, uploaded: Int64? = nil, seeders: Int? = nil, downloadSpeed: Int64? = nil, uploadSpeed: Int64? = nil, processingPerc: Double? = nil, uploadDate: Int64? = nil, completionDate: Int64? = nil, type: String? = nil, notified: Bool? = nil, version: Int? = nil, nbLinks: Int? = nil, files: [FileNode]? = nil, deleted: Bool? = nil) {
        self.id = id
        self.filename = filename
        self.size = size
        self.hash = hash
        self.status = status
        self.statusCode = statusCode
        self.downloaded = downloaded
        self.uploaded = uploaded
        self.seeders = seeders
        self.downloadSpeed = downloadSpeed
        self.uploadSpeed = uploadSpeed
        self.processingPerc = processingPerc
        self.uploadDate = uploadDate
        self.completionDate = completionDate
        self.type = type
        self.notified = notified
        self.version = version
        self.nbLinks = nbLinks
        self.files = files
        self.deleted = deleted
    }

    private enum CodingKeys: String, CodingKey {
        case id
        case filename
        case size
        case hash
        case status
        case statusCode
        case downloaded
        case uploaded
        case seeders
        case downloadSpeed
        case uploadSpeed
        case processingPerc
        case uploadDate
        case completionDate
        case type
        case notified
        case version
        case nbLinks
        case files
        case deleted
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        id = try c.decode(Int64.self, forKey: .id)
        filename = try c.decodeIfPresent(String.self, forKey: .filename)
        size = try c.decodeIfPresent(Int64.self, forKey: .size)
        hash = try c.decodeIfPresent(String.self, forKey: .hash)
        status = try c.decodeIfPresent(String.self, forKey: .status)
        statusCode = try c.decodeIfPresent(Int.self, forKey: .statusCode)
        downloaded = try c.decodeIfPresent(Int64.self, forKey: .downloaded)
        uploaded = try c.decodeIfPresent(Int64.self, forKey: .uploaded)
        seeders = try c.decodeIfPresent(Int.self, forKey: .seeders)
        downloadSpeed = try c.decodeIfPresent(Int64.self, forKey: .downloadSpeed)
        uploadSpeed = try c.decodeIfPresent(Int64.self, forKey: .uploadSpeed)
        processingPerc = try c.decodeIfPresent(Double.self, forKey: .processingPerc)
        uploadDate = try c.decodeIfPresent(Int64.self, forKey: .uploadDate)
        completionDate = try c.decodeIfPresent(Int64.self, forKey: .completionDate)
        type = try c.decodeIfPresent(String.self, forKey: .type)
        notified = try c.decodeIfPresent(Bool.self, forKey: .notified)
        version = try c.decodeIfPresent(Int.self, forKey: .version)
        nbLinks = try c.decodeIfPresent(Int.self, forKey: .nbLinks)
        files = try c.decodeIfPresent([FileNode].self, forKey: .files)
        deleted = try c.decodeIfPresent(Bool.self, forKey: .deleted)
    }
}

public struct MagnetStatusData: Decodable, Sendable, Equatable {
    public let magnets: [Magnet]
    /// Live mode. The counter to send on the next call.
    public let counter: Int64?
    /// Live mode. `true` when the response holds the full list.
    public let fullsync: Bool?

    public init(magnets: [Magnet], counter: Int64? = nil, fullsync: Bool? = nil) {
        self.magnets = magnets
        self.counter = counter
        self.fullsync = fullsync
    }

    private enum CodingKeys: String, CodingKey {
        case magnets
        case counter
        case fullsync
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        magnets = try c.decodeMagnets(forKey: .magnets)
        counter = try c.decodeIfPresent(Int64.self, forKey: .counter)
        fullsync = try c.decodeIfPresent(Bool.self, forKey: .fullsync)
    }
}

public struct MagnetFiles: Decodable, Sendable, Equatable {
    /// The magnet ID (sent back as a string).
    public let id: FlexString
    public let files: [FileNode]?
    public let error: ApiError?

    public init(id: FlexString, files: [FileNode]? = nil, error: ApiError? = nil) {
        self.id = id
        self.files = files
        self.error = error
    }

    private enum CodingKeys: String, CodingKey {
        case id
        case files
        case error
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        id = try c.decode(FlexString.self, forKey: .id)
        files = try c.decodeIfPresent([FileNode].self, forKey: .files)
        error = try c.decodeIfPresent(ApiError.self, forKey: .error)
    }
}

public struct MagnetFilesData: Decodable, Sendable, Equatable {
    public let magnets: [MagnetFiles]

    public init(magnets: [MagnetFiles]) {
        self.magnets = magnets
    }

    private enum CodingKeys: String, CodingKey {
        case magnets
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        magnets = try c.decode([MagnetFiles].self, forKey: .magnets)
    }
}

public struct RestartedMagnet: Decodable, Sendable, Equatable {
    /// The magnet ID.
    public let magnet: String?
    public let message: String?
    public let error: ApiError?

    public init(magnet: String? = nil, message: String? = nil, error: ApiError? = nil) {
        self.magnet = magnet
        self.message = message
        self.error = error
    }

    private enum CodingKeys: String, CodingKey {
        case magnet
        case message
        case error
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        magnet = try c.decodeIfPresent(String.self, forKey: .magnet)
        message = try c.decodeIfPresent(String.self, forKey: .message)
        error = try c.decodeIfPresent(ApiError.self, forKey: .error)
    }
}

public struct RestartMagnetsData: Decodable, Sendable, Equatable {
    /// With `id`. The restarted magnet.
    public let id: Int64?
    /// With `id`. A confirmation message.
    public let message: String?
    /// With `ids[]`. One entry per magnet.
    public let magnets: [RestartedMagnet]?

    public init(id: Int64? = nil, message: String? = nil, magnets: [RestartedMagnet]? = nil) {
        self.id = id
        self.message = message
        self.magnets = magnets
    }

    private enum CodingKeys: String, CodingKey {
        case id
        case message
        case magnets
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        id = try c.decodeIfPresent(Int64.self, forKey: .id)
        message = try c.decodeIfPresent(String.self, forKey: .message)
        magnets = try c.decodeIfPresent([RestartedMagnet].self, forKey: .magnets)
    }
}

public struct SavedLink: Decodable, Sendable, Equatable {
    public let link: String
    public let filename: String?
    public let size: Int64?
    /// Unix time the link was saved or unlocked.
    public let date: Int64?
    public let host: String?
    /// History only. The download link that was generated.
    public let linkDl: String?
    /// History only. The error the unlock gave, if any.
    public let error: String?

    public init(link: String, filename: String? = nil, size: Int64? = nil, date: Int64? = nil, host: String? = nil, linkDl: String? = nil, error: String? = nil) {
        self.link = link
        self.filename = filename
        self.size = size
        self.date = date
        self.host = host
        self.linkDl = linkDl
        self.error = error
    }

    private enum CodingKeys: String, CodingKey {
        case link
        case filename
        case size
        case date
        case host
        case linkDl = "link_dl"
        case error
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        link = try c.decode(String.self, forKey: .link)
        filename = try c.decodeIfPresent(String.self, forKey: .filename)
        size = try c.decodeIfPresent(Int64.self, forKey: .size)
        date = try c.decodeIfPresent(Int64.self, forKey: .date)
        host = try c.decodeIfPresent(String.self, forKey: .host)
        linkDl = try c.decodeIfPresent(String.self, forKey: .linkDl)
        error = try c.decodeIfPresent(String.self, forKey: .error)
    }
}

public struct SavedLinksData: Decodable, Sendable, Equatable {
    public let links: [SavedLink]

    public init(links: [SavedLink]) {
        self.links = links
    }

    private enum CodingKeys: String, CodingKey {
        case links
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        links = try c.decode([SavedLink].self, forKey: .links)
    }
}

public struct VoucherBalanceData: Decodable, Sendable, Equatable {
    /// The reseller balance.
    public let balance: Double

    public init(balance: Double) {
        self.balance = balance
    }

    private enum CodingKeys: String, CodingKey {
        case balance
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        balance = try c.decode(Double.self, forKey: .balance)
    }
}

public struct VouchersData: Decodable, Sendable, Equatable {
    public let codes: [String]
    /// `true` when fewer vouchers than requested were available.
    public let partialList: Bool?

    public init(codes: [String], partialList: Bool? = nil) {
        self.codes = codes
        self.partialList = partialList
    }

    private enum CodingKeys: String, CodingKey {
        case codes
        case partialList
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        codes = try c.decode([String].self, forKey: .codes)
        partialList = try c.decodeIfPresent(Bool.self, forKey: .partialList)
    }
}

public struct GeneratedVouchersData: Decodable, Sendable, Equatable {
    public let codes: [String]
    public let pricePerVoucher: Double?
    /// Total price of the request.
    public let total: Double?
    /// The balance left after the generation.
    public let balance: Double?

    public init(codes: [String], pricePerVoucher: Double? = nil, total: Double? = nil, balance: Double? = nil) {
        self.codes = codes
        self.pricePerVoucher = pricePerVoucher
        self.total = total
        self.balance = balance
    }

    private enum CodingKeys: String, CodingKey {
        case codes
        case pricePerVoucher
        case total
        case balance
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        codes = try c.decode([String].self, forKey: .codes)
        pricePerVoucher = try c.decodeIfPresent(Double.self, forKey: .pricePerVoucher)
        total = try c.decodeIfPresent(Double.self, forKey: .total)
        balance = try c.decodeIfPresent(Double.self, forKey: .balance)
    }
}

/// Error codes the API is known to return. The API may return others.
public enum ErrorCode {
    public static let generic = "GENERIC"
    public static let notFound = "404"
    public static let maintenance = "MAINTENANCE"
    public static let discontinued = "DISCONTINUED"
    public static let authMissingApikey = "AUTH_MISSING_APIKEY"
    public static let authBadApikey = "AUTH_BAD_APIKEY"
    public static let authBlocked = "AUTH_BLOCKED"
    public static let authUserBanned = "AUTH_USER_BANNED"
    public static let alreadySent = "ALREADY_SENT"
    public static let noServer = "NO_SERVER"
    public static let linkIsMissing = "LINK_IS_MISSING"
    public static let badLink = "BAD_LINK"
    public static let linkHostNotSupported = "LINK_HOST_NOT_SUPPORTED"
    public static let linkDown = "LINK_DOWN"
    public static let linkPassProtected = "LINK_PASS_PROTECTED"
    public static let linkHostUnavailable = "LINK_HOST_UNAVAILABLE"
    public static let linkTooManyDownloads = "LINK_TOO_MANY_DOWNLOADS"
    public static let linkHostFull = "LINK_HOST_FULL"
    public static let linkHostLimitReached = "LINK_HOST_LIMIT_REACHED"
    public static let linkError = "LINK_ERROR"
    public static let linkTemporaryUnavailable = "LINK_TEMPORARY_UNAVAILABLE"
    public static let linkNotSupported = "LINK_NOT_SUPPORTED"
    public static let redirectorNotSupported = "REDIRECTOR_NOT_SUPPORTED"
    public static let redirectorError = "REDIRECTOR_ERROR"
    public static let redirectorRateLimited = "REDIRECTOR_RATE_LIMITED"
    public static let streamInvalidGenId = "STREAM_INVALID_GEN_ID"
    public static let streamInvalidStreamId = "STREAM_INVALID_STREAM_ID"
    public static let delayedInvalidId = "DELAYED_INVALID_ID"
    public static let freeTrialLimitReached = "FREE_TRIAL_LIMIT_REACHED"
    public static let mustBePremium = "MUST_BE_PREMIUM"
    public static let magnetInvalidId = "MAGNET_INVALID_ID"
    public static let magnetInvalidUri = "MAGNET_INVALID_URI"
    public static let magnetInvalidFile = "MAGNET_INVALID_FILE"
    public static let magnetFileUploadFailed = "MAGNET_FILE_UPLOAD_FAILED"
    public static let magnetNoUri = "MAGNET_NO_URI"
    public static let magnetProcessing = "MAGNET_PROCESSING"
    public static let magnetTooManyActive = "MAGNET_TOO_MANY_ACTIVE"
    public static let magnetTooMany = "MAGNET_TOO_MANY"
    public static let magnetMustBePremium = "MAGNET_MUST_BE_PREMIUM"
    public static let magnetNoServer = "MAGNET_NO_SERVER"
    public static let magnetTooLarge = "MAGNET_TOO_LARGE"
    public static let magnetUploadFailed = "MAGNET_UPLOAD_FAILED"
    public static let magnetInternalError = "MAGNET_INTERNAL_ERROR"
    public static let magnetCantBootstrap = "MAGNET_CANT_BOOTSTRAP"
    public static let magnetMagnetTooBig = "MAGNET_MAGNET_TOO_BIG"
    public static let magnetTookTooLong = "MAGNET_TOOK_TOO_LONG"
    public static let magnetLinksRemoved = "MAGNET_LINKS_REMOVED"
    public static let magnetProcessingFailed = "MAGNET_PROCESSING_FAILED"
    public static let pinAlreadyAuthed = "PIN_ALREADY_AUTHED"
    public static let pinExpired = "PIN_EXPIRED"
    public static let pinInvalid = "PIN_INVALID"
    public static let userLinkMissing = "USER_LINK_MISSING"
    public static let userLinkInvalid = "USER_LINK_INVALID"
    public static let missingNotifEndpoint = "MISSING_NOTIF_ENDPOINT"
    public static let voucherDurationInvalid = "VOUCHER_DURATION_INVALID"
    public static let voucherNbInvalid = "VOUCHER_NB_INVALID"
    public static let noMoreVoucher = "NO_MORE_VOUCHER"
    public static let insufficientBalance = "INSUFFICIENT_BALANCE"
    public static let downloadFailed = "DOWNLOAD_FAILED"
    public static let accountInvalid = "ACCOUNT_INVALID"
    public static let noJsonParam = "NO_JSON_PARAM"
    public static let jsonInvalid = "JSON_INVALID"
    public static let freedaysInvalidCountry = "FREEDAYS_INVALID_COUNTRY"
    public static let freedaysInvalidPhone = "FREEDAYS_INVALID_PHONE"
    public static let freedaysInvalidProvider = "FREEDAYS_INVALID_PROVIDER"
    public static let freedaysUsedPhone = "FREEDAYS_USED_PHONE"
    public static let freedaysAlreadyRequested = "FREEDAYS_ALREADY_REQUESTED"
    public static let freedaysInvalidStatus = "FREEDAYS_INVALID_STATUS"
    public static let freedaysTooMuchRetries = "FREEDAYS_TOO_MUCH_RETRIES"
    public static let freedaysInvalidCode = "FREEDAYS_INVALID_CODE"
    public static let freedaysDenied = "FREEDAYS_DENIED"
    public static let freedaysErrorSending = "FREEDAYS_ERROR_SENDING"

    /// Every code of the spec.
    public static let all: [String] = [
        generic,
        notFound,
        maintenance,
        discontinued,
        authMissingApikey,
        authBadApikey,
        authBlocked,
        authUserBanned,
        alreadySent,
        noServer,
        linkIsMissing,
        badLink,
        linkHostNotSupported,
        linkDown,
        linkPassProtected,
        linkHostUnavailable,
        linkTooManyDownloads,
        linkHostFull,
        linkHostLimitReached,
        linkError,
        linkTemporaryUnavailable,
        linkNotSupported,
        redirectorNotSupported,
        redirectorError,
        redirectorRateLimited,
        streamInvalidGenId,
        streamInvalidStreamId,
        delayedInvalidId,
        freeTrialLimitReached,
        mustBePremium,
        magnetInvalidId,
        magnetInvalidUri,
        magnetInvalidFile,
        magnetFileUploadFailed,
        magnetNoUri,
        magnetProcessing,
        magnetTooManyActive,
        magnetTooMany,
        magnetMustBePremium,
        magnetNoServer,
        magnetTooLarge,
        magnetUploadFailed,
        magnetInternalError,
        magnetCantBootstrap,
        magnetMagnetTooBig,
        magnetTookTooLong,
        magnetLinksRemoved,
        magnetProcessingFailed,
        pinAlreadyAuthed,
        pinExpired,
        pinInvalid,
        userLinkMissing,
        userLinkInvalid,
        missingNotifEndpoint,
        voucherDurationInvalid,
        voucherNbInvalid,
        noMoreVoucher,
        insufficientBalance,
        downloadFailed,
        accountInvalid,
        noJsonParam,
        jsonInvalid,
        freedaysInvalidCountry,
        freedaysInvalidPhone,
        freedaysInvalidProvider,
        freedaysUsedPhone,
        freedaysAlreadyRequested,
        freedaysInvalidStatus,
        freedaysTooMuchRetries,
        freedaysInvalidCode,
        freedaysDenied,
        freedaysErrorSending,
    ]
}
