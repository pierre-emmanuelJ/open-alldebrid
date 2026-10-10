//! Fields whose JSON type varies between responses.

use std::collections::HashMap;
use std::fmt;

use serde::de::{self, Deserializer};
use serde::{Deserialize, Serialize, Serializer};
use serde_json::Value;

use crate::models::Magnet;

/// A value the API sends either as a JSON string or as a JSON number, such as
/// a delayed link ID. It keeps the value as a string.
#[derive(Debug, Clone, PartialEq, Eq, Hash, Default)]
pub struct FlexString(pub String);

impl FlexString {
    /// The value as sent by the API.
    #[must_use]
    pub fn as_str(&self) -> &str {
        &self.0
    }

    /// The value as an integer, when it is one.
    #[must_use]
    pub fn as_i64(&self) -> Option<i64> {
        self.0.parse().ok()
    }
}

impl fmt::Display for FlexString {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(&self.0)
    }
}

impl<'de> Deserialize<'de> for FlexString {
    fn deserialize<D: Deserializer<'de>>(d: D) -> Result<Self, D::Error> {
        match Value::deserialize(d)? {
            Value::String(s) => Ok(FlexString(s)),
            Value::Number(n) => Ok(FlexString(n.to_string())),
            other => Err(de::Error::custom(format!("expected a string or a number, got {other}"))),
        }
    }
}

impl Serialize for FlexString {
    fn serialize<S: Serializer>(&self, s: S) -> Result<S::Ok, S::Error> {
        s.serialize_str(&self.0)
    }
}

fn convert<'de, D: Deserializer<'de>, T: serde::de::DeserializeOwned>(value: Value) -> Result<T, D::Error> {
    serde_json::from_value(value).map_err(de::Error::custom)
}

/// One string or a list of strings.
pub(crate) fn strings<'de, D: Deserializer<'de>>(d: D) -> Result<Vec<String>, D::Error> {
    match Value::deserialize(d)? {
        Value::String(s) => Ok(vec![s]),
        other => convert::<D, _>(other),
    }
}

pub(crate) fn opt_strings<'de, D: Deserializer<'de>>(d: D) -> Result<Option<Vec<String>>, D::Error> {
    strings(d).map(Some)
}

/// A list of magnets, or a single magnet object.
pub(crate) fn magnets<'de, D: Deserializer<'de>>(d: D) -> Result<Vec<Magnet>, D::Error> {
    match Value::deserialize(d)? {
        Value::Null => Ok(Vec::new()),
        one @ Value::Object(_) => Ok(vec![convert::<D, _>(one)?]),
        other => convert::<D, _>(other),
    }
}

#[allow(dead_code)]
pub(crate) fn opt_magnets<'de, D: Deserializer<'de>>(d: D) -> Result<Option<Vec<Magnet>>, D::Error> {
    magnets(d).map(Some)
}

/// Quotas per host; the API sends an empty map as an empty list.
pub(crate) fn quotas<'de, D: Deserializer<'de>>(d: D) -> Result<HashMap<String, i64>, D::Error> {
    match Value::deserialize(d)? {
        Value::Array(a) if a.is_empty() => Ok(HashMap::new()),
        Value::Array(_) => Err(de::Error::custom("quotas sent as a non-empty list")),
        other => convert::<D, _>(other),
    }
}

#[allow(dead_code)]
pub(crate) fn opt_quotas<'de, D: Deserializer<'de>>(d: D) -> Result<Option<HashMap<String, i64>>, D::Error> {
    quotas(d).map(Some)
}
