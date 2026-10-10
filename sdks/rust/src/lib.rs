//! Async client for the [AllDebrid](https://alldebrid.com) API (v4 and v4.1),
//! built from the [open-alldebrid](https://github.com/pierre-emmanuelJ/open-alldebrid)
//! OpenAPI spec.
//!
//! ```no_run
//! use open_alldebrid::{AllDebrid, error_codes};
//!
//! # async fn run() -> open_alldebrid::Result<()> {
//! let client = AllDebrid::new("api key").user_agent("my-app");
//! let user = client.user().await?;
//!
//! match client.unlock("https://1fichier.com/?abcdefghij", None).await {
//!     Ok(unlocked) => println!("{:?}", unlocked.link),
//!     Err(e) if e.is_code(error_codes::LINK_DOWN) => println!("the file is gone"),
//!     Err(e) => return Err(e),
//! }
//! # Ok(()) }
//! ```
//!
//! The models in [`models`] are generated from the spec; the client methods
//! are written by hand on top of them, one per operation.

mod client;
mod error;
mod flex;
mod live;
#[allow(missing_docs)]
pub mod models;

pub use client::{AllDebrid, DEFAULT_BASE_URL, TorrentFile, delayed_status, demo_keys, magnet_status};
pub use error::{Error, Result};
pub use flex::FlexString;
pub use live::LiveMagnets;
pub use models::error_codes;
