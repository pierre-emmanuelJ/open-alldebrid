import Foundation
import Testing
#if canImport(FoundationNetworking)
    import FoundationNetworking
#endif
@testable import OpenAllDebrid

/// Records requests and answers canned bodies in order.
final class StubProtocol: URLProtocol {
    struct Recorded {
        let method: String
        let url: URL
        let headers: [String: String]
        let body: Data
    }

    nonisolated(unsafe) static var answers: [(Int, String)] = []
    nonisolated(unsafe) static var requests: [Recorded] = []
    static let lock = NSLock()

    static func reset() {
        lock.withLock {
            answers = []
            requests = []
        }
    }

    static func answer(_ status: Int, _ body: String) {
        lock.withLock { answers.append((status, body)) }
    }

    static func ok(_ data: String) {
        answer(200, #"{"status":"success","data":"# + data + "}")
    }

    static var last: Recorded { lock.withLock { requests.last! } }

    override class func canInit(with request: URLRequest) -> Bool { true }
    override class func canonicalRequest(for request: URLRequest) -> URLRequest { request }

    override func startLoading() {
        var body = request.httpBody ?? Data()
        if body.isEmpty, let stream = request.httpBodyStream {
            stream.open()
            var buffer = [UInt8](repeating: 0, count: 4096)
            while stream.hasBytesAvailable {
                let n = stream.read(&buffer, maxLength: buffer.count)
                if n <= 0 { break }
                body.append(buffer, count: n)
            }
            stream.close()
        }
        let (status, text) = Self.lock.withLock { () -> (Int, String) in
            Self.requests.append(
                Recorded(
                    method: request.httpMethod ?? "GET",
                    url: request.url!,
                    headers: request.allHTTPHeaderFields ?? [:],
                    body: body))
            return Self.answers.isEmpty ? (200, #"{"status":"success","data":{}}"#) : Self.answers.removeFirst()
        }
        let response = HTTPURLResponse(url: request.url!, statusCode: status, httpVersion: "HTTP/1.1", headerFields: nil)!
        client?.urlProtocol(self, didReceive: response, cacheStoragePolicy: .notAllowed)
        client?.urlProtocol(self, didLoad: Data(text.utf8))
        client?.urlProtocolDidFinishLoading(self)
    }

    override func stopLoading() {}
}

func form(_ r: StubProtocol.Recorded) -> [String: [String]] {
    var out: [String: [String]] = [:]
    for pair in String(decoding: r.body, as: UTF8.self).split(separator: "&") {
        let kv = pair.split(separator: "=", maxSplits: 1).map { String($0).removingPercentEncoding ?? String($0) }
        out[kv[0], default: []].append(kv.count > 1 ? kv[1] : "")
    }
    return out
}

let specPath = URL(fileURLWithPath: #filePath)
    .deletingLastPathComponent().deletingLastPathComponent().deletingLastPathComponent()
    .deletingLastPathComponent().deletingLastPathComponent()
    .appendingPathComponent("openapi.yaml")

func specOperations() throws -> Set<String> {
    let text = try String(contentsOf: specPath, encoding: .utf8)
    var ops = Set<String>()
    var path: String?
    for line in text.split(separator: "\n", omittingEmptySubsequences: false).map(String.init) {
        if line == "components:" { break }
        if line.hasPrefix("  /"), line.hasSuffix(":"), !line.hasPrefix("   ") {
            path = String(line.dropFirst(2).dropLast())
        } else if let p = path, line == "    get:" || line == "    post:" {
            ops.insert("\(line.trimmingCharacters(in: .whitespaces).dropLast().uppercased()) \(p)")
        }
    }
    return ops
}

let userJSON =
    #"{"user":{"username":"u","email":"e","isPremium":true,"isSubscribed":false,"isTrial":false,"premiumUntil":0,"lang":"en","preferedDomain":"com","fidelityPoints":0,"limitedHostersQuotas":[],"notifications":[]}}"#

@Suite(.serialized)
struct ClientTests {
    let client: AllDebrid

    init() {
        StubProtocol.reset()
        let config = URLSessionConfiguration.ephemeral
        config.protocolClasses = [StubProtocol.self]
        client = AllDebrid(
            apiKey: "test-key", baseURL: URL(string: "https://api.test")!, userAgent: "unit",
            session: URLSession(configuration: config))
    }

    @Test func listParametersArePHPStyle() async throws {
        StubProtocol.ok(#"{"magnets":[{"magnet":"a","id":1},{"magnet":"b","error":{"code":"MAGNET_INVALID_URI","message":"m"}}]}"#)
        let result = try await client.uploadMagnets(["a", "b"])
        let last = StubProtocol.last
        #expect(last.method == "POST")
        #expect(last.url.path == "/v4/magnet/upload")
        #expect(form(last) == ["magnets[]": ["a", "b"]])
        #expect(last.headers["Authorization"] == "Bearer test-key")
        #expect(last.headers["User-Agent"] == "unit")
        #expect(result[0].id == 1)
        #expect(result[1].error?.code == "MAGNET_INVALID_URI")
    }

    @Test func noKeyNoAuthorization() async throws {
        StubProtocol.ok(#"{"ping":"pong"}"#)
        try await client.withAPIKey(nil).ping()
        #expect(StubProtocol.last.headers["Authorization"] == nil)
    }

    @Test func getParametersInQuery() async throws {
        StubProtocol.ok(#"{"hosts":{}}"#)
        _ = try await client.hosts(hostsOnly: true)
        #expect(StubProtocol.last.method == "GET")
        #expect(StubProtocol.last.url.query == "hostsOnly=true")
    }

    @Test func torrentFilesAreMultipart() async throws {
        StubProtocol.ok(#"{"files":[{"magnet":"m","id":9}]}"#)
        _ = try await client.uploadTorrentFiles([
            TorrentFile(name: "a.torrent", content: Data("d4:infoe".utf8)),
            TorrentFile(name: "b.torrent", content: Data("x".utf8)),
        ])
        let body = String(decoding: StubProtocol.last.body, as: UTF8.self)
        #expect(StubProtocol.last.headers["Content-Type"]?.hasPrefix("multipart/form-data") == true)
        #expect(body.components(separatedBy: #"name="files[]""#).count == 3)
        #expect(body.contains(#"filename="a.torrent""#) && body.contains(#"filename="b.torrent""#))
    }

    typealias Route = (String, @Sendable (AllDebrid) async throws -> Void, [String: [String]], String)

    static let routes: [Route] = [
        ("GET /v4/ping", { try await $0.ping() }, [:], "{}"),
        (
            "GET /v4.1/pin/get", { _ = try await $0.getPin() }, [:],
            #"{"pin":"P","check":"C","expires_in":600,"user_url":"u","base_url":"b"}"#
        ),
        (
            "POST /v4/pin/check", { _ = try await $0.checkPin(pin: "P", check: "C") }, ["pin": ["P"], "check": ["C"]],
            #"{"activated":false,"expires_in":600}"#
        ),
        ("GET /v4/hosts", { _ = try await $0.hosts() }, [:], #"{"hosts":{}}"#),
        ("GET /v4/hosts/domains", { _ = try await $0.hostDomains() }, [:], #"{"hosts":[]}"#),
        ("GET /v4/hosts/priority", { _ = try await $0.hostPriority() }, [:], #"{"hosts":{}}"#),
        ("GET /v4/user", { _ = try await $0.user() }, [:], userJSON),
        ("GET /v4.1/user/hosts", { _ = try await $0.userHosts() }, [:], #"{"hosts":{}}"#),
        ("POST /v4/user/verif", { _ = try await $0.verificationStatus(token: "T") }, ["token": ["T"]], #"{"verif":"waiting"}"#),
        ("POST /v4/user/verif/resend", { try await $0.resendVerification(token: "T") }, ["token": ["T"]], "{}"),
        ("POST /v4/user/notification/clear", { try await $0.clearNotification(code: "N") }, ["code": ["N"]], "{}"),
        (
            "POST /v4/link/infos", { _ = try await $0.linkInfos(["a", "b"], password: "p") }, ["link[]": ["a", "b"], "password": ["p"]],
            #"{"infos":[]}"#
        ),
        ("POST /v4/link/redirector", { _ = try await $0.extractRedirector("L") }, ["link": ["L"]], #"{"links":[]}"#),
        ("POST /v4/link/unlock", { _ = try await $0.unlock("L") }, ["link": ["L"]], "{}"),
        ("POST /v4/link/streaming", { _ = try await $0.streamingLink(id: "G", stream: "S") }, ["id": ["G"], "stream": ["S"]], "{}"),
        ("POST /v4/link/delayed", { _ = try await $0.delayedLink("42") }, ["id": ["42"]], #"{"status":1,"time_left":5}"#),
        ("POST /v4/magnet/upload", { _ = try await $0.uploadMagnets(["M"]) }, ["magnets[]": ["M"]], #"{"magnets":[]}"#),
        (
            "POST /v4/magnet/upload/file", { _ = try await $0.uploadTorrentFiles([TorrentFile(name: "a", content: Data([1]))]) }, [:],
            #"{"files":[]}"#
        ),
        ("POST /v4.1/magnet/status", { _ = try await $0.magnets() }, [:], #"{"magnets":[]}"#),
        ("POST /v4/magnet/files", { _ = try await $0.magnetFiles(ids: [1, 2]) }, ["id[]": ["1", "2"]], #"{"magnets":[]}"#),
        ("POST /v4/magnet/delete", { try await $0.deleteMagnet(id: 5) }, ["id": ["5"]], "{}"),
        ("POST /v4/magnet/restart", { _ = try await $0.restartMagnets(ids: [5, 6]) }, ["ids[]": ["5", "6"]], "{}"),
        ("GET /v4/user/links", { _ = try await $0.savedLinks() }, [:], #"{"links":[]}"#),
        ("POST /v4/user/links/save", { try await $0.saveLinks(["L"]) }, ["links[]": ["L"]], "{}"),
        ("POST /v4/user/links/delete", { try await $0.deleteSavedLinks(["L"]) }, ["links[]": ["L"]], "{}"),
        ("GET /v4/user/history", { _ = try await $0.history() }, [:], #"{"links":[]}"#),
        ("POST /v4/user/history/delete", { try await $0.purgeHistory() }, [:], "{}"),
        ("GET /v4/voucher/balance", { _ = try await $0.voucherBalance() }, [:], #"{"balance":1.5}"#),
        (
            "POST /v4/voucher/get", { _ = try await $0.vouchers(durationDays: 30, count: 2) }, ["duration": ["30"], "nb": ["2"]],
            #"{"codes":[]}"#
        ),
        (
            "POST /v4/voucher/generate", { _ = try await $0.generateVouchers(durationDays: 90, count: 1) },
            ["duration": ["90"], "nb": ["1"]], #"{"codes":[]}"#
        ),
    ]

    @Test func routesEveryOperationOfTheSpec() async throws {
        let ops = try specOperations()
        for (want, call, expected, data) in Self.routes {
            StubProtocol.ok(data)
            try await call(client)
            let got = "\(StubProtocol.last.method) \(StubProtocol.last.url.path)"
            #expect(got == want)
            #expect(ops.contains(got), "\(got) is not in the spec")
            if !expected.isEmpty {
                #expect(form(StubProtocol.last) == expected, "\(want)")
            }
        }
        #expect(Set(Self.routes.map(\.0)) == ops)
    }

    @Test func errorCodesMatchTheSpec() throws {
        let text = try String(contentsOf: specPath, encoding: .utf8)
        let block = text.components(separatedBy: "\n    ErrorCode:\n")[1]
        var codes: [String] = []
        for line in block.split(separator: "\n").map(String.init) {
            if line.hasPrefix("    "), !line.hasPrefix("     ") { break }
            if line.hasPrefix("        - ") {
                codes.append(String(line.dropFirst(10)).trimmingCharacters(in: CharacterSet(charactersIn: "\"")))
            }
        }
        #expect(codes.sorted() == ErrorCode.all.sorted())
    }

    @Test func apiErrorsCarryTheirCode() async throws {
        StubProtocol.answer(200, #"{"status":"error","error":{"code":"AUTH_BLOCKED","message":"m","token":"tok"}}"#)
        await #expect(throws: AllDebridError.api(code: ErrorCode.authBlocked, message: "m", token: "tok")) {
            _ = try await client.user()
        }
    }

    @Test func httpErrors() async throws {
        StubProtocol.answer(429, "Too Many Requests")
        StubProtocol.answer(404, #"{"status":"error","error":{"code":"404","message":"Endpoint doesn't exist"}}"#)
        StubProtocol.answer(200, "<html>")
        do {
            try await client.ping()
            Issue.record("no error")
        } catch let e as AllDebridError {
            #expect(e.isRateLimited)
        }
        await #expect(throws: AllDebridError.http(status: 404, code: ErrorCode.notFound)) { try await client.ping() }
        await #expect(throws: AllDebridError.self) { try await client.ping() }
    }

    @Test func flexibleShapes() async throws {
        StubProtocol.ok(#"{"magnets":{"id":7,"statusCode":4,"files":[{"n":"a","s":1,"l":"L"}]}}"#)
        let magnet = try await client.magnet(id: 7)
        #expect(magnet.id == 7 && magnet.isReady)
        #expect(magnet.files?.first?.l == "L")

        StubProtocol.ok(userJSON)
        #expect(try await client.user().limitedHostersQuotas.isEmpty)

        StubProtocol.ok(#"{"link":"","delayed":123,"streaming":[{"id":"mp3","quality":"mp3"},{"id":"hd","quality":1080}]}"#)
        let unlocked = try await client.unlock("L")
        #expect(unlocked.delayed?.value == "123")
        #expect(unlocked.streaming?.map { $0.quality?.value } == ["mp3", "1080"])

        StubProtocol.ok(#"{"hosts":{"a":{"name":"a","type":"free","domains":[],"regexp":"x"}}}"#)
        #expect(try await client.hosts().hosts["a"]?.regexp == ["x"])
    }

    @Test func missingRequiredFieldsAreRejected() async throws {
        StubProtocol.ok(#"{"user":{"username":"u"}}"#)
        await #expect(throws: AllDebridError.self) { _ = try await client.user() }
    }

    @Test func waitForPin() async throws {
        StubProtocol.ok(#"{"activated":false,"expires_in":600}"#)
        StubProtocol.ok(#"{"activated":true,"expires_in":590,"apikey":"K"}"#)
        let pin = PinGetData(pin: "ABCD", check: "chk", expiresIn: 600, userUrl: "u", baseUrl: "b")
        #expect(try await client.waitForPin(pin, interval: .milliseconds(1)) == "K")
        #expect(form(StubProtocol.last) == ["pin": ["ABCD"], "check": ["chk"]])
    }

    @Test func waitDelayedLink() async throws {
        StubProtocol.ok(#"{"status":1,"time_left":5}"#)
        StubProtocol.ok(#"{"status":2,"time_left":0,"link":"L"}"#)
        StubProtocol.ok(#"{"status":3,"time_left":0}"#)
        #expect(try await client.waitDelayedLink("1", interval: .milliseconds(1)) == "L")
        await #expect(throws: AllDebridError.self) { _ = try await client.waitDelayedLink("1", interval: .milliseconds(1)) }
    }

    @Test func liveMagnets() async throws {
        StubProtocol.ok(
            #"{"fullsync":true,"counter":1,"magnets":[{"id":1,"filename":"a","uploadDate":10,"downloaded":5},{"id":2,"filename":"b","uploadDate":20}]}"#
        )
        StubProtocol.ok(#"{"counter":2,"magnets":{"id":1,"downloaded":9}}"#)
        StubProtocol.ok(#"{"counter":3,"magnets":[{"id":2,"deleted":true},{"id":3,"filename":"c","uploadDate":30}]}"#)
        StubProtocol.ok(#"{"counter":4,"magnets":[]}"#)
        StubProtocol.ok(#"{"fullsync":true,"counter":1,"magnets":[{"id":9,"uploadDate":1}]}"#)
        let live = client.liveMagnets()

        #expect(try await live.poll().map(\.id) == [2, 1])
        let first = form(StubProtocol.last)
        #expect(first["counter"] == ["0"])

        let merged = try await live.poll()
        let second = form(StubProtocol.last)
        #expect(second["counter"] == ["1"])
        #expect(first["session"] == second["session"])
        let one = merged.first { $0.id == 1 }
        #expect(one?.filename == "a" && one?.downloaded == 9)

        #expect(try await live.poll().map(\.id) == [3, 1])
        #expect(try await live.poll().map(\.id) == [3, 1])
        #expect(try await live.poll().map(\.id) == [9])
    }
}
