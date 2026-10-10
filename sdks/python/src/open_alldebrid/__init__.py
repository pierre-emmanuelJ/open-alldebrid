"""Typed client for the AllDebrid API, built from the open-alldebrid OpenAPI spec.

    from open_alldebrid import AllDebrid

    with AllDebrid(api_key, user_agent="my-app") as client:
        user = client.user()

``AsyncAllDebrid`` has the same methods, as coroutines.
"""

from ._async_client import AsyncAllDebrid, AsyncLiveMagnets
from ._errors import ERROR_CODES, AllDebridError, HTTPError, is_error_code
from ._sync_client import DEFAULT_BASE_URL, AllDebrid, LiveMagnets, TorrentFile

__version__ = "1.0.0"

DEMO_KEYS = {
    "premium": "staticDemoApikeyPrem",
    "trial": "staticDemoApikeyTria",
    "free": "staticDemoApikeyFree",
}
"""Public demo keys: the API answers them with static sample data."""


class MagnetStatus:
    """Magnet status codes. 5 and above are errors."""

    QUEUED = 0
    DOWNLOADING = 1
    COMPRESSING = 2
    UPLOADING = 3
    READY = 4


class DelayedStatus:
    """Delayed link statuses."""

    PROCESSING = 1
    READY = 2
    FAILED = 3


__all__ = [
    "DEFAULT_BASE_URL",
    "DEMO_KEYS",
    "ERROR_CODES",
    "AllDebrid",
    "AllDebridError",
    "AsyncAllDebrid",
    "AsyncLiveMagnets",
    "DelayedStatus",
    "HTTPError",
    "LiveMagnets",
    "MagnetStatus",
    "TorrentFile",
    "is_error_code",
]
