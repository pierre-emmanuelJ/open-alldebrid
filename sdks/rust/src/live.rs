//! Live mode: incremental magnet status.

use std::collections::HashMap;
use std::hash::{BuildHasher, RandomState};
use std::time::{SystemTime, UNIX_EPOCH};

use serde_json::{Map, Value};

use crate::client::AllDebrid;
use crate::error::{Error, Result};
use crate::models::Magnet;

/// Live mode state: the merged list of magnets, updated by [`LiveMagnets::poll`].
#[derive(Debug)]
pub struct LiveMagnets {
    client: AllDebrid,
    session: u64,
    counter: i64,
    state: HashMap<i64, Map<String, Value>>,
}

impl LiveMagnets {
    pub(crate) fn new(client: AllDebrid) -> Self {
        let nanos = SystemTime::now().duration_since(UNIX_EPOCH).map(|d| d.as_nanos()).unwrap_or_default();
        let session = RandomState::new().hash_one(nanos) % (1 << 31);
        Self { client, session, counter: 0, state: HashMap::new() }
    }

    /// The session ID sent to the API.
    #[must_use]
    pub fn session(&self) -> u64 {
        self.session
    }

    /// Applies the changes since the last call and returns every magnet, newest first.
    pub async fn poll(&mut self) -> Result<Vec<Magnet>> {
        let params =
            vec![("session".to_owned(), self.session.to_string()), ("counter".to_owned(), self.counter.to_string())];
        let data = self.client.call_raw(reqwest::Method::POST, "/v4.1/magnet/status", params, None).await?;
        if data.get("fullsync").and_then(Value::as_bool) == Some(true) {
            self.state.clear();
        }
        let changes = match data.get("magnets") {
            Some(Value::Array(list)) => list.clone(),
            Some(one @ Value::Object(_)) => vec![one.clone()],
            _ => Vec::new(),
        };
        for change in changes {
            let Value::Object(change) = change else { continue };
            let Some(id) = change.get("id").and_then(Value::as_i64) else { continue };
            if change.get("deleted").and_then(Value::as_bool) == Some(true) {
                self.state.remove(&id);
            } else {
                self.state.entry(id).or_default().extend(change);
            }
        }
        self.counter = data.get("counter").and_then(Value::as_i64).unwrap_or(0);

        let mut magnets = self
            .state
            .values()
            .map(|m| serde_json::from_value::<Magnet>(Value::Object(m.clone())))
            .collect::<std::result::Result<Vec<_>, _>>()
            .map_err(|e| Error::Decode(e.to_string()))?;
        magnets.sort_by_key(|m| std::cmp::Reverse((m.upload_date.unwrap_or(0), m.id)));
        Ok(magnets)
    }
}
