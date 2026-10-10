# Errata for the official AllDebrid API docs

Differences between [docs.alldebrid.com](https://docs.alldebrid.com/) and the
live API, found while writing `openapi.yaml`. Checked on 2026-10-09 and
2026-10-10 with the public demo keys and with a real premium account.

The contract tests in [`contract/`](contract/) replay the docs' demo cases and
a real-account session against the spec, so each "API behaviour" item below
is reproducible.

## 1. The API does not do what the docs say

| # | Endpoint | The docs say | The API does | Effect on clients |
|---|----------|--------------|--------------|-------------------|
| 1.1 | `POST /v4.1/magnet/status` | The `status` parameter filters magnets (`active`, `ready`, `expired`, `error`). | The filter is ignored: every value returns every magnet (tried with GET, form and multipart bodies). | Apps that rely on it download the whole list and show finished magnets as active. |
| 1.2 | `POST /v4.1/magnet/status` with `id` | `magnets` is an array of magnet objects. | `magnets` is a single object. Live mode deltas can also be a single object. | Typed clients crash on decoding. |
| 1.3 | `POST /v4/link/unlock` | Alternative qualities are in `streams`; `id` is an integer; `delayed` is an integer. | The field is `streaming`; `id` is a string (`"3hk19hz3205"`); `delayed` is a string (`"123456"`), while the `/v4/link/streaming` example shows it as an integer. | Code written from the docs never sees the qualities and fails to parse IDs. |
| 1.4 | `POST /v4/magnet/upload/file` | Each result has `file`, `name`, `hash`, `size`, `ready`, `id`. | Successful results have `magnet` instead of `file` (plus `filename_original`); only failed ones have `file`. | Apps cannot match results to uploaded files by the documented field. |
| 1.5 | `POST /v4/magnet/files` | `id` is an integer, and the parameter is `id` (array of integer). | `id` comes back as a string. Only `id[]` works: a scalar `id` answers `MAGNET_INVALID_ID`. A JSON body is not accepted either. | |
| 1.6 | `POST /v4/magnet/restart` with an unknown `id` | `MAGNET_INVALID_ID`. | `MAGNET_INVALID_URI`. | Error handling keyed on the documented code misses it. |
| 1.7 | `/v4/voucher/*` from a non-reseller account | `ACCOUNT_INVALID` exists for "doesn't have access to this endpoint". | `GENERIC`. | |
| 1.8 | `POST /v4/user/notification/clear` without `code` | `code` is required. | Answers success. | |
| 1.9 | All endpoints | "Some errors will return an HTTP code 401 or 429" (All errors section). | Every error tried came back with HTTP 200, except unknown endpoints (404); the 14/11/2023 changelog entry says 400 and 401 were removed. Rate limiting (429/503) was not triggered on purpose. | |
| 1.10 | `/v4/magnet/status` | Deprecated in favour of `/v4.1/`. | Discontinued: it answers the undocumented error code `DISCONTINUED`. | Old clients get an error code they cannot know. |

## 2. Undocumented behaviour and fields

These are in the spec so typed clients can use them; the docs should list them.

- **Live mode deletions.** A deleted magnet arrives in a delta as
  `{"id": 123, "deleted": true}`. Without this a client keeps showing deleted
  magnets until the next full sync.
- **`/v4.1/magnet/status` magnet fields:** `hash`, `type` (`m` magnet, `b`
  torrent file), `notified`, `version`, `nbLinks`, and `processingPerc` while
  processing.
- **Magnet status texts** differ from the status code table: code 7 is
  "No peer after 30 minutes" (docs: not downloaded in 20 min), code 10 "Download
  took more than 3 days" (docs: 72h), code 11 "Expired - Files removed". Codes 12
  and 13 have the same description in the table.
- **`/v4/user`:** a `history` boolean (whether history logging is enabled).
- **`/v4/hosts` and `/v4.1/user/hosts`:** every entry has both `regexp` and
  `regexps`, and `regexp` is sometimes an array; hosts can have `display` and
  `hardRedirect`; redirectors have `heavy`.
- **`/v4/user/history`:** entries have `link_dl` and `error` (seen in demo responses).
- **`/v4/magnet/upload`:** results have `filename_original`; the endpoint
  accepts `magnets[]`, `magnets` and `magnet`.
- **`/v4/link/redirector`:** the error `REDIRECTOR_RATE_LIMITED` ("Too much
  request to unlock folder"), returned by the docs' own demo button.
- **`/v4.1/ping`:** a `subversion` boolean.
- **Authentication:** besides the documented `Authorization: Bearer` header,
  `apikey` works as a query parameter and as a form field (the docs' own demo
  buttons use the form field).
- **Uploading a magnet that is already in the account** returns its existing ID.

## 3. Mistakes in the documentation

### Wrong URLs or parameters in examples

- `/user/links/save`: the curl example posts to `/v4/magnet/delete`.
- `/magnet/files`: the curl example uses `/v4.1/magnet/files`; the HTTP Request
  line and "All endpoints" say `/v4/magnet/files` (both work).
- `/pin/check`: the full PHP flow still calls the deprecated `/v4/pin/get`.
- `/user/history/delete`, `/voucher/balance`, `/voucher/get` and
  `/voucher/generate`: the HTTP Request lines lack `/v4/`.
- `/user/history/delete`: the curl URL ends with a stray `'`.
- `/user/notification/clear`: the curl URL is missing its closing quote.
- `/magnet/upload`: the PHP and curl examples still pass `?apikey=` in the URL
  while also sending the header. The PHP example has a 20-character key in place
  of the placeholder; if it was ever a real key, it should be revoked.
- `/magnet/upload/file`: PHP uses `files[0]`, curl uses `files[]`, the table
  says `files`.
- Examples mix `http://` and `https://`, and curl `-F` (multipart) with PHP
  `http_build_query` (urlencoded).
- The voucher endpoints are missing from "All endpoints".

### Invalid JSON in examples

Trailing commas or missing braces in the examples of `/pin/get`, `/hosts`,
`/user/verif`, `/link/infos`, `/link/unlock` (missing closing brace),
`/link/streaming` (second example missing its opening brace),
`/magnet/status`, live mode, `/magnet/files`, `/user/links` and
`/user/history`. They cannot be pasted into a JSON parser or a test.

### Wrong types in tables

- `/user`: `premiumUntil` is documented as Integer but the example is a string;
  `limitedHostersQuotas` is documented as Array but is an object.
- `/v4.1/user/hosts`: the example shows `regexp`, the table `regexps`.
- Deprecated `/v4/magnet/status`: the Link object's `link` is typed Integer.
- Demo responses: the example shows `"demo": "true"` (a string); the API sends
  `true`.

### Structure

- The "All errors" table is cut at `FREEDAYS_INVALID_PROVIDER`
  ("Invalid pro16/10/2024") and the following `FREEDAYS_*` rows are missing;
  they only appear in the PHP array above it.
- The third reseller section is titled "Balance" but documents
  `/voucher/generate`, with the description of `/voucher/get`; its curl example
  asks for 5 vouchers, the PHP one for 2, and the response shows 3.
- The changelog lists the 2020 entries twice; the 01/02/2020 entry mentions
  `user/links/add` (the endpoint is `user/links/save`).
- `/user/links/delete` marks `links` as optional.
- The hosts introduction links to `/user/host` (missing `s`); the
  `/magnet/status` table links to `/magnets/files`.
- `/user/verif` and `/user/verif/resend` document `PIN_EXPIRED` as "Verification
  token is invalid".

### Typos

"Rgexp", "Interger"/"interger", "Thos encrypted links", "se againnt",
"We can referrer", "An orror occured", "doens't", "Sent link in not valid",
"No enough balance", "accross", "Proccessing", "the endpoin", "successfull",
"reponse", "reseted", and in API messages "Links were sucessfully saved".
