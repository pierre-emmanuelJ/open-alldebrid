package io.github.pierreemmanuelj.alldebrid

import java.io.IOException

/** An error answered by the API (`"status": "error"`). Compare [code] with [ErrorCodes]. */
public class AllDebridException(
    /** The error code, e.g. `MAGNET_INVALID_ID`. */
    public val code: String,
    message: String,
    /** With `AUTH_BLOCKED`, the token to poll with [AllDebrid.verificationStatus]. */
    public val token: String? = null,
) : Exception("$code: $message") {
    /** The message of the API, without the code. */
    public val apiMessage: String = message
}

/** An HTTP status other than 200, such as 429 when rate limited. */
public class HttpException(
    public val status: Int,
    /** The API error in the body, when there is one. */
    public val apiError: AllDebridException? = null,
) : IOException(if (apiError != null) "HTTP $status: ${apiError.message}" else "HTTP $status") {
    /** Whether the API asked to slow down (HTTP 429 or 503). */
    public val isRateLimited: Boolean get() = status == 429 || status == 503
}

/** Whether this is an API error with the given code. */
public fun Throwable.isErrorCode(code: String): Boolean = when (this) {
    is AllDebridException -> this.code == code
    is HttpException -> apiError?.code == code
    else -> false
}
