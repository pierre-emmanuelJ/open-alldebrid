# Open AllDebrid

[![CI](https://github.com/pierre-emmanuelJ/open-alldebrid/actions/workflows/ci.yml/badge.svg)](https://github.com/pierre-emmanuelJ/open-alldebrid/actions/workflows/ci.yml)
[![Drift check](https://github.com/pierre-emmanuelJ/open-alldebrid/actions/workflows/drift.yml/badge.svg)](https://github.com/pierre-emmanuelJ/open-alldebrid/actions/workflows/drift.yml)

An OpenAPI 3.1 description of the [AllDebrid](https://alldebrid.com) API (v4
and v4.1), checked against the live API, with SDKs for Go, TypeScript,
Python, Kotlin, PHP, Rust and C# built on it.

**[Interactive reference →](https://pierre-emmanuelj.github.io/open-alldebrid/)**
Every endpoint, with real response examples, SDK samples, and a "Try it"
button that calls the API with the public demo key.

## What is here

| | |
|---|---|
| [`openapi.yaml`](openapi.yaml) | The spec: 30 operations, named schemas, every error code, real examples. |
| [`sdks/go`](sdks/go) | Go client, no dependencies. `go get github.com/pierre-emmanuelJ/open-alldebrid/sdks/go` |
| [`sdks/typescript`](sdks/typescript) | TypeScript client for Node, Deno, Bun and browsers, no runtime dependencies. |
| [`sdks/python`](sdks/python) | Python client, sync and async, on httpx. |
| [`sdks/kotlin`](sdks/kotlin) | Kotlin client for the JVM and Android: coroutines, OkHttp, kotlinx.serialization. |
| [`sdks/php`](sdks/php) | PHP 8.2+ client with readonly typed models, no dependencies beyond ext-curl. |
| [`sdks/rust`](sdks/rust) | Async Rust client on reqwest, serde and tokio. |
| [`sdks/csharp`](sdks/csharp) | .NET 8/10 client on HttpClient and System.Text.Json, no dependencies. |
| [`contract/`](contract) | Tests that replay the API against the spec. |
| [`ERRATA.md`](ERRATA.md) | Where the official docs and the live API disagree. |

The TypeScript, Python, Kotlin, PHP, Rust and C# packages are ready to publish
(on npm, PyPI, Maven Central, Packagist, crates.io and NuGet) but are not
published yet.

## Why a spec "checked against the live API"

The [official docs](https://docs.alldebrid.com) are written by hand, and the
API has drifted from them: `/v4/link/unlock` returns `streaming` where the docs
say `streams`, `/v4.1/magnet/status` returns one object instead of a list when
you ask for one magnet, the `status` filter of that endpoint is ignored, IDs
documented as integers are strings. [`ERRATA.md`](ERRATA.md) lists every case.

This spec follows what the API actually sends, and keeps following it:

- **Contract tests** call every operation (with the docs' own demo cases and
  the public demo keys), validate each response against the spec, and fail on
  any field the spec does not declare. With a real API key they also run
  against an account (`make test-live`).
- **Every response example** in the spec is a real demo response, validated
  against its schema.
- **A weekly job** reruns the contract tests and watches the changelog of the
  official docs, and opens an issue when either changes.

## The SDKs

The SDKs share one design: models generated from the spec, and a small
hand-written client on top, one method per operation. Each SDK's tests check
that every operation of the spec has a method that calls the right path, so
the spec and the SDKs cannot drift apart.

They handle what makes this API awkward:

- errors come with HTTP 200 and `"status": "error"`: the SDKs raise typed
  errors you can match by code;
- list parameters are PHP-style (`magnets[]=a&magnets[]=b`), which most
  OpenAPI generators get wrong;
- fields whose JSON type varies (string or number IDs, a list or a single
  object, an empty map sent as `[]`) come out normalized;
- the PIN login, delayed links and live mode (incremental magnet status,
  including the undocumented deletion marker) have helpers.

```go
client := alldebrid.New(apiKey)
magnets, err := client.UploadMagnets(ctx, "magnet:?xt=urn:btih:...")
```

```ts
const client = new AllDebrid({ apiKey });
const [magnet] = await client.uploadMagnets(["magnet:?xt=urn:btih:..."]);
```

```python
client = AllDebrid(api_key)
[magnet] = client.upload_magnets(["magnet:?xt=urn:btih:..."])
```

## Use the spec for another language

The spec works with any OpenAPI 3.1 tool. If you generate a client, watch
for the PHP-style list parameters: the `magnets[]`, `link[]`, `id[]`, `ids[]`
and `links[]` fields must be sent as repeated keys with the brackets.

## Development

```sh
make lint       # Redocly lint of openapi.yaml
make generate   # regenerate every SDK's models from the spec
make test       # contract tests and all SDK tests (calls the demo API)
make test-live  # contract tests on a real account (ALLDEBRID_API_KEY)
make drift      # has the official changelog changed?
```

When the API changes: update `openapi.yaml`, run `make generate`, and add the
change to the SDK clients if an operation was added or renamed; the tests say
what is missing.

## History

The first version (2022) translated the AllDebrid PHP documentation into an
OpenAPI 3.0 spec with a generated Go client. Version 1.0 (2026) rewrote it
for the v4.1 API, checked it against the live API, and added the SDKs, the
contract tests and the docs page.

## License

[MIT](LICENSE). AllDebrid is a trademark of its owner; this project is not
affiliated with it.
