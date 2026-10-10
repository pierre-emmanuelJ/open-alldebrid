# AllDebrid Go SDK

A Go client for the [AllDebrid](https://alldebrid.com) API, built on the
[open-alldebrid OpenAPI spec](../../openapi.yaml). No dependencies beyond the
standard library.

```sh
go get github.com/pierre-emmanuelJ/open-alldebrid/sdks/go
```

```go
import alldebrid "github.com/pierre-emmanuelJ/open-alldebrid/sdks/go"

client := alldebrid.New(apiKey, alldebrid.WithUserAgent("my-app"))

user, err := client.User(ctx)

magnets, err := client.UploadMagnets(ctx, "magnet:?xt=urn:btih:...")
m, err := client.Magnet(ctx, magnets[0].ID)
files, err := client.MagnetFiles(ctx, m.ID)

link, err := client.Unlock(ctx, "https://1fichier.com/?abcdefghij", "")
```

## What it handles for you

- **Errors.** The API answers errors with HTTP 200 and `"status": "error"`.
  Every method returns them as `*alldebrid.Error`; compare codes with
  `errors.Is(err, alldebrid.CodeMagnetInvalidID)`. HTTP 429/503 are
  `*alldebrid.HTTPError` (`alldebrid.IsRateLimited(err)`).
- **List parameters** are sent the way the API (PHP) expects:
  `magnets[]=a&magnets[]=b`.
- **Inconsistent types.** IDs that come as strings or numbers
  (`FlexString`), a magnet status that is an object for one magnet and a list
  otherwise (`Magnets`), regexps that are a string or a list (`Strings`), and
  empty quota maps sent as `[]` (`Quotas`).
- **PIN login.** `GetPin` then `WaitForPin` returns the new key.
- **Delayed links.** `WaitDelayedLink` polls until the link is ready.
- **Live mode.** `LiveMagnets().Poll` only transfers what changed and keeps the
  merged list, including deletions.

## Try it without an account

The demo keys `alldebrid.DemoPremiumKey`, `DemoTrialKey` and `DemoFreeKey`
make the API answer static sample data.

## Development

`types.gen.go` is generated from the spec with
[oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) and the overlay
in `overlay.yaml`:

```sh
go generate ./...
go test ./...          # -short skips the calls to the demo API
```

The tests check that every operation of the spec has a client method that
calls the right path, and that the error codes match the spec.
