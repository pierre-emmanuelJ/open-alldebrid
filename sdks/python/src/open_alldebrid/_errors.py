from __future__ import annotations

from typing import get_args

from .models import ErrorCode

ERROR_CODES: tuple[str, ...] = get_args(ErrorCode)
"""Every error code the spec lists. The API may return others."""


class AllDebridError(Exception):
    """An error answered by the API (``"status": "error"``)."""

    def __init__(self, code: str, message: str, token: str | None = None) -> None:
        super().__init__(f"{code}: {message}")
        self.code = code
        """The error code, e.g. ``MAGNET_INVALID_ID``."""
        self.message = message
        self.token = token
        """With ``AUTH_BLOCKED``, the token to poll with ``verification_status``."""


class HTTPError(Exception):
    """An HTTP status other than 200, such as 429 when rate limited."""

    def __init__(self, status_code: int, api_error: AllDebridError | None = None) -> None:
        detail = f": {api_error}" if api_error else ""
        super().__init__(f"HTTP {status_code}{detail}")
        self.status_code = status_code
        self.api_error = api_error
        """The API error in the body, when there is one."""

    @property
    def rate_limited(self) -> bool:
        """Whether the API asked to slow down (HTTP 429 or 503)."""
        return self.status_code in (429, 503)


def is_error_code(err: BaseException, code: str) -> bool:
    """Whether ``err`` is an API error with the given code."""
    if isinstance(err, HTTPError):
        return err.api_error is not None and err.api_error.code == code
    return isinstance(err, AllDebridError) and err.code == code
