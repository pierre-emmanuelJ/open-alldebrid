import Foundation
#if canImport(FoundationNetworking)
    import FoundationNetworking
#endif

/// A .torrent file to upload.
public struct TorrentFile: Sendable {
    public let name: String
    public let content: Data

    public init(name: String, content: Data) {
        self.name = name
        self.content = content
    }
}

/// Magnet status codes. 5 and above are errors.
public enum MagnetStatus {
    public static let queued = 0
    public static let downloading = 1
    public static let compressing = 2
    public static let uploading = 3
    public static let ready = 4
}

/// Delayed link statuses.
public enum DelayedStatus {
    public static let processing = 1
    public static let ready = 2
    public static let failed = 3
}

extension Magnet {
    /// Whether the magnet's files are available.
    public var isReady: Bool { statusCode == MagnetStatus.ready }

    /// Whether the magnet failed.
    public var isError: Bool { (statusCode ?? 0) > MagnetStatus.ready }
}

/// A client for the AllDebrid API.
///
/// ```swift
/// let client = AllDebrid(apiKey: apiKey, userAgent: "my-app")
/// let user = try await client.user()
/// ```
///
/// Calls throw ``AllDebridError``.
public struct AllDebrid: Sendable {
    public static let defaultBaseURL = URL(string: "https://api.alldebrid.com")!

    /// Public demo keys: the API answers them with static sample data.
    public enum DemoKey {
        public static let premium = "staticDemoApikeyPrem"
        public static let trial = "staticDemoApikeyTria"
        public static let free = "staticDemoApikeyFree"
    }

    private let apiKey: String?
    private let baseURL: URL
    private let userAgent: String
    private let session: URLSession

    /// - Parameters:
    ///   - apiKey: An API key; nil for public endpoints and the PIN flow.
    ///   - baseURL: The API base URL, for tests or a proxy.
    ///   - userAgent: Identifies your application to AllDebrid.
    ///   - session: The URL session.
    public init(
        apiKey: String? = nil,
        baseURL: URL = AllDebrid.defaultBaseURL,
        userAgent: String = "open-alldebrid-swift",
        session: URLSession = .shared
    ) {
        self.apiKey = apiKey?.isEmpty == false ? apiKey : nil
        self.baseURL = baseURL
        self.userAgent = userAgent
        self.session = session
    }

    /// A client with another API key, sharing this one's session.
    public func withAPIKey(_ apiKey: String?) -> AllDebrid {
        AllDebrid(apiKey: apiKey, baseURL: baseURL, userAgent: userAgent, session: session)
    }

    typealias Params = [(String, String)]

    /// Calls the API and returns the JSON of the `data` of a successful answer.
    func callRaw(_ method: String, _ path: String, _ params: Params = [], files: [TorrentFile]? = nil) async throws -> Data {
        var components = URLComponents(url: baseURL.appendingPathComponent(path), resolvingAgainstBaseURL: false)!
        if method == "GET", !params.isEmpty {
            components.queryItems = params.map { URLQueryItem(name: $0.0, value: $0.1) }
        }
        var request = URLRequest(url: components.url!)
        request.httpMethod = method
        request.setValue("application/json", forHTTPHeaderField: "Accept")
        request.setValue(userAgent, forHTTPHeaderField: "User-Agent")
        if let apiKey {
            request.setValue("Bearer \(apiKey)", forHTTPHeaderField: "Authorization")
        }
        if let files {
            let boundary = "open-alldebrid-\(UUID().uuidString)"
            request.setValue("multipart/form-data; boundary=\(boundary)", forHTTPHeaderField: "Content-Type")
            request.httpBody = Self.multipart(files, boundary: boundary)
        } else if method != "GET" {
            // List parameters keep their PHP-style names and repeat: magnets[]=a&magnets[]=b.
            request.setValue("application/x-www-form-urlencoded", forHTTPHeaderField: "Content-Type")
            request.httpBody = Data(Self.formEncode(params).utf8)
        }

        let (body, response) = try await session.data(for: request)
        let status = (response as? HTTPURLResponse)?.statusCode ?? 0
        let envelope = (try? JSONSerialization.jsonObject(with: body)) as? [String: Any]
        let apiError = Self.apiError(envelope?["error"])
        if status != 200 {
            throw AllDebridError.http(status: status, code: apiError?.code)
        }
        guard let envelope else {
            throw AllDebridError.decoding("\(path): the response is not a JSON object")
        }
        guard envelope["status"] as? String == "success" else {
            throw apiError ?? AllDebridError.decoding("\(path): error without details")
        }
        let data = envelope["data"] ?? NSNull()
        return try JSONSerialization.data(withJSONObject: data, options: [.fragmentsAllowed])
    }

    private func call<T: Decodable>(_ method: String, _ path: String, _ params: Params = [], files: [TorrentFile]? = nil) async throws -> T
    {
        let data = try await callRaw(method, path, params, files: files)
        do {
            return try JSONDecoder().decode(T.self, from: data)
        } catch {
            throw AllDebridError.decoding("\(path): \(error)")
        }
    }

    private static func apiError(_ value: Any?) -> AllDebridError? {
        guard let error = value as? [String: Any] else { return nil }
        let text = { (key: String) -> String? in error[key].map { "\($0)" } }
        return .api(code: text("code") ?? "", message: text("message") ?? "", token: text("token"))
    }

    static func formEncode(_ params: Params) -> String {
        var allowed = CharacterSet.alphanumerics
        allowed.insert(charactersIn: "-._~")
        let encode = { (s: String) in s.addingPercentEncoding(withAllowedCharacters: allowed) ?? s }
        return params.map { "\(encode($0.0))=\(encode($0.1))" }.joined(separator: "&")
    }

    private static func multipart(_ files: [TorrentFile], boundary: String) -> Data {
        var body = Data()
        for file in files {
            let name = file.name.replacingOccurrences(of: "\"", with: "\\\"")
            body.append(Data("--\(boundary)\r\n".utf8))
            body.append(Data("Content-Disposition: form-data; name=\"files[]\"; filename=\"\(name)\"\r\n".utf8))
            body.append(Data("Content-Type: application/x-bittorrent\r\n\r\n".utf8))
            body.append(file.content)
            body.append(Data("\r\n".utf8))
        }
        body.append(Data("--\(boundary)--\r\n".utf8))
        return body
    }

    // MARK: - General

    /// Checks that the API is up.
    public func ping() async throws {
        _ = try await callRaw("GET", "/v4/ping")
    }

    // MARK: - PIN

    /// Starts a PIN authentication. Show `pin` or `userUrl`, then ``waitForPin(_:interval:)``.
    public func getPin() async throws -> PinGetData {
        try await call("GET", "/v4.1/pin/get")
    }

    /// The state of a PIN authentication; `apikey` is set once activated.
    public func checkPin(pin: String, check: String) async throws -> PinCheckData {
        try await call("POST", "/v4/pin/check", [("pin", pin), ("check", check)])
    }

    /// Polls until the user confirms the PIN and returns the new API key.
    public func waitForPin(_ pin: PinGetData, interval: Duration = .seconds(5)) async throws -> String {
        while true {
            let state = try await checkPin(pin: pin.pin, check: pin.check)
            if state.activated, let key = state.apikey, !key.isEmpty {
                return key
            }
            try await Task.sleep(for: interval)
        }
    }

    // MARK: - Hosts

    /// Supported hosts, streaming sites and redirectors. Cache the result.
    public func hosts(hostsOnly: Bool = false) async throws -> HostsData {
        try await call("GET", "/v4/hosts", hostsOnly ? [("hostsOnly", "true")] : [])
    }

    /// Every supported domain.
    public func hostDomains() async throws -> HostDomainsData {
        try await call("GET", "/v4/hosts/domains")
    }

    /// Host names ranked from the most open (1) to the most restricted.
    public func hostPriority() async throws -> [String: Int] {
        let data: HostPriorityData = try await call("GET", "/v4/hosts/priority")
        return data.hosts
    }

    // MARK: - User

    /// The account information.
    public func user() async throws -> User {
        let data: UserData = try await call("GET", "/v4/user")
        return data.user
    }

    /// The hosts this user can use now, with their quotas.
    public func userHosts(hostsOnly: Bool = false) async throws -> ServiceMap {
        let data: UserHostsData = try await call("GET", "/v4.1/user/hosts", hostsOnly ? [("hostsOnly", "true")] : [])
        return data.hosts
    }

    /// Follows a new-location login verification (token from an `AUTH_BLOCKED` error).
    public func verificationStatus(token: String) async throws -> VerificationData {
        try await call("POST", "/v4/user/verif", [("token", token)])
    }

    /// Sends the verification email again (once per token).
    public func resendVerification(token: String) async throws {
        _ = try await callRaw("POST", "/v4/user/verif/resend", [("token", token)])
    }

    /// Clears one of the codes in `User.notifications`.
    public func clearNotification(code: String) async throws {
        _ = try await callRaw("POST", "/v4/user/notification/clear", [("code", code)])
    }

    // MARK: - Links

    /// File name and size of each link, in order; unreadable links have an `error`.
    public func linkInfos(_ links: [String], password: String? = nil) async throws -> [LinkInfo] {
        var params = links.map { ("link[]", $0) }
        if let password { params.append(("password", password)) }
        let data: LinkInfosData = try await call("POST", "/v4/link/infos", params)
        return data.infos
    }

    /// Links behind a redirector; they only work with ``linkInfos(_:password:)`` and ``unlock(_:password:)``.
    public func extractRedirector(_ link: String) async throws -> [String] {
        let data: RedirectorData = try await call("POST", "/v4/link/redirector", [("link", link)])
        return data.links
    }

    /// Unlocks a hoster link. When `delayed` is set, wait with ``waitDelayedLink(_:interval:)``;
    /// when `streaming` lists qualities, choose one with ``streamingLink(id:stream:)``.
    public func unlock(_ link: String, password: String? = nil) async throws -> UnlockData {
        var params = [("link", link)]
        if let password { params.append(("password", password)) }
        return try await call("POST", "/v4/link/unlock", params)
    }

    /// The link of one quality: `id` from ``unlock(_:password:)``, `stream` from one of its streaming entries.
    public func streamingLink(id: String, stream: String) async throws -> StreamingData {
        try await call("POST", "/v4/link/streaming", [("id", id), ("stream", stream)])
    }

    /// The status of a delayed link.
    public func delayedLink(_ id: String) async throws -> DelayedData {
        try await call("POST", "/v4/link/delayed", [("id", id)])
    }

    /// Polls a delayed link until it is ready and returns the download link.
    public func waitDelayedLink(_ id: String, interval: Duration = .seconds(5)) async throws -> String {
        while true {
            let state = try await delayedLink(id)
            if state.status == DelayedStatus.ready, let link = state.link, !link.isEmpty {
                return link
            }
            if state.status == DelayedStatus.failed {
                throw AllDebridError.decoding("delayed link \(id) failed")
            }
            try await Task.sleep(for: interval)
        }
    }

    // MARK: - Magnets

    /// Adds magnet URIs or info hashes. Results are in order; invalid ones have an `error`.
    public func uploadMagnets(_ magnets: [String]) async throws -> [UploadedMagnet] {
        let data: UploadMagnetsData = try await call("POST", "/v4/magnet/upload", magnets.map { ("magnets[]", $0) })
        return data.magnets
    }

    /// Adds .torrent files. Results are in order; invalid files have an `error`.
    public func uploadTorrentFiles(_ files: [TorrentFile]) async throws -> [UploadedMagnet] {
        let data: UploadTorrentFilesData = try await call("POST", "/v4/magnet/upload/file", files: files)
        return data.files
    }

    /// Every magnet of the account.
    public func magnets() async throws -> [Magnet] {
        let data: MagnetStatusData = try await call("POST", "/v4.1/magnet/status")
        return data.magnets
    }

    /// One magnet, with its files.
    public func magnet(id: Int64) async throws -> Magnet {
        let data: MagnetStatusData = try await call("POST", "/v4.1/magnet/status", [("id", String(id))])
        guard data.magnets.count == 1 else {
            throw AllDebridError.decoding("magnet \(id): got \(data.magnets.count) magnets")
        }
        return data.magnets[0]
    }

    /// The file tree of each magnet, in order; unknown IDs have an `error`.
    public func magnetFiles(ids: [Int64]) async throws -> [MagnetFiles] {
        let data: MagnetFilesData = try await call("POST", "/v4/magnet/files", ids.map { ("id[]", String($0)) })
        return data.magnets
    }

    /// Deletes a magnet.
    public func deleteMagnet(id: Int64) async throws {
        _ = try await callRaw("POST", "/v4/magnet/delete", [("id", String(id))])
    }

    /// Restarts a failed magnet.
    public func restartMagnet(id: Int64) async throws {
        _ = try await callRaw("POST", "/v4/magnet/restart", [("id", String(id))])
    }

    /// Restarts failed magnets; each result has its own `error`.
    public func restartMagnets(ids: [Int64]) async throws -> [RestartedMagnet] {
        let data: RestartMagnetsData = try await call("POST", "/v4/magnet/restart", ids.map { ("ids[]", String($0)) })
        return data.magnets ?? []
    }

    /// Follows the magnets with live mode: each ``LiveMagnets/poll()`` only transfers changes.
    public func liveMagnets() -> LiveMagnets {
        LiveMagnets(client: self)
    }

    // MARK: - Saved links

    /// The links the user saved.
    public func savedLinks() async throws -> [SavedLink] {
        let data: SavedLinksData = try await call("GET", "/v4/user/links")
        return data.links
    }

    /// Saves links for later.
    public func saveLinks(_ links: [String]) async throws {
        _ = try await callRaw("POST", "/v4/user/links/save", links.map { ("links[]", $0) })
    }

    /// Deletes saved links.
    public func deleteSavedLinks(_ links: [String]) async throws {
        _ = try await callRaw("POST", "/v4/user/links/delete", links.map { ("links[]", $0) })
    }

    /// Links unlocked in the last 3 days, when history logging is enabled.
    public func history() async throws -> [SavedLink] {
        let data: SavedLinksData = try await call("GET", "/v4/user/history")
        return data.links
    }

    /// Deletes the whole recent links history.
    public func purgeHistory() async throws {
        _ = try await callRaw("POST", "/v4/user/history/delete")
    }

    // MARK: - Reseller

    /// The reseller balance. Reseller accounts only.
    public func voucherBalance() async throws -> Double {
        let data: VoucherBalanceData = try await call("GET", "/v4/voucher/balance")
        return data.balance
    }

    /// Takes vouchers (15, 30, 90, 180 or 365 days) from the reseller stock.
    public func vouchers(durationDays: Int, count: Int) async throws -> VouchersData {
        try await call("POST", "/v4/voucher/get", [("duration", String(durationDays)), ("nb", String(count))])
    }

    /// Generates new vouchers from the reseller balance.
    public func generateVouchers(durationDays: Int, count: Int) async throws -> GeneratedVouchersData {
        try await call("POST", "/v4/voucher/generate", [("duration", String(durationDays)), ("nb", String(count))])
    }
}

/// Live mode state: the merged list of magnets, updated by ``poll()``.
public actor LiveMagnets {
    private let client: AllDebrid
    /// The session ID sent to the API.
    public let session = Int64.random(in: 0..<(1 << 31))
    private var counter: Int64 = 0
    private var state: [Int64: [String: Any]] = [:]

    init(client: AllDebrid) {
        self.client = client
    }

    /// Applies the changes since the last call and returns every magnet, newest first.
    public func poll() async throws -> [Magnet] {
        let raw = try await client.callRaw("POST", "/v4.1/magnet/status", [("session", String(session)), ("counter", String(counter))])
        guard let data = try JSONSerialization.jsonObject(with: raw) as? [String: Any] else {
            throw AllDebridError.decoding("live mode answered without data")
        }
        if data["fullsync"] as? Bool == true {
            state = [:]
        }
        let changes: [[String: Any]]
        switch data["magnets"] {
        case let list as [[String: Any]]: changes = list
        case let one as [String: Any]: changes = [one]
        default: changes = []
        }
        for change in changes {
            guard let id = (change["id"] as? NSNumber)?.int64Value else { continue }
            if change["deleted"] as? Bool == true {
                state[id] = nil
            } else {
                state[id, default: [:]].merge(change) { _, new in new }
            }
        }
        counter = (data["counter"] as? NSNumber)?.int64Value ?? 0

        let decoder = JSONDecoder()
        return try state.values
            .map { try decoder.decode(Magnet.self, from: JSONSerialization.data(withJSONObject: $0)) }
            .sorted { ($0.uploadDate ?? 0, $0.id) > ($1.uploadDate ?? 0, $1.id) }
    }
}
