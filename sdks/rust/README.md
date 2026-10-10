# open-alldebrid (Rust)

An async client for the [AllDebrid](https://alldebrid.com) API, built on the
[open-alldebrid OpenAPI spec](../../openapi.yaml): reqwest, serde and tokio.

```toml
[dependencies]
open-alldebrid = "1"
```

```rust
use open_alldebrid::{AllDebrid, error_codes};

let client = AllDebrid::new(api_key).user_agent("my-app");

let user = client.user().await?;

let uploaded = client.upload_magnets(&["magnet:?xt=urn:btih:..."]).await?;
let magnet = client.magnet(uploaded[0].id.unwrap()).await?;
let files = client.magnet_files(&[magnet.id]).await?;

match client.unlock("https://1fichier.com/?abcdefghij", None).await {
    Ok(u) => {
        let link = match u.delayed {
            Some(id) => client.wait_delayed_link(id.as_str(), Duration::from_secs(5)).await?,
            None => u.link.unwrap_or_default(),
        };
    }
    Err(e) if e.is_code(error_codes::LINK_DOWN) => println!("the file is gone"),
    Err(e) => return Err(e.into()),
}
```

## What it handles for you

- **Errors.** The API answers errors with HTTP 200 and `"status": "error"`.
  Every call returns them as `Error::Api` (with `code`, and `token` for
  `AUTH_BLOCKED`); other HTTP statuses are `Error::Http`
  (`is_rate_limited()` for 429/503).
- **List parameters** are sent the way the API expects:
  `magnets[]=a&magnets[]=b`.
- **Inconsistent shapes.** IDs sent as strings or numbers (`FlexString`), a
  magnet status that is an object for one magnet and a list otherwise, regexps
  that are a string or a list, and empty quota maps sent as `[]` all decode.
- **PIN login.** `get_pin()` then `wait_for_pin(&pin, interval)` returns the
  new key.
- **Delayed links.** `wait_delayed_link(id, interval)` polls until the link is
  ready.
- **Live mode.** `live_magnets().poll()` only transfers what changed and keeps
  the merged list, including deletions.

## Try it without an account

`demo_keys::PREMIUM`, `demo_keys::TRIAL` and `demo_keys::FREE` make the API
answer static sample data.

## Development

`src/models.rs` is generated from the spec by `scripts/gen_models.py` at the
root of the repository (needs python3 with PyYAML):

```sh
python3 ../../scripts/gen_models.py rust src/models.rs
cargo test                 # OFFLINE=1 skips the calls to the demo API
cargo clippy --all-targets -- -D warnings
```

The tests check that every operation of the spec has a method calling the
right path, and that `error_codes` matches the spec.
