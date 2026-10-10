"""The sync client, generated from _async_client.py by scripts/unasync.py. DO NOT EDIT."""

from __future__ import annotations

import random
import time
from collections.abc import Iterable
from typing import IO, Any, cast

import httpx

from ._errors import AllDebridError, HTTPError
from .models import (
    DelayedData,
    GeneratedVouchersData,
    HostDomainsData,
    HostsData,
    LinkInfo,
    Magnet,
    MagnetFiles,
    PinCheckData,
    PinGetData,
    RestartedMagnet,
    SavedLink,
    Service,
    StreamingData,
    UnlockData,
    UploadedMagnet,
    User,
    VerificationData,
    VouchersData,
)

DEFAULT_BASE_URL = "https://api.alldebrid.com"
DEFAULT_USER_AGENT = "open-alldebrid-python"

Params = list[tuple[str, str]]
TorrentFile = tuple[str, bytes | IO[bytes]]
"""A .torrent file to upload: (file name, content)."""

__all__ = ["AllDebrid", "LiveMagnets", "TorrentFile"]


class AllDebrid:
    """A client for the AllDebrid API.

    Use it as a context manager, or call ``close()`` when done::

        with AllDebrid(api_key, user_agent="my-app") as client:
            user = client.user()
    """

    def __init__(
        self,
        api_key: str | None = None,
        *,
        base_url: str = DEFAULT_BASE_URL,
        user_agent: str = DEFAULT_USER_AGENT,
        timeout: float = 30.0,
        http_client: httpx.Client | None = None,
    ) -> None:
        self._api_key = api_key
        self._base_url = base_url.rstrip("/")
        self._user_agent = user_agent
        self._owns_http = http_client is None
        self._http = http_client or httpx.Client(timeout=timeout)

    def __enter__(self) -> AllDebrid:
        return self

    def __exit__(self, *exc: object) -> None:
        self.close()

    def close(self) -> None:
        """Closes the HTTP client, unless it was passed in."""
        if self._owns_http:
            self._http.close()

    def with_api_key(self, api_key: str | None) -> AllDebrid:
        """A client with another API key, sharing this one's HTTP client."""
        clone = AllDebrid(api_key, base_url=self._base_url, user_agent=self._user_agent, http_client=self._http)
        clone._owns_http = False
        return clone

    def _call(
        self,
        method: str,
        path: str,
        params: Params | None = None,
        files: list[tuple[str, tuple[str, bytes | IO[bytes], str]]] | None = None,
    ) -> Any:
        headers = {"Accept": "application/json", "User-Agent": self._user_agent}
        if self._api_key:
            headers["Authorization"] = f"Bearer {self._api_key}"
        url = self._base_url + path
        # List parameters keep their PHP-style names and repeat: magnets[]=a&magnets[]=b.
        if method == "GET":
            query = httpx.QueryParams(tuple(params)) if params else None
            response = self._http.request(method, url, params=query, headers=headers)
        elif files:
            response = self._http.request(method, url, files=files, headers=headers)
        else:
            response = self._http.request(method, url, data=dict_list(params), headers=headers)

        try:
            envelope = response.json()
        except ValueError:
            envelope = None
        if response.status_code != 200:
            raise HTTPError(response.status_code, api_error(envelope))
        if not isinstance(envelope, dict):
            raise ValueError(f"AllDebrid {method} {path}: the response is not a JSON object")
        if envelope.get("status") != "success":
            err = api_error(envelope)
            if err is None:
                raise ValueError(f"AllDebrid {method} {path}: status {envelope.get('status')!r} without error")
            raise err
        return envelope.get("data")

    # ------------------------------------------------------------------ general

    def ping(self) -> None:
        """Checks that the API is up."""
        self._call("GET", "/v4/ping")

    # ---------------------------------------------------------------------- pin

    def get_pin(self) -> PinGetData:
        """Starts a PIN authentication. Show ``pin`` or ``user_url`` to the user, then ``wait_for_pin``."""
        return cast(PinGetData, self._call("GET", "/v4.1/pin/get"))

    def check_pin(self, pin: str, check: str) -> PinCheckData:
        """The state of a PIN authentication; ``apikey`` is set once activated."""
        return cast(PinCheckData, self._call("POST", "/v4/pin/check", [("pin", pin), ("check", check)]))

    def wait_for_pin(self, pin: PinGetData, interval: float = 5.0) -> str:
        """Polls until the user confirms the PIN and returns the new API key."""
        while True:
            state = self.check_pin(pin["pin"], pin["check"])
            if state["activated"] and state.get("apikey"):
                return state["apikey"]
            time.sleep(interval)

    # -------------------------------------------------------------------- hosts

    def hosts(self, hosts_only: bool = False) -> HostsData:
        """Supported hosts, streaming sites and redirectors. Cache the result."""
        return cast(HostsData, self._call("GET", "/v4/hosts", [("hostsOnly", "true")] if hosts_only else None))

    def host_domains(self) -> HostDomainsData:
        """Every supported domain."""
        return cast(HostDomainsData, self._call("GET", "/v4/hosts/domains"))

    def host_priority(self) -> dict[str, int]:
        """Host names ranked from the most open (1) to the most restricted."""
        return cast(dict[str, int], (self._call("GET", "/v4/hosts/priority"))["hosts"])

    # --------------------------------------------------------------------- user

    def user(self) -> User:
        """The account information. Empty quotas are always a dict."""
        user = cast(User, (self._call("GET", "/v4/user"))["user"])
        if isinstance(user.get("limitedHostersQuotas"), list):
            user["limitedHostersQuotas"] = {}
        return user

    def user_hosts(self, hosts_only: bool = False) -> dict[str, Service]:
        """The hosts this user can use now, with their quotas."""
        data = self._call("GET", "/v4.1/user/hosts", [("hostsOnly", "true")] if hosts_only else None)
        return cast(dict[str, Service], data["hosts"])

    def verification_status(self, token: str) -> VerificationData:
        """Follows a new-location login verification (token from an AUTH_BLOCKED error)."""
        return cast(VerificationData, self._call("POST", "/v4/user/verif", [("token", token)]))

    def resend_verification(self, token: str) -> None:
        """Sends the verification email again (once per token)."""
        self._call("POST", "/v4/user/verif/resend", [("token", token)])

    def clear_notification(self, code: str) -> None:
        """Clears one of the codes in ``user()["notifications"]``."""
        self._call("POST", "/v4/user/notification/clear", [("code", code)])

    # -------------------------------------------------------------------- links

    def link_infos(self, links: Iterable[str], password: str | None = None) -> list[LinkInfo]:
        """File name and size of each link, in order; unreadable links have ``error`` set."""
        params: Params = [("link[]", link) for link in links]
        if password:
            params.append(("password", password))
        return cast(list[LinkInfo], (self._call("POST", "/v4/link/infos", params))["infos"])

    def extract_redirector(self, link: str) -> list[str]:
        """Links behind a redirector; they only work with ``link_infos`` and ``unlock``."""
        return cast(list[str], (self._call("POST", "/v4/link/redirector", [("link", link)]))["links"])

    def unlock(self, link: str, password: str | None = None) -> UnlockData:
        """Unlocks a hoster link.

        When ``delayed`` is set, wait with ``wait_delayed_link``; when
        ``streaming`` lists qualities, choose one with ``streaming_link``.
        """
        params: Params = [("link", link)]
        if password:
            params.append(("password", password))
        return cast(UnlockData, self._call("POST", "/v4/link/unlock", params))

    def streaming_link(self, id: str, stream: str) -> StreamingData:
        """The link of one quality: ``id`` from ``unlock``, ``stream`` from one of its ``streaming`` entries."""
        return cast(StreamingData, self._call("POST", "/v4/link/streaming", [("id", id), ("stream", stream)]))

    def delayed_link(self, id: str | int) -> DelayedData:
        """The status of a delayed link."""
        return cast(DelayedData, self._call("POST", "/v4/link/delayed", [("id", str(id))]))

    def wait_delayed_link(self, id: str | int, interval: float = 5.0) -> str:
        """Polls a delayed link until it is ready and returns the download link."""
        while True:
            state = self.delayed_link(id)
            if state["status"] == 2 and state.get("link"):
                return state["link"]
            if state["status"] == 3:
                raise RuntimeError(f"AllDebrid: delayed link {id} failed")
            time.sleep(interval)

    # ------------------------------------------------------------------ magnets

    def upload_magnets(self, magnets: Iterable[str]) -> list[UploadedMagnet]:
        """Adds magnet URIs or info hashes. Results are in order; invalid ones have ``error`` set."""
        params: Params = [("magnets[]", m) for m in magnets]
        return cast(list[UploadedMagnet], (self._call("POST", "/v4/magnet/upload", params))["magnets"])

    def upload_torrent_files(self, files: Iterable[TorrentFile]) -> list[UploadedMagnet]:
        """Adds .torrent files, given as (name, content). Invalid files have ``error`` set."""
        parts = [("files[]", (name, content, "application/x-bittorrent")) for name, content in files]
        return cast(list[UploadedMagnet], (self._call("POST", "/v4/magnet/upload/file", files=parts))["files"])

    def magnets(self) -> list[Magnet]:
        """Every magnet of the account."""
        return as_list((self._call("POST", "/v4.1/magnet/status"))["magnets"])

    def magnet(self, id: int) -> Magnet:
        """One magnet, with its ``files``."""
        found = as_list((self._call("POST", "/v4.1/magnet/status", [("id", str(id))]))["magnets"])
        if len(found) != 1:
            raise ValueError(f"AllDebrid: magnet {id}: got {len(found)} magnets")
        return found[0]

    def magnet_files(self, ids: Iterable[int]) -> list[MagnetFiles]:
        """The file tree of each magnet, in order; unknown IDs have ``error`` set."""
        params: Params = [("id[]", str(i)) for i in ids]
        return cast(list[MagnetFiles], (self._call("POST", "/v4/magnet/files", params))["magnets"])

    def delete_magnet(self, id: int) -> None:
        """Deletes a magnet."""
        self._call("POST", "/v4/magnet/delete", [("id", str(id))])

    def restart_magnet(self, id: int) -> None:
        """Restarts a failed magnet."""
        self._call("POST", "/v4/magnet/restart", [("id", str(id))])

    def restart_magnets(self, ids: Iterable[int]) -> list[RestartedMagnet]:
        """Restarts failed magnets; each result has its own ``error``."""
        params: Params = [("ids[]", str(i)) for i in ids]
        return cast(list[RestartedMagnet], (self._call("POST", "/v4/magnet/restart", params)).get("magnets", []))

    def live_magnets(self) -> LiveMagnets:
        """Follows the magnets with live mode: each ``poll`` only transfers changes."""
        return LiveMagnets(self)

    # -------------------------------------------------------------- saved links

    def saved_links(self) -> list[SavedLink]:
        """The links the user saved."""
        return cast(list[SavedLink], (self._call("GET", "/v4/user/links"))["links"])

    def save_links(self, links: Iterable[str]) -> None:
        """Saves links for later."""
        self._call("POST", "/v4/user/links/save", [("links[]", link) for link in links])

    def delete_saved_links(self, links: Iterable[str]) -> None:
        """Deletes saved links."""
        self._call("POST", "/v4/user/links/delete", [("links[]", link) for link in links])

    def history(self) -> list[SavedLink]:
        """Links unlocked in the last 3 days, when history logging is enabled."""
        return cast(list[SavedLink], (self._call("GET", "/v4/user/history"))["links"])

    def purge_history(self) -> None:
        """Deletes the whole recent links history."""
        self._call("POST", "/v4/user/history/delete")

    # ----------------------------------------------------------------- reseller

    def voucher_balance(self) -> float:
        """The reseller balance. Reseller accounts only."""
        return float((self._call("GET", "/v4/voucher/balance"))["balance"])

    def vouchers(self, duration_days: int, count: int) -> VouchersData:
        """Takes vouchers (15, 30, 90, 180 or 365 days) from the reseller stock."""
        params: Params = [("duration", str(duration_days)), ("nb", str(count))]
        return cast(VouchersData, self._call("POST", "/v4/voucher/get", params))

    def generate_vouchers(self, duration_days: int, count: int) -> GeneratedVouchersData:
        """Generates new vouchers from the reseller balance."""
        params: Params = [("duration", str(duration_days)), ("nb", str(count))]
        return cast(GeneratedVouchersData, self._call("POST", "/v4/voucher/generate", params))


class LiveMagnets:
    """Live mode state: the merged list of magnets, updated by ``poll``."""

    def __init__(self, client: AllDebrid) -> None:
        self._client = client
        self.session = random.randrange(1 << 31)
        self._counter = 0
        self._state: dict[int, Magnet] = {}

    def poll(self) -> list[Magnet]:
        """Applies the changes since the last call and returns every magnet, newest first."""
        params: Params = [("session", str(self.session)), ("counter", str(self._counter))]
        data = self._client._call("POST", "/v4.1/magnet/status", params)
        if data.get("fullsync"):
            self._state = {}
        for change in as_list(data.get("magnets")):
            if change.get("deleted"):
                self._state.pop(change["id"], None)
            else:
                merged = cast(Magnet, {**self._state.get(change["id"], {}), **change})
                self._state[change["id"]] = merged
        self._counter = int(data["counter"])
        return sorted(self._state.values(), key=lambda m: (m.get("uploadDate", 0), m["id"]), reverse=True)


def as_list(value: Any) -> list[Magnet]:
    """The API sends magnets as a list, or as one object for a single magnet."""
    if value is None:
        return []
    return cast(list[Magnet], value if isinstance(value, list) else [value])


def dict_list(params: Params | None) -> dict[str, list[str]] | None:
    """Groups repeated parameters, which httpx then sends once per value."""
    if not params:
        return None
    grouped: dict[str, list[str]] = {}
    for key, value in params:
        grouped.setdefault(key, []).append(value)
    return grouped


def api_error(envelope: Any) -> AllDebridError | None:
    if isinstance(envelope, dict) and isinstance(envelope.get("error"), dict):
        e = envelope["error"]
        return AllDebridError(str(e.get("code", "")), str(e.get("message", "")), e.get("token"))
    return None
