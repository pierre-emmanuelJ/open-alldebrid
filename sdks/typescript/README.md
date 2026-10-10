# open-alldebrid (TypeScript)

A typed client for the [AllDebrid](https://alldebrid.com) API, built on the
[open-alldebrid OpenAPI spec](https://github.com/pierre-emmanuelJ/open-alldebrid/blob/master/openapi.yaml).
No runtime dependencies: it uses `fetch`, so it runs in Node 18+, Deno, Bun,
browsers and browser extensions.

```sh
npm install open-alldebrid
```

```ts
import { AllDebrid, isErrorCode } from "open-alldebrid";

const client = new AllDebrid({ apiKey: process.env.ALLDEBRID_API_KEY, userAgent: "my-app" });

const user = await client.user();

const [uploaded] = await client.uploadMagnets(["magnet:?xt=urn:btih:..."]);
const magnet = await client.magnet(uploaded!.id!);
const files = await client.magnetFiles([magnet.id]);

try {
  const unlocked = await client.unlock("https://1fichier.com/?abcdefghij");
  const link = unlocked.delayed ? await client.waitDelayedLink(unlocked.delayed) : unlocked.link;
} catch (err) {
  if (isErrorCode(err, "LINK_DOWN")) console.log("the file is gone");
  else throw err;
}
```

## What it handles for you

- **Errors.** The API answers errors with HTTP 200 and `"status": "error"`.
  Every method throws them as `AllDebridError` (with `code`, and `token` for
  `AUTH_BLOCKED`). Other HTTP statuses throw `HttpError`; `err.rateLimited`
  tells a 429/503 apart.
- **List parameters** are sent the way the API (PHP) expects:
  `magnets[]=a&magnets[]=b`.
- **Inconsistent shapes.** `magnets()` always returns a list and `magnet(id)`
  one object, though the API switches between the two; empty quota maps sent
  as `[]` come back as `{}`.
- **PIN login.** `getPin()` then `waitForPin(pin)` returns the new key.
- **Delayed links.** `waitDelayedLink(id)` polls until the link is ready.
- **Live mode.** `liveMagnets().poll()` only transfers what changed and keeps
  the merged list, including deletions.

Every poller takes `{ intervalMs, signal }`, so an `AbortSignal` stops it.

## Try it without an account

`DEMO_KEYS.premium`, `DEMO_KEYS.trial` and `DEMO_KEYS.free` make the API
answer static sample data.

## Development

```sh
npm install
npm run generate   # src/schema.gen.ts from ../../openapi.yaml
npm run check      # typecheck, tests, build, publint, attw
```

`OFFLINE=1 npm test` skips the calls to the demo API. The tests check that
every operation of the spec has a method calling the right path.
