# open-alldebrid (Python)

A typed client for the [AllDebrid](https://alldebrid.com) API, built on the
[open-alldebrid OpenAPI spec](https://github.com/pierre-emmanuelJ/open-alldebrid/blob/master/openapi.yaml).
Sync and async, with [httpx](https://www.python-httpx.org/) as its only
dependency. Python 3.11+.

```sh
pip install open-alldebrid
```

```python
from open_alldebrid import AllDebrid, AllDebridError

with AllDebrid(api_key, user_agent="my-app") as client:
    user = client.user()

    [uploaded] = client.upload_magnets(["magnet:?xt=urn:btih:..."])
    magnet = client.magnet(uploaded["id"])
    files = client.magnet_files([magnet["id"]])

    try:
        unlocked = client.unlock("https://1fichier.com/?abcdefghij")
        delayed = unlocked.get("delayed")
        link = client.wait_delayed_link(delayed) if delayed else unlocked["link"]
    except AllDebridError as err:
        if err.code != "LINK_DOWN":
            raise
```

The async client has the same methods:

```python
from open_alldebrid import AsyncAllDebrid

async with AsyncAllDebrid(api_key) as client:
    magnets = await client.magnets()
```

## What it handles for you

- **Errors.** The API answers errors with HTTP 200 and `"status": "error"`.
  Every method raises them as `AllDebridError` (with `code`, and `token` for
  `AUTH_BLOCKED`). Other HTTP statuses raise `HTTPError`; `err.rate_limited`
  tells a 429/503 apart.
- **List parameters** are sent the way the API (PHP) expects:
  `magnets[]=a&magnets[]=b`.
- **Inconsistent shapes.** `magnets()` always returns a list and `magnet(id)`
  one dict, though the API switches between the two; empty quota maps sent as
  `[]` come back as `{}`.
- **PIN login.** `get_pin()` then `wait_for_pin(pin)` returns the new key.
- **Delayed links.** `wait_delayed_link(id)` polls until the link is ready.
- **Live mode.** `live_magnets().poll()` only transfers what changed and keeps
  the merged list, including deletions.

Responses are plain dicts typed with `TypedDict`s generated from the spec
(`open_alldebrid.models`), so editors and mypy know every field.

## Try it without an account

`DEMO_KEYS["premium"]`, `DEMO_KEYS["trial"]` and `DEMO_KEYS["free"]` make the
API answer static sample data.

## Development

```sh
pip install -e . --group dev        # pip >= 25.1
./scripts/generate.sh               # models.py from ../../openapi.yaml
python scripts/unasync.py           # _sync_client.py from _async_client.py
ruff check . && ruff format --check . && mypy && pytest
```

`OFFLINE=1 pytest` skips the calls to the demo API. The tests check that every
operation of the spec has a method calling the right path, and that the sync
client is up to date with the async one.
