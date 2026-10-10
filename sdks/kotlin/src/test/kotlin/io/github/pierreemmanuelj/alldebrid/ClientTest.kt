package io.github.pierreemmanuelj.alldebrid

import java.io.File
import java.net.URLDecoder
import kotlin.test.AfterTest
import kotlin.test.BeforeTest
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNull
import kotlin.test.assertTrue
import kotlin.time.Duration.Companion.milliseconds
import kotlinx.coroutines.test.runTest
import mockwebserver3.MockResponse
import mockwebserver3.MockWebServer
import mockwebserver3.RecordedRequest
import org.yaml.snakeyaml.Yaml

class ClientTest {
    private lateinit var server: MockWebServer
    private lateinit var client: AllDebrid

    @BeforeTest
    fun start() {
        server = MockWebServer()
        server.start()
        client = AllDebrid("test-key", baseUrl = server.url("/").toString(), userAgent = "unit")
    }

    @AfterTest
    fun stop() {
        server.close()
    }

    private fun answer(body: String, code: Int = 200) {
        server.enqueue(MockResponse.Builder().code(code).body(body).build())
    }

    private fun ok(data: String) = answer("""{"status":"success","data":$data}""")

    private fun apiError(code: String, extra: String = "") =
        answer("""{"status":"error","error":{"code":"$code","message":"m"$extra}}""")

    private fun RecordedRequest.form(): Map<String, List<String>> {
        val text = body?.utf8().orEmpty()
        if (text.isEmpty()) return emptyMap()
        return text.split("&").map { it.split("=", limit = 2) }
            .groupBy({ URLDecoder.decode(it[0], "UTF-8") }, { URLDecoder.decode(it.getOrElse(1) { "" }, "UTF-8") })
    }

    @Test
    fun listParametersArePhpStyle() = runTest {
        ok("""{"magnets":[{"magnet":"a","id":1}]}""")
        val result = client.uploadMagnets(listOf("a", "b"))
        val req = server.takeRequest()
        assertEquals("POST", req.method)
        assertEquals("/v4/magnet/upload", req.url.encodedPath)
        assertEquals(listOf("a", "b"), req.form()["magnets[]"])
        assertEquals("Bearer test-key", req.headers["Authorization"])
        assertEquals("unit", req.headers["User-Agent"])
        assertEquals(1L, result[0].id)
    }

    @Test
    fun noKeyNoAuthorization() = runTest {
        ok("""{"ping":"pong"}""")
        client.withApiKey(null).ping()
        assertNull(server.takeRequest().headers["Authorization"])
    }

    @Test
    fun getParametersInQuery() = runTest {
        ok("""{"hosts":{}}""")
        client.hosts(hostsOnly = true)
        val req = server.takeRequest()
        assertEquals("GET", req.method)
        assertEquals("true", req.url.queryParameter("hostsOnly"))
    }

    @Test
    fun torrentFilesAreMultipart() = runTest {
        ok("""{"files":[{"magnet":"m","id":9}]}""")
        client.uploadTorrentFiles(listOf(TorrentFile("a.torrent", byteArrayOf(1)), TorrentFile("b.torrent", byteArrayOf(2))))
        val req = server.takeRequest()
        val body = req.body?.utf8().orEmpty()
        assertTrue(req.headers["Content-Type"]!!.startsWith("multipart/form-data"))
        assertEquals(2, Regex("name=\"files\\[]\"").findAll(body).count())
        assertTrue("filename=\"a.torrent\"" in body && "filename=\"b.torrent\"" in body)
    }

    @Test
    fun routesEveryOperationOfTheSpec() = runTest {
        @Suppress("UNCHECKED_CAST")
        val spec = Yaml().load<Map<String, Any>>(File(System.getenv("OPENAPI_SPEC")).readText())
        @Suppress("UNCHECKED_CAST")
        val operations = (spec["paths"] as Map<String, Map<String, Any>>)
            .flatMap { (path, item) -> item.keys.map { "${it.uppercase()} $path" } }.toSet()

        val calls: List<Triple<String, suspend () -> Unit, Map<String, List<String>>>> = listOf(
            Triple("GET /v4/ping", { client.ping() }, emptyMap()),
            Triple("GET /v4.1/pin/get", { client.getPin() }, emptyMap()),
            Triple("POST /v4/pin/check", { client.checkPin("P", "C") }, mapOf("pin" to listOf("P"), "check" to listOf("C"))),
            Triple("GET /v4/hosts", { client.hosts() }, emptyMap()),
            Triple("GET /v4/hosts/domains", { client.hostDomains() }, emptyMap()),
            Triple("GET /v4/hosts/priority", { client.hostPriority() }, emptyMap()),
            Triple("GET /v4/user", { client.user() }, emptyMap()),
            Triple("GET /v4.1/user/hosts", { client.userHosts() }, emptyMap()),
            Triple("POST /v4/user/verif", { client.verificationStatus("T") }, mapOf("token" to listOf("T"))),
            Triple("POST /v4/user/verif/resend", { client.resendVerification("T") }, mapOf("token" to listOf("T"))),
            Triple("POST /v4/user/notification/clear", { client.clearNotification("N") }, mapOf("code" to listOf("N"))),
            Triple("POST /v4/link/infos", { client.linkInfos(listOf("a", "b"), "p") }, mapOf("link[]" to listOf("a", "b"), "password" to listOf("p"))),
            Triple("POST /v4/link/redirector", { client.extractRedirector("L") }, mapOf("link" to listOf("L"))),
            Triple("POST /v4/link/unlock", { client.unlock("L") }, mapOf("link" to listOf("L"))),
            Triple("POST /v4/link/streaming", { client.streamingLink("G", "S") }, mapOf("id" to listOf("G"), "stream" to listOf("S"))),
            Triple("POST /v4/link/delayed", { client.delayedLink("42") }, mapOf("id" to listOf("42"))),
            Triple("POST /v4/magnet/upload", { client.uploadMagnets(listOf("M")) }, mapOf("magnets[]" to listOf("M"))),
            Triple("POST /v4/magnet/upload/file", { client.uploadTorrentFiles(listOf(TorrentFile("a", byteArrayOf(1)))) }, emptyMap()),
            Triple("POST /v4.1/magnet/status", { client.magnets() }, emptyMap()),
            Triple("POST /v4/magnet/files", { client.magnetFiles(listOf(1, 2)) }, mapOf("id[]" to listOf("1", "2"))),
            Triple("POST /v4/magnet/delete", { client.deleteMagnet(5) }, mapOf("id" to listOf("5"))),
            Triple("POST /v4/magnet/restart", { client.restartMagnets(listOf(5, 6)) }, mapOf("ids[]" to listOf("5", "6"))),
            Triple("GET /v4/user/links", { client.savedLinks() }, emptyMap()),
            Triple("POST /v4/user/links/save", { client.saveLinks(listOf("L")) }, mapOf("links[]" to listOf("L"))),
            Triple("POST /v4/user/links/delete", { client.deleteSavedLinks(listOf("L")) }, mapOf("links[]" to listOf("L"))),
            Triple("GET /v4/user/history", { client.history() }, emptyMap()),
            Triple("POST /v4/user/history/delete", { client.purgeHistory() }, emptyMap()),
            Triple("GET /v4/voucher/balance", { client.voucherBalance() }, emptyMap()),
            Triple("POST /v4/voucher/get", { client.vouchers(30, 2) }, mapOf("duration" to listOf("30"), "nb" to listOf("2"))),
            Triple("POST /v4/voucher/generate", { client.generateVouchers(90, 1) }, mapOf("duration" to listOf("90"), "nb" to listOf("1"))),
        )
        val data = mapOf(
            "GET /v4.1/pin/get" to """{"pin":"P","check":"C","expires_in":600,"user_url":"u","base_url":"b"}""",
            "POST /v4/pin/check" to """{"activated":false,"expires_in":600}""",
            "GET /v4/hosts" to """{"hosts":{}}""",
            "GET /v4/hosts/domains" to """{"hosts":[]}""",
            "GET /v4/hosts/priority" to """{"hosts":{}}""",
            "GET /v4/user" to USER,
            "GET /v4.1/user/hosts" to """{"hosts":{}}""",
            "POST /v4/user/verif" to """{"verif":"waiting"}""",
            "POST /v4/link/infos" to """{"infos":[]}""",
            "POST /v4/link/redirector" to """{"links":[]}""",
            "POST /v4/link/delayed" to """{"status":1,"time_left":5}""",
            "POST /v4/magnet/upload" to """{"magnets":[]}""",
            "POST /v4/magnet/upload/file" to """{"files":[]}""",
            "POST /v4.1/magnet/status" to """{"magnets":[]}""",
            "POST /v4/magnet/files" to """{"magnets":[]}""",
            "GET /v4/user/links" to """{"links":[]}""",
            "GET /v4/user/history" to """{"links":[]}""",
            "GET /v4/voucher/balance" to """{"balance":1.5}""",
            "POST /v4/voucher/get" to """{"codes":[]}""",
            "POST /v4/voucher/generate" to """{"codes":[]}""",
        )
        for ((want, call, form) in calls) {
            ok(data[want] ?: """{"message":"ok"}""")
            call()
            val req = server.takeRequest()
            val got = "${req.method} ${req.url.encodedPath}"
            assertEquals(want, got)
            assertTrue(got in operations, "$got is not in the spec")
            if (form.isNotEmpty()) assertEquals(form, req.form())
        }
        assertEquals(operations, calls.map { it.first }.toSet())
    }

    @Test
    fun errorCodesMatchTheSpec() {
        @Suppress("UNCHECKED_CAST")
        val spec = Yaml().load<Map<String, Any>>(File(System.getenv("OPENAPI_SPEC")).readText())
        @Suppress("UNCHECKED_CAST")
        val schemas = (spec["components"] as Map<String, Any>)["schemas"] as Map<String, Map<String, Any>>
        val enum = (schemas.getValue("ErrorCode")["enum"] as List<*>).map { it.toString() }
        assertEquals(enum.sorted(), ErrorCodes.ALL.sorted())
    }

    @Test
    fun apiErrorsCarryTheirCode() = runTest {
        apiError("MAGNET_INVALID_ID")
        val err = assertFailsWith<AllDebridException> { client.deleteMagnet(1) }
        assertEquals(ErrorCodes.MAGNET_INVALID_ID, err.code)
        assertTrue(err.isErrorCode("MAGNET_INVALID_ID"))
    }

    @Test
    fun blockedLoginKeepsToken() = runTest {
        apiError("AUTH_BLOCKED", ""","token":"tok"""")
        val err = assertFailsWith<AllDebridException> { client.user() }
        assertEquals("tok", err.token)
    }

    @Test
    fun httpErrors() = runTest {
        answer("Too Many Requests", 429)
        assertTrue(assertFailsWith<HttpException> { client.ping() }.isRateLimited)
        answer("""{"status":"error","error":{"code":"404","message":"Endpoint doesn't exist"}}""", 404)
        assertTrue(assertFailsWith<HttpException> { client.ping() }.isErrorCode(ErrorCodes.NOT_FOUND))
        answer("<html>")
        assertFailsWith<java.io.IOException> { client.ping() }
    }

    @Test
    fun flexibleShapes() = runTest {
        ok("""{"magnets":{"id":7,"statusCode":4,"files":[{"n":"a","s":1,"l":"L"}]}}""")
        val m = client.magnet(7)
        assertEquals(7L, m.id)
        assertTrue(m.isReady)
        assertEquals("L", m.files!![0].l)

        ok(USER)
        assertEquals(emptyMap(), client.user().limitedHostersQuotas)

        ok("""{"link":"","delayed":"123","streaming":[{"id":"mp3","quality":"mp3"},{"id":"hd","quality":1080}]}""")
        val u = client.unlock("L")
        assertEquals("123", u.delayed?.value)
        assertEquals(listOf("mp3", "1080"), u.streaming!!.map { it.quality?.value })

        ok("""{"hosts":{"a":{"name":"a","type":"free","domains":[],"regexp":["x","y"]}}}""")
        assertEquals(listOf("x", "y"), client.hosts().hosts["a"]!!.regexp)
    }

    @Test
    fun waitForPin() = runTest {
        ok("""{"activated":false,"expires_in":600}""")
        ok("""{"activated":true,"expires_in":590,"apikey":"K"}""")
        val pin = PinGetData("ABCD", "chk", 600, "u", "b")
        assertEquals("K", client.waitForPin(pin, 1.milliseconds))
        apiError("PIN_EXPIRED")
        assertFailsWith<AllDebridException> { client.waitForPin(pin, 1.milliseconds) }
    }

    @Test
    fun waitDelayedLink() = runTest {
        ok("""{"status":1,"time_left":5}""")
        ok("""{"status":2,"time_left":0,"link":"L"}""")
        assertEquals("L", client.waitDelayedLink("1", 1.milliseconds))
        ok("""{"status":3,"time_left":0}""")
        assertFailsWith<java.io.IOException> { client.waitDelayedLink("1", 1.milliseconds) }
    }

    @Test
    fun liveMagnets() = runTest {
        ok("""{"fullsync":true,"counter":1,"magnets":[{"id":1,"filename":"a","uploadDate":10,"downloaded":5},{"id":2,"filename":"b","uploadDate":20}]}""")
        ok("""{"counter":2,"magnets":{"id":1,"downloaded":9}}""")
        ok("""{"counter":3,"magnets":[{"id":2,"deleted":true},{"id":3,"filename":"c","uploadDate":30}]}""")
        ok("""{"counter":4,"magnets":[]}""")
        ok("""{"fullsync":true,"counter":1,"magnets":[{"id":9,"uploadDate":1}]}""")
        val live = client.liveMagnets()

        assertEquals(listOf(2L, 1L), live.poll().map { it.id })
        val first = server.takeRequest().form()
        assertEquals(listOf("0"), first["counter"])

        val merged = live.poll()
        val second = server.takeRequest().form()
        assertEquals(listOf("1"), second["counter"])
        assertEquals(first["session"], second["session"])
        val one = merged.first { it.id == 1L }
        assertEquals("a", one.filename)
        assertEquals(9L, one.downloaded)

        assertEquals(listOf(3L, 1L), live.poll().map { it.id })
        assertEquals(listOf(3L, 1L), live.poll().map { it.id })
        assertEquals(listOf(9L), live.poll().map { it.id })
    }

    private companion object {
        const val USER = """{"user":{"username":"u","email":"e","isPremium":true,"isSubscribed":false,"isTrial":false,"premiumUntil":0,"lang":"en","preferedDomain":"com","fidelityPoints":0,"limitedHostersQuotas":[],"notifications":[]}}"""
    }
}
