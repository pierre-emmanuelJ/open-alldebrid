package io.github.pierreemmanuelj.alldebrid

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue
import kotlinx.coroutines.test.runTest
import org.junit.jupiter.api.Assumptions.assumeFalse
import org.junit.jupiter.api.BeforeEach

/** Calls the real API with the public demo keys. Set OFFLINE=1 to skip. */
class DemoTest {
    private val client = AllDebrid(AllDebrid.DEMO_PREMIUM_KEY, userAgent = "open-alldebrid-kotlin-tests")
    private val anonymous = AllDebrid(userAgent = "open-alldebrid-kotlin-tests")

    @BeforeEach
    fun online() {
        assumeFalse(System.getenv("OFFLINE") == "1", "OFFLINE=1")
    }

    @Test
    fun publicEndpoints() = runTest {
        anonymous.ping()
        assertTrue(anonymous.hosts().hosts.isNotEmpty())
        assertTrue(anonymous.hostDomains().hosts.isNotEmpty())
        assertTrue(anonymous.hostPriority().isNotEmpty())
    }

    @Test
    fun pin() = runTest {
        assertTrue(anonymous.getPin().userUrl.startsWith("https://"))
        assertTrue(anonymous.checkPin("123B", "664c3ca2635c99f291d28e11ea18e154750bd21a").activated)
    }

    @Test
    fun user() = runTest {
        assertEquals("demoUserPremium", client.user().username)
        assertTrue(client.userHosts().isNotEmpty())
        client.clearNotification("NOTIF_CODE")
        assertEquals("allowed", anonymous.verificationStatus("verificationTokenAllowed3f9eb0db").verif)
        val err = assertFailsWith<AllDebridException> { anonymous.resendVerification("verificationTokenAlreadySentb0db") }
        assertEquals(ErrorCodes.ALREADY_SENT, err.code)
    }

    @Test
    fun links() = runTest {
        val infos = client.linkInfos(listOf("https://example.com/file2MB.txt", "https://example.com/down"))
        assertTrue(!infos[0].filename.isNullOrEmpty())
        assertTrue(client.extractRedirector("https://example.com/redir/folder").isNotEmpty())
        assertTrue(client.unlock("https://example.com/streaming").streaming!!.isNotEmpty())
        assertTrue(client.unlock("https://example.com/delayed").delayed != null)
        assertTrue(assertFailsWith<AllDebridException> { client.unlock("https://example.com/down") }.isErrorCode(ErrorCodes.LINK_DOWN))
        assertTrue(!client.streamingLink("demov1bc797", "hd-720").link.isNullOrEmpty())
        assertEquals(DelayedStatus.READY, client.delayedLink("1234567").status)
    }

    @Test
    fun magnets() = runTest {
        assertTrue(client.uploadMagnets(listOf("magnet:?xt=urn:btih:3A842783e3005495d5d1637f5364b59343c7844707"))[0].id != null)
        client.uploadTorrentFiles(listOf(TorrentFile("demo.torrent", "d4:infod4:name4:demoee".toByteArray())))
        assertTrue(client.magnets().isNotEmpty())
        assertTrue(client.magnet(56789).files!!.isNotEmpty())
        assertEquals(1, client.magnetFiles(listOf(56789)).size)
        client.deleteMagnet(56789)
        client.restartMagnet(99999)
        assertTrue(client.liveMagnets().poll().isNotEmpty())
    }

    @Test
    fun savedLinks() = runTest {
        assertTrue(client.savedLinks().isNotEmpty())
        client.saveLinks(listOf("https://example.com/file2MB.txt"))
        client.deleteSavedLinks(listOf("https://example.com/file2MB.txt"))
        client.history()
        client.purgeHistory()
    }

    @Test
    fun badKey() = runTest {
        val err = assertFailsWith<AllDebridException> { client.withApiKey("not-a-real-key").user() }
        assertEquals(ErrorCodes.AUTH_BAD_APIKEY, err.code)
    }
}
