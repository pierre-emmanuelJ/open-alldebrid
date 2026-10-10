# open-alldebrid (Kotlin)

A Kotlin client for the [AllDebrid](https://alldebrid.com) API, built on the
[open-alldebrid OpenAPI spec](../../openapi.yaml): coroutines, OkHttp and
kotlinx.serialization. Runs on the JVM 11+ and on Android.

```kotlin
val client = AllDebrid(apiKey, userAgent = "my-app")

val user = client.user()

val (uploaded) = client.uploadMagnets(listOf("magnet:?xt=urn:btih:..."))
val magnet = client.magnet(uploaded.id!!)
val files = client.magnetFiles(listOf(magnet.id))

try {
    val unlocked = client.unlock("https://1fichier.com/?abcdefghij")
    val link = unlocked.delayed?.let { client.waitDelayedLink(it.value) } ?: unlocked.link
} catch (e: AllDebridException) {
    if (e.code != ErrorCodes.LINK_DOWN) throw e
}
```

## What it handles for you

- **Errors.** The API answers errors with HTTP 200 and `"status": "error"`.
  Every call throws them as `AllDebridException` (with `code`, and `token` for
  `AUTH_BLOCKED`); other HTTP statuses throw `HttpException`
  (`isRateLimited` for 429/503).
- **List parameters** are sent the way the API (PHP) expects:
  `magnets[]=a&magnets[]=b`.
- **Inconsistent shapes.** IDs sent as strings or numbers (`FlexString`), a
  magnet status that is an object for one magnet and a list otherwise, regexps
  that are a string or a list, and empty quota maps sent as `[]` all decode.
- **PIN login.** `getPin()` then `waitForPin(pin)` returns the new key.
- **Delayed links.** `waitDelayedLink(id)` polls until the link is ready.
- **Live mode.** `liveMagnets().poll()` only transfers what changed and keeps
  the merged list, including deletions.

Every call is cancellable: cancelling the coroutine cancels the HTTP call.

## Try it without an account

`AllDebrid.DEMO_PREMIUM_KEY`, `DEMO_TRIAL_KEY` and `DEMO_FREE_KEY` make the API
answer static sample data.

## Development

`Models.kt` is generated from the spec by `scripts/gen_models.py` at the root
of the repository (`./gradlew generateModels`, needs python3 with PyYAML).

```sh
./gradlew build       # compiles and runs the tests
OFFLINE=1 ./gradlew test   # without the calls to the demo API
```

The tests check that every operation of the spec has a method calling the
right path, and that `ErrorCodes` matches the spec.
