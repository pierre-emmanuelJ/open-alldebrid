package io.github.pierreemmanuelj.alldebrid

import java.io.IOException
import kotlin.coroutines.resume
import kotlin.coroutines.resumeWithException
import kotlin.random.Random
import kotlin.time.Duration
import kotlin.time.Duration.Companion.seconds
import kotlinx.coroutines.delay
import kotlinx.coroutines.suspendCancellableCoroutine
import kotlinx.coroutines.sync.Mutex
import kotlinx.coroutines.sync.withLock
import kotlinx.serialization.DeserializationStrategy
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonArray
import kotlinx.serialization.json.JsonElement
import kotlinx.serialization.json.JsonNull
import kotlinx.serialization.json.JsonObject
import kotlinx.serialization.json.booleanOrNull
import kotlinx.serialization.json.contentOrNull
import kotlinx.serialization.json.jsonObject
import kotlinx.serialization.json.jsonPrimitive
import kotlinx.serialization.json.longOrNull
import okhttp3.Call
import okhttp3.Callback
import okhttp3.FormBody
import okhttp3.HttpUrl.Companion.toHttpUrl
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.MultipartBody
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody
import okhttp3.RequestBody.Companion.toRequestBody
import okhttp3.Response

/** A .torrent file to upload. */
public class TorrentFile(public val name: String, public val content: ByteArray)

/**
 * A client for the AllDebrid API. Every call is a suspend function; API
 * errors throw [AllDebridException], other HTTP statuses [HttpException].
 *
 * ```
 * val client = AllDebrid(apiKey, userAgent = "my-app")
 * val user = client.user()
 * ```
 */
public class AllDebrid(
    private val apiKey: String? = null,
    baseUrl: String = DEFAULT_BASE_URL,
    private val userAgent: String = DEFAULT_USER_AGENT,
    private val http: OkHttpClient = defaultHttpClient,
) {
    private val baseUrl = baseUrl.trimEnd('/')

    public companion object {
        public const val DEFAULT_BASE_URL: String = "https://api.alldebrid.com"
        public const val DEFAULT_USER_AGENT: String = "open-alldebrid-kotlin"

        /** Public demo keys: the API answers them with static sample data. */
        public const val DEMO_PREMIUM_KEY: String = "staticDemoApikeyPrem"
        public const val DEMO_TRIAL_KEY: String = "staticDemoApikeyTria"
        public const val DEMO_FREE_KEY: String = "staticDemoApikeyFree"

        private val defaultHttpClient: OkHttpClient by lazy { OkHttpClient() }

        internal val json: Json = Json {
            ignoreUnknownKeys = true
            explicitNulls = false
        }
    }

    /** A client with another API key, sharing this one's HTTP client. */
    public fun withApiKey(apiKey: String?): AllDebrid = AllDebrid(apiKey, baseUrl, userAgent, http)

    private suspend fun <T> call(
        method: String,
        path: String,
        params: List<Pair<String, String>> = emptyList(),
        serializer: DeserializationStrategy<T>,
        body: RequestBody? = null,
    ): T {
        val data = callRaw(method, path, params, body)
        return json.decodeFromJsonElement(serializer, data)
    }

    internal suspend fun callRaw(
        method: String,
        path: String,
        params: List<Pair<String, String>> = emptyList(),
        body: RequestBody? = null,
    ): JsonElement {
        val url = (baseUrl + path).toHttpUrl().newBuilder()
        var requestBody = body
        if (method == "GET") {
            params.forEach { (k, v) -> url.addQueryParameter(k, v) }
        } else if (requestBody == null) {
            // List parameters keep their PHP-style names and repeat: magnets[]=a&magnets[]=b.
            val form = FormBody.Builder()
            params.forEach { (k, v) -> form.add(k, v) }
            requestBody = if (params.isEmpty()) ByteArray(0).toRequestBody(FORM) else form.build()
        }
        val request = Request.Builder()
            .url(url.build())
            .header("Accept", "application/json")
            .header("User-Agent", userAgent)
            .apply { if (!apiKey.isNullOrEmpty()) header("Authorization", "Bearer $apiKey") }
            .method(method, requestBody)
            .build()

        val (status, text) = http.newCall(request).await().use { it.code to it.body.string() }
        val envelope = runCatching { json.parseToJsonElement(text).jsonObject }.getOrNull()
        val apiError = envelope?.get("error")?.let(::toException)
        if (status != 200) throw HttpException(status, apiError)
        if (envelope == null) throw IOException("AllDebrid $method $path: the response is not a JSON object")
        if (envelope["status"]?.jsonPrimitive?.contentOrNull != "success") {
            throw apiError ?: IOException("AllDebrid $method $path: error without details")
        }
        return envelope["data"] ?: JsonNull
    }

    private fun toException(element: JsonElement): AllDebridException? {
        val error = element as? JsonObject ?: return null
        return AllDebridException(
            code = error["code"]?.jsonPrimitive?.contentOrNull ?: "",
            message = error["message"]?.jsonPrimitive?.contentOrNull ?: "",
            token = error["token"]?.jsonPrimitive?.contentOrNull,
        )
    }

    // ---------------------------------------------------------------- general

    /** Checks that the API is up. */
    public suspend fun ping() {
        callRaw("GET", "/v4/ping")
    }

    // -------------------------------------------------------------------- pin

    /** Starts a PIN authentication. Show [PinGetData.pin] or [PinGetData.userUrl], then [waitForPin]. */
    public suspend fun getPin(): PinGetData = call("GET", "/v4.1/pin/get", serializer = PinGetData.serializer())

    /** The state of a PIN authentication; [PinCheckData.apikey] is set once activated. */
    public suspend fun checkPin(pin: String, check: String): PinCheckData =
        call("POST", "/v4/pin/check", listOf("pin" to pin, "check" to check), PinCheckData.serializer())

    /** Polls until the user confirms the PIN and returns the new API key. */
    public suspend fun waitForPin(pin: PinGetData, interval: Duration = 5.seconds): String {
        while (true) {
            val state = checkPin(pin.pin, pin.check)
            val key = state.apikey
            if (state.activated && !key.isNullOrEmpty()) return key
            delay(interval)
        }
    }

    // ------------------------------------------------------------------ hosts

    /** Supported hosts, streaming sites and redirectors. Cache the result. */
    public suspend fun hosts(hostsOnly: Boolean = false): HostsData =
        call("GET", "/v4/hosts", hostsOnlyParam(hostsOnly), HostsData.serializer())

    /** Every supported domain. */
    public suspend fun hostDomains(): HostDomainsData =
        call("GET", "/v4/hosts/domains", serializer = HostDomainsData.serializer())

    /** Host names ranked from the most open (1) to the most restricted. */
    public suspend fun hostPriority(): Map<String, Int> =
        call("GET", "/v4/hosts/priority", serializer = HostPriorityData.serializer()).hosts

    // ------------------------------------------------------------------- user

    /** The account information. */
    public suspend fun user(): User = call("GET", "/v4/user", serializer = UserData.serializer()).user

    /** The hosts this user can use now, with their quotas. */
    public suspend fun userHosts(hostsOnly: Boolean = false): ServiceMap =
        call("GET", "/v4.1/user/hosts", hostsOnlyParam(hostsOnly), UserHostsData.serializer()).hosts

    /** Follows a new-location login verification (token from an `AUTH_BLOCKED` error). */
    public suspend fun verificationStatus(token: String): VerificationData =
        call("POST", "/v4/user/verif", listOf("token" to token), VerificationData.serializer())

    /** Sends the verification email again (once per token). */
    public suspend fun resendVerification(token: String) {
        callRaw("POST", "/v4/user/verif/resend", listOf("token" to token))
    }

    /** Clears one of the codes in [User.notifications]. */
    public suspend fun clearNotification(code: String) {
        callRaw("POST", "/v4/user/notification/clear", listOf("code" to code))
    }

    // ------------------------------------------------------------------ links

    /** File name and size of each link, in order; unreadable links have [LinkInfo.error] set. */
    public suspend fun linkInfos(links: List<String>, password: String? = null): List<LinkInfo> {
        val params = links.map { "link[]" to it } + listOfNotNull(password?.let { "password" to it })
        return call("POST", "/v4/link/infos", params, LinkInfosData.serializer()).infos
    }

    /** Links behind a redirector; they only work with [linkInfos] and [unlock]. */
    public suspend fun extractRedirector(link: String): List<String> =
        call("POST", "/v4/link/redirector", listOf("link" to link), RedirectorData.serializer()).links

    /**
     * Unlocks a hoster link. When [UnlockData.delayed] is set, wait with [waitDelayedLink];
     * when [UnlockData.streaming] lists qualities, choose one with [streamingLink].
     */
    public suspend fun unlock(link: String, password: String? = null): UnlockData {
        val params = listOf("link" to link) + listOfNotNull(password?.let { "password" to it })
        return call("POST", "/v4/link/unlock", params, UnlockData.serializer())
    }

    /** The link of one quality: [id] from [unlock], [stream] from one of its streaming entries. */
    public suspend fun streamingLink(id: String, stream: String): StreamingData =
        call("POST", "/v4/link/streaming", listOf("id" to id, "stream" to stream), StreamingData.serializer())

    /** The status of a delayed link. */
    public suspend fun delayedLink(id: String): DelayedData =
        call("POST", "/v4/link/delayed", listOf("id" to id), DelayedData.serializer())

    /** Polls a delayed link until it is ready and returns the download link. */
    public suspend fun waitDelayedLink(id: String, interval: Duration = 5.seconds): String {
        while (true) {
            val state = delayedLink(id)
            val link = state.link
            when {
                state.status == DelayedStatus.READY && !link.isNullOrEmpty() -> return link
                state.status == DelayedStatus.FAILED -> throw IOException("AllDebrid: delayed link $id failed")
            }
            delay(interval)
        }
    }

    // ---------------------------------------------------------------- magnets

    /** Adds magnet URIs or info hashes. Results are in order; invalid ones have an error. */
    public suspend fun uploadMagnets(magnets: List<String>): List<UploadedMagnet> =
        call("POST", "/v4/magnet/upload", magnets.map { "magnets[]" to it }, UploadMagnetsData.serializer()).magnets

    /** Adds .torrent files. Results are in order; invalid files have an error. */
    public suspend fun uploadTorrentFiles(files: List<TorrentFile>): List<UploadedMagnet> {
        val body = MultipartBody.Builder().setType(MultipartBody.FORM).apply {
            files.forEach { addFormDataPart("files[]", it.name, it.content.toRequestBody(TORRENT)) }
        }.build()
        return call("POST", "/v4/magnet/upload/file", serializer = UploadTorrentFilesData.serializer(), body = body).files
    }

    /** Every magnet of the account. */
    public suspend fun magnets(): List<Magnet> =
        call("POST", "/v4.1/magnet/status", serializer = MagnetStatusData.serializer()).magnets

    /** One magnet, with its files. */
    public suspend fun magnet(id: Long): Magnet {
        val found = call("POST", "/v4.1/magnet/status", listOf("id" to id.toString()), MagnetStatusData.serializer()).magnets
        return found.singleOrNull() ?: throw IOException("AllDebrid: magnet $id: got ${found.size} magnets")
    }

    /** The file tree of each magnet, in order; unknown IDs have an error. */
    public suspend fun magnetFiles(ids: List<Long>): List<MagnetFiles> =
        call("POST", "/v4/magnet/files", ids.map { "id[]" to it.toString() }, MagnetFilesData.serializer()).magnets

    /** Deletes a magnet. */
    public suspend fun deleteMagnet(id: Long) {
        callRaw("POST", "/v4/magnet/delete", listOf("id" to id.toString()))
    }

    /** Restarts a failed magnet. */
    public suspend fun restartMagnet(id: Long) {
        callRaw("POST", "/v4/magnet/restart", listOf("id" to id.toString()))
    }

    /** Restarts failed magnets; each result has its own error. */
    public suspend fun restartMagnets(ids: List<Long>): List<RestartedMagnet> =
        call("POST", "/v4/magnet/restart", ids.map { "ids[]" to it.toString() }, RestartMagnetsData.serializer())
            .magnets.orEmpty()

    /** Follows the magnets with live mode: each [LiveMagnets.poll] only transfers changes. */
    public fun liveMagnets(): LiveMagnets = LiveMagnets(this)

    // ------------------------------------------------------------ saved links

    /** The links the user saved. */
    public suspend fun savedLinks(): List<SavedLink> =
        call("GET", "/v4/user/links", serializer = SavedLinksData.serializer()).links

    /** Saves links for later. */
    public suspend fun saveLinks(links: List<String>) {
        callRaw("POST", "/v4/user/links/save", links.map { "links[]" to it })
    }

    /** Deletes saved links. */
    public suspend fun deleteSavedLinks(links: List<String>) {
        callRaw("POST", "/v4/user/links/delete", links.map { "links[]" to it })
    }

    /** Links unlocked in the last 3 days, when history logging is enabled. */
    public suspend fun history(): List<SavedLink> =
        call("GET", "/v4/user/history", serializer = SavedLinksData.serializer()).links

    /** Deletes the whole recent links history. */
    public suspend fun purgeHistory() {
        callRaw("POST", "/v4/user/history/delete")
    }

    // --------------------------------------------------------------- reseller

    /** The reseller balance. Reseller accounts only. */
    public suspend fun voucherBalance(): Double =
        call("GET", "/v4/voucher/balance", serializer = VoucherBalanceData.serializer()).balance

    /** Takes vouchers (15, 30, 90, 180 or 365 days) from the reseller stock. */
    public suspend fun vouchers(durationDays: Int, count: Int): VouchersData =
        call("POST", "/v4/voucher/get", voucherParams(durationDays, count), VouchersData.serializer())

    /** Generates new vouchers from the reseller balance. */
    public suspend fun generateVouchers(durationDays: Int, count: Int): GeneratedVouchersData =
        call("POST", "/v4/voucher/generate", voucherParams(durationDays, count), GeneratedVouchersData.serializer())

    private fun hostsOnlyParam(hostsOnly: Boolean) = if (hostsOnly) listOf("hostsOnly" to "true") else emptyList()

    private fun voucherParams(durationDays: Int, count: Int) =
        listOf("duration" to durationDays.toString(), "nb" to count.toString())
}

/** Magnet status codes. 5 and above are errors. */
public object MagnetStatus {
    public const val QUEUED: Int = 0
    public const val DOWNLOADING: Int = 1
    public const val COMPRESSING: Int = 2
    public const val UPLOADING: Int = 3
    public const val READY: Int = 4
}

/** Delayed link statuses. */
public object DelayedStatus {
    public const val PROCESSING: Int = 1
    public const val READY: Int = 2
    public const val FAILED: Int = 3
}

/** Whether the magnet's files are available. */
public val Magnet.isReady: Boolean get() = statusCode == MagnetStatus.READY

/** Whether the magnet failed. */
public val Magnet.isError: Boolean get() = (statusCode ?: 0) > MagnetStatus.READY

/** Live mode state: the merged list of magnets, updated by [poll]. */
public class LiveMagnets internal constructor(private val client: AllDebrid) {
    public val session: Long = Random.nextLong(1L shl 31)
    private val mutex = Mutex()
    private var counter = 0L
    private var state = LinkedHashMap<Long, JsonObject>()

    /** Applies the changes since the last call and returns every magnet, newest first. */
    public suspend fun poll(): List<Magnet> = mutex.withLock {
        val data = client.callRaw(
            "POST", "/v4.1/magnet/status",
            listOf("session" to session.toString(), "counter" to counter.toString()),
        ).jsonObject
        if (data["fullsync"]?.jsonPrimitive?.booleanOrNull == true) state = LinkedHashMap()
        val changes = when (val m = data["magnets"]) {
            is JsonArray -> m.map { it.jsonObject }
            is JsonObject -> listOf(m)
            else -> emptyList()
        }
        for (change in changes) {
            val id = change["id"]?.jsonPrimitive?.longOrNull ?: continue
            if (change["deleted"]?.jsonPrimitive?.booleanOrNull == true) {
                state.remove(id)
            } else {
                state[id] = JsonObject(state[id].orEmpty() + change)
            }
        }
        counter = data["counter"]?.jsonPrimitive?.longOrNull ?: 0L
        state.values
            .map { AllDebrid.json.decodeFromJsonElement(Magnet.serializer(), it) }
            .sortedWith(compareByDescending<Magnet> { it.uploadDate ?: 0L }.thenByDescending { it.id })
    }
}

private val FORM = "application/x-www-form-urlencoded".toMediaType()
private val TORRENT = "application/x-bittorrent".toMediaType()

private suspend fun Call.await(): Response = suspendCancellableCoroutine { cont ->
    cont.invokeOnCancellation { cancel() }
    enqueue(object : Callback {
        override fun onFailure(call: Call, e: IOException) {
            cont.resumeWithException(e)
        }

        override fun onResponse(call: Call, response: Response) {
            cont.resume(response) { _, _, _ -> response.close() }
        }
    })
}
