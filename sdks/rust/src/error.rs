//! Errors returned by the client.

use std::fmt;

/// An error of the client.
#[derive(Debug)]
#[non_exhaustive]
pub enum Error {
    /// An error answered by the API (`"status": "error"`). Compare `code` with
    /// the constants of [`crate::error_codes`].
    Api {
        /// The error code, e.g. `MAGNET_INVALID_ID`.
        code: String,
        /// The message of the API.
        message: String,
        /// With `AUTH_BLOCKED`, the token to poll with
        /// [`crate::AllDebrid::verification_status`].
        token: Option<String>,
    },
    /// An HTTP status other than 200, such as 429 when rate limited.
    Http {
        /// The HTTP status.
        status: u16,
        /// The API error code in the body, when there is one.
        code: Option<String>,
    },
    /// The request could not be sent or the response could not be read.
    Transport(reqwest::Error),
    /// The response did not have the expected shape.
    Decode(String),
}

impl Error {
    /// The API error code, when this error carries one.
    #[must_use]
    pub fn code(&self) -> Option<&str> {
        match self {
            Error::Api { code, .. } => Some(code),
            Error::Http { code, .. } => code.as_deref(),
            _ => None,
        }
    }

    /// Whether this is an API error with the given code.
    #[must_use]
    pub fn is_code(&self, code: &str) -> bool {
        self.code() == Some(code)
    }

    /// Whether the API asked to slow down (HTTP 429 or 503).
    #[must_use]
    pub fn is_rate_limited(&self) -> bool {
        matches!(self, Error::Http { status: 429 | 503, .. })
    }
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Error::Api { code, message, .. } => write!(f, "alldebrid: {code}: {message}"),
            Error::Http { status, code: Some(code) } => write!(f, "alldebrid: HTTP {status}: {code}"),
            Error::Http { status, code: None } => write!(f, "alldebrid: HTTP {status}"),
            Error::Transport(e) => write!(f, "alldebrid: {e}"),
            Error::Decode(e) => write!(f, "alldebrid: unexpected response: {e}"),
        }
    }
}

impl std::error::Error for Error {
    fn source(&self) -> Option<&(dyn std::error::Error + 'static)> {
        match self {
            Error::Transport(e) => Some(e),
            _ => None,
        }
    }
}

impl From<reqwest::Error> for Error {
    fn from(e: reqwest::Error) -> Self {
        Error::Transport(e)
    }
}

/// The result type of the client.
pub type Result<T> = std::result::Result<T, Error>;
