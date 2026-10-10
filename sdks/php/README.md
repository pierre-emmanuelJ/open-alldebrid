# open-alldebrid (PHP)

A PHP client for the [AllDebrid](https://alldebrid.com) API, built on the
[open-alldebrid OpenAPI spec](../../openapi.yaml). PHP 8.2+, no dependencies
beyond ext-curl and ext-json; responses are readonly typed objects.

```sh
composer require pierre-emmanuelj/open-alldebrid
```

```php
use OpenAllDebrid\AllDebrid;
use OpenAllDebrid\AllDebridException;
use OpenAllDebrid\ErrorCode;

$client = new AllDebrid($apiKey, userAgent: 'my-app');

$user = $client->user();

[$uploaded] = $client->uploadMagnets(['magnet:?xt=urn:btih:...']);
$magnet = $client->magnet($uploaded->id);
$files = $client->magnetFiles([$magnet->id]);

try {
    $unlocked = $client->unlock('https://1fichier.com/?abcdefghij');
    $link = $unlocked->delayed !== null ? $client->waitDelayedLink($unlocked->delayed) : $unlocked->link;
} catch (AllDebridException $e) {
    if ($e->errorCode !== ErrorCode::LINK_DOWN) {
        throw $e;
    }
}
```

## What it handles for you

- **Errors.** The API answers errors with HTTP 200 and `"status": "error"`.
  Every call throws them as `AllDebridException` (with `errorCode`, and `token`
  for `AUTH_BLOCKED`); other HTTP statuses throw `HttpException`
  (`isRateLimited()` for 429/503).
- **List parameters** are sent the way the API expects:
  `magnets[]=a&magnets[]=b`.
- **Inconsistent shapes.** IDs sent as strings or numbers, a magnet status
  that is an object for one magnet and a list otherwise, regexps that are a
  string or a list, and empty quota maps sent as `[]` all decode to one type.
  Values of an unexpected type throw `UnexpectedValueException` instead of
  going through silently.
- **PIN login.** `getPin()` then `waitForPin($pin)` returns the new key.
- **Delayed links.** `waitDelayedLink($id)` polls until the link is ready.
- **Live mode.** `liveMagnets()->poll()` only transfers what changed and keeps
  the merged list, including deletions.

The HTTP layer is a one-method `Transport` interface; the default uses curl.

## Try it without an account

`AllDebrid::DEMO_PREMIUM_KEY`, `DEMO_TRIAL_KEY` and `DEMO_FREE_KEY` make the
API answer static sample data.

## Development

`src/Model/` and `src/ErrorCode.php` are generated from the spec by
`scripts/gen_models.py` at the root of the repository (`composer generate`,
needs python3 with PyYAML).

```sh
composer install
composer test        # OFFLINE=1 skips the calls to the demo API
composer analyse     # PHPStan, level max
```

The tests check that every operation of the spec has a method calling the
right path, and that `ErrorCode` matches the spec.
