import Foundation
import Testing
@testable import OpenAllDebrid

/// Calls the real API with the public demo keys. Set OFFLINE=1 to skip.
@Suite(.enabled(if: ProcessInfo.processInfo.environment["OFFLINE"] != "1"))
struct DemoTests {
    let client = AllDebrid(apiKey: AllDebrid.DemoKey.premium, userAgent: "open-alldebrid-swift-tests")
    let anonymous = AllDebrid(userAgent: "open-alldebrid-swift-tests")

    @Test func publicEndpoints() async throws {
        try await anonymous.ping()
        #expect(try await !anonymous.hosts().hosts.isEmpty)
        #expect(try await !anonymous.hostDomains().hosts.isEmpty)
        #expect(try await !anonymous.hostPriority().isEmpty)
    }

    @Test func pin() async throws {
        #expect(try await anonymous.getPin().userUrl.hasPrefix("https://"))
        #expect(try await anonymous.checkPin(pin: "123B", check: "664c3ca2635c99f291d28e11ea18e154750bd21a").activated)
    }

    @Test func user() async throws {
        #expect(try await client.user().username == "demoUserPremium")
        #expect(try await !client.userHosts().isEmpty)
        try await client.clearNotification(code: "NOTIF_CODE")
        #expect(try await anonymous.verificationStatus(token: "verificationTokenAllowed3f9eb0db").verif == "allowed")
        do {
            try await anonymous.resendVerification(token: "verificationTokenAlreadySentb0db")
            Issue.record("no error")
        } catch let e as AllDebridError {
            #expect(e.code == ErrorCode.alreadySent)
        }
    }

    @Test func links() async throws {
        let infos = try await client.linkInfos(["https://example.com/file2MB.txt", "https://example.com/down"])
        #expect(infos.first?.filename?.isEmpty == false)
        #expect(try await !client.extractRedirector("https://example.com/redir/folder").isEmpty)
        #expect(try await client.unlock("https://example.com/streaming").streaming?.isEmpty == false)
        #expect(try await client.unlock("https://example.com/delayed").delayed != nil)
        do {
            _ = try await client.unlock("https://example.com/down")
            Issue.record("no error")
        } catch let e as AllDebridError {
            #expect(e.code == ErrorCode.linkDown)
        }
        #expect(try await client.streamingLink(id: "demov1bc797", stream: "hd-720").link?.isEmpty == false)
        #expect(try await client.delayedLink("1234567").status == DelayedStatus.ready)
    }

    @Test func magnets() async throws {
        #expect(try await client.uploadMagnets(["magnet:?xt=urn:btih:3A842783e3005495d5d1637f5364b59343c7844707"]).first?.id != nil)
        _ = try await client.uploadTorrentFiles([TorrentFile(name: "demo.torrent", content: Data("d4:infod4:name4:demoee".utf8))])
        #expect(try await !client.magnets().isEmpty)
        #expect(try await client.magnet(id: 56789).files?.isEmpty == false)
        #expect(try await client.magnetFiles(ids: [56789]).count == 1)
        try await client.deleteMagnet(id: 56789)
        try await client.restartMagnet(id: 99999)
        #expect(try await !client.liveMagnets().poll().isEmpty)
    }

    @Test func savedLinks() async throws {
        #expect(try await !client.savedLinks().isEmpty)
        try await client.saveLinks(["https://example.com/file2MB.txt"])
        try await client.deleteSavedLinks(["https://example.com/file2MB.txt"])
        _ = try await client.history()
        try await client.purgeHistory()
    }

    @Test func badKey() async throws {
        do {
            _ = try await client.withAPIKey("not-a-real-key").user()
            Issue.record("no error")
        } catch let e as AllDebridError {
            #expect(e.code == ErrorCode.authBadApikey)
        }
    }
}
