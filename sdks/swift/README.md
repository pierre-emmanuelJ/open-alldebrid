# OpenAllDebrid (Swift)

A Swift client for the [AllDebrid](https://alldebrid.com) API, built on the
[open-alldebrid OpenAPI spec](../../openapi.yaml): async/await on URLSession,
no dependencies. iOS 15, macOS 12, tvOS 15, watchOS 8 and Linux; Swift 6
strict concurrency.

```swift
// Package.swift
.package(url: "https://github.com/pierre-emmanuelJ/open-alldebrid", from: "1.0.0")
```

```swift
import OpenAllDebrid

let client = AllDebrid(apiKey: apiKey, userAgent: "my-app")

let user = try await client.user()

let uploaded = try await client.uploadMagnets(["magnet:?xt=urn:btih:..."])
let magnet = try await client.magnet(id: uploaded[0].id!)
let files = try await client.magnetFiles(ids: [magnet.id])

do {
    let unlocked = try await client.unlock("https://1fichier.com/?abcdefghij")
    let link = if let id = unlocked.delayed { try await client.waitDelayedLink(id.value) } else { unlocked.link }
} catch let error as AllDebridError where error.code == ErrorCode.linkDown {
    print("the file is gone")
}
```

## What it handles for you

- **Errors.** The API answers errors with HTTP 200 and `"status": "error"`.
  Every call throws them as `AllDebridError.api(code:message:token:)`; other
  HTTP statuses are `.http(status:code:)` (`isRateLimited` for 429/503).
- **List parameters** are sent the way the API expects:
  `magnets[]=a&magnets[]=b`.
- **Inconsistent shapes.** IDs sent as strings or numbers (`FlexString`), a
  magnet status that is an object for one magnet and a list otherwise, regexps
  that are a string or a list, and empty quota maps sent as `[]` all decode.
  Responses missing a required field throw.
- **PIN login.** `getPin()` then `waitForPin(pin)` returns the new key.
- **Delayed links.** `waitDelayedLink(id)` polls until the link is ready.
- **Live mode.** `liveMagnets().poll()` (an actor) only transfers what changed
  and keeps the merged list, including deletions.

## Try it without an account

`AllDebrid.DemoKey.premium`, `.trial` and `.free` make the API answer static
sample data.

## Development

`Sources/OpenAllDebrid/Models.generated.swift` is generated from the spec by
`scripts/gen_models.py` at the root of the repository (needs python3 with
PyYAML):

```sh
python3 ../../scripts/gen_models.py swift Sources/OpenAllDebrid/Models.generated.swift
swift test                       # OFFLINE=1 skips the calls to the demo API
swift format lint -r --strict Sources Tests
```

The tests check that every operation of the spec has a method calling the
right path, and that `ErrorCode` matches the spec.
