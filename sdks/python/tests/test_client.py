from __future__ import annotations

import json
import subprocess
import sys
from collections.abc import Callable
from pathlib import Path
from typing import Any
from urllib.parse import parse_qs

import httpx
import pytest
import yaml

from open_alldebrid import (
    ERROR_CODES,
    AllDebrid,
    AllDebridError,
    AsyncAllDebrid,
    HTTPError,
    is_error_code,
)

ROOT = Path(__file__).resolve().parent.parent
SPEC = yaml.safe_load((ROOT.parent.parent / "openapi.yaml").read_text())


class Recorder:
    """An httpx transport that records requests and answers canned bodies."""

    def __init__(self, *answers: tuple[int, str]) -> None:
        self.answers = list(answers)
        self.requests: list[httpx.Request] = []

    def __call__(self, request: httpx.Request) -> httpx.Response:
        request.read()
        self.requests.append(request)
        status, body = self.answers.pop(0) if self.answers else (200, '{"status":"success","data":{}}')
        return httpx.Response(status, text=body, headers={"Content-Type": "application/json"})

    @property
    def last(self) -> httpx.Request:
        return self.requests[-1]

    def form(self, request: httpx.Request | None = None) -> dict[str, list[str]]:
        return parse_qs((request or self.last).content.decode(), keep_blank_values=True)


def ok(data: Any) -> tuple[int, str]:
    return 200, json.dumps({"status": "success", "data": data})


def api_error(code: str, **extra: str) -> tuple[int, str]:
    return 200, json.dumps({"status": "error", "error": {"code": code, "message": "m", **extra}})


def make(*answers: tuple[int, str]) -> tuple[AllDebrid, Recorder]:
    rec = Recorder(*answers)
    http = httpx.Client(transport=httpx.MockTransport(rec))
    return AllDebrid("test-key", base_url="https://api.test", user_agent="unit", http_client=http), rec


def make_async(*answers: tuple[int, str]) -> tuple[AsyncAllDebrid, Recorder]:
    rec = Recorder(*answers)
    http = httpx.AsyncClient(transport=httpx.MockTransport(rec))
    return AsyncAllDebrid("test-key", base_url="https://api.test", user_agent="unit", http_client=http), rec


# ------------------------------------------------------------------ requests


def test_list_parameters_are_php_style() -> None:
    client, rec = make(ok({"magnets": [{"magnet": "a", "id": 1}]}))
    result = client.upload_magnets(["a", "b"])
    assert rec.last.method == "POST"
    assert rec.last.url.path == "/v4/magnet/upload"
    assert rec.form() == {"magnets[]": ["a", "b"]}
    assert rec.last.headers["Authorization"] == "Bearer test-key"
    assert rec.last.headers["User-Agent"] == "unit"
    assert result[0].get("id") == 1


def test_no_key_no_authorization() -> None:
    client, rec = make(ok({"ping": "pong"}))
    client.with_api_key(None).ping()
    assert "Authorization" not in rec.last.headers


def test_get_parameters_in_query() -> None:
    client, rec = make(ok({"hosts": {}}))
    client.hosts(hosts_only=True)
    assert rec.last.method == "GET"
    assert rec.last.url.params["hostsOnly"] == "true"


def test_password_only_when_given() -> None:
    client, rec = make()
    client.unlock("L")
    assert "password" not in rec.form()
    client.unlock("L", password="pw")
    assert rec.form()["password"] == ["pw"]


def test_torrent_files_are_multipart() -> None:
    client, rec = make(ok({"files": [{"magnet": "m", "id": 9}]}))
    client.upload_torrent_files([("a.torrent", b"d4:infoe"), ("b.torrent", b"x")])
    body = rec.last.content
    assert rec.last.headers["Content-Type"].startswith("multipart/form-data")
    assert body.count(b'name="files[]"') == 2
    assert b'filename="a.torrent"' in body and b'filename="b.torrent"' in body


# ------------------------------------------------------------------- routing

OPERATIONS = {f"{method.upper()} {path}" for path, item in SPEC["paths"].items() for method in item}

CALLS: list[tuple[str, Callable[[AllDebrid], Any], dict[str, list[str]], Any]] = [
    ("GET /v4/ping", lambda c: c.ping(), {}, {}),
    ("GET /v4.1/pin/get", lambda c: c.get_pin(), {}, {}),
    ("POST /v4/pin/check", lambda c: c.check_pin("P", "C"), {"pin": ["P"], "check": ["C"]}, {}),
    ("GET /v4/hosts", lambda c: c.hosts(), {}, {}),
    ("GET /v4/hosts/domains", lambda c: c.host_domains(), {}, {}),
    ("GET /v4/hosts/priority", lambda c: c.host_priority(), {}, {"hosts": {}}),
    ("GET /v4/user", lambda c: c.user(), {}, {"user": {"limitedHostersQuotas": []}}),
    ("GET /v4.1/user/hosts", lambda c: c.user_hosts(), {}, {"hosts": {}}),
    ("POST /v4/user/verif", lambda c: c.verification_status("T"), {"token": ["T"]}, {}),
    ("POST /v4/user/verif/resend", lambda c: c.resend_verification("T"), {"token": ["T"]}, {}),
    ("POST /v4/user/notification/clear", lambda c: c.clear_notification("N"), {"code": ["N"]}, {}),
    (
        "POST /v4/link/infos",
        lambda c: c.link_infos(["a", "b"], password="p"),
        {"link[]": ["a", "b"], "password": ["p"]},
        {"infos": []},
    ),
    ("POST /v4/link/redirector", lambda c: c.extract_redirector("L"), {"link": ["L"]}, {"links": []}),
    ("POST /v4/link/unlock", lambda c: c.unlock("L"), {"link": ["L"]}, {}),
    ("POST /v4/link/streaming", lambda c: c.streaming_link("G", "S"), {"id": ["G"], "stream": ["S"]}, {}),
    ("POST /v4/link/delayed", lambda c: c.delayed_link(42), {"id": ["42"]}, {}),
    ("POST /v4/magnet/upload", lambda c: c.upload_magnets(["M"]), {"magnets[]": ["M"]}, {"magnets": []}),
    ("POST /v4/magnet/upload/file", lambda c: c.upload_torrent_files([("a", b"x")]), {}, {"files": []}),
    ("POST /v4.1/magnet/status", lambda c: c.magnets(), {}, {"magnets": []}),
    ("POST /v4/magnet/files", lambda c: c.magnet_files([1, 2]), {"id[]": ["1", "2"]}, {"magnets": []}),
    ("POST /v4/magnet/delete", lambda c: c.delete_magnet(5), {"id": ["5"]}, {}),
    ("POST /v4/magnet/restart", lambda c: c.restart_magnets([5, 6]), {"ids[]": ["5", "6"]}, {}),
    ("GET /v4/user/links", lambda c: c.saved_links(), {}, {"links": []}),
    ("POST /v4/user/links/save", lambda c: c.save_links(["L"]), {"links[]": ["L"]}, {}),
    ("POST /v4/user/links/delete", lambda c: c.delete_saved_links(["L"]), {"links[]": ["L"]}, {}),
    ("GET /v4/user/history", lambda c: c.history(), {}, {"links": []}),
    ("POST /v4/user/history/delete", lambda c: c.purge_history(), {}, {}),
    ("GET /v4/voucher/balance", lambda c: c.voucher_balance(), {}, {"balance": 1.5}),
    ("POST /v4/voucher/get", lambda c: c.vouchers(30, 2), {"duration": ["30"], "nb": ["2"]}, {}),
    ("POST /v4/voucher/generate", lambda c: c.generate_vouchers(90, 1), {"duration": ["90"], "nb": ["1"]}, {}),
]


@pytest.mark.parametrize(("want", "call", "form", "data"), CALLS, ids=[c[0] for c in CALLS])
def test_routing(want: str, call: Callable[[AllDebrid], Any], form: dict[str, list[str]], data: Any) -> None:
    client, rec = make(ok(data))
    call(client)
    got = f"{rec.last.method} {rec.last.url.path}"
    assert got == want
    assert got in OPERATIONS
    if rec.last.method == "POST" and not rec.last.headers.get("Content-Type", "").startswith("multipart"):
        assert rec.form() == form


def test_every_operation_is_routed() -> None:
    assert {c[0] for c in CALLS} == OPERATIONS


def test_error_codes_match_spec() -> None:
    assert sorted(ERROR_CODES) == sorted(SPEC["components"]["schemas"]["ErrorCode"]["enum"])


def test_sync_client_is_up_to_date() -> None:
    result = subprocess.run(
        [sys.executable, str(ROOT / "scripts" / "unasync.py"), "--check"], capture_output=True, text=True
    )
    assert result.returncode == 0, result.stderr


def test_sync_and_async_have_the_same_methods() -> None:
    def public(cls: type) -> set[str]:
        return {"close" if n == "aclose" else n for n in vars(cls) if not n.startswith("_")}

    assert public(AllDebrid) == public(AsyncAllDebrid)


# ----------------------------------------------------------------- responses


def test_api_errors_carry_their_code() -> None:
    client, _ = make(api_error("MAGNET_INVALID_ID"))
    with pytest.raises(AllDebridError) as caught:
        client.delete_magnet(1)
    assert caught.value.code == "MAGNET_INVALID_ID"
    assert is_error_code(caught.value, "MAGNET_INVALID_ID")
    assert not is_error_code(caught.value, "LINK_DOWN")


def test_blocked_login_keeps_token() -> None:
    client, _ = make(api_error("AUTH_BLOCKED", token="tok"))
    with pytest.raises(AllDebridError) as caught:
        client.user()
    assert caught.value.token == "tok"


def test_http_errors() -> None:
    client, _ = make((429, "Too Many Requests"), (404, api_error("404")[1]))
    with pytest.raises(HTTPError) as limited:
        client.ping()
    assert limited.value.rate_limited
    with pytest.raises(HTTPError) as missing:
        client.ping()
    assert is_error_code(missing.value, "404")


def test_non_json_body() -> None:
    client, _ = make((200, "<html>"))
    with pytest.raises(ValueError, match="not a JSON object"):
        client.ping()


def test_single_magnet_object() -> None:
    client, _ = make(ok({"magnets": {"id": 7, "statusCode": 4, "files": [{"n": "a", "s": 1, "l": "L"}]}}))
    magnet = client.magnet(7)
    assert magnet["id"] == 7
    files = magnet.get("files") or []
    assert files[0].get("l") == "L"


def test_magnets_always_a_list() -> None:
    client, _ = make(ok({"magnets": {"id": 1}}))
    assert client.magnets() == [{"id": 1}]


def test_empty_quotas_become_a_dict() -> None:
    client, _ = make(ok({"user": {"username": "u", "limitedHostersQuotas": []}}))
    assert client.user()["limitedHostersQuotas"] == {}


# ------------------------------------------------------------------- helpers


def test_wait_for_pin() -> None:
    client, rec = make(
        ok({"activated": False, "expires_in": 600}), ok({"activated": True, "expires_in": 5, "apikey": "K"})
    )
    pin = {"pin": "ABCD", "check": "chk", "expires_in": 600, "user_url": "u", "base_url": "b"}
    assert client.wait_for_pin(pin, interval=0) == "K"  # type: ignore[arg-type]
    assert rec.form()["pin"] == ["ABCD"]


def test_wait_delayed_link() -> None:
    client, _ = make(ok({"status": 1, "time_left": 5}), ok({"status": 2, "time_left": 0, "link": "L"}))
    assert client.wait_delayed_link("1", interval=0) == "L"
    failed, _ = make(ok({"status": 3, "time_left": 0}))
    with pytest.raises(RuntimeError, match="failed"):
        failed.wait_delayed_link("1", interval=0)


def test_live_magnets() -> None:
    client, rec = make(
        ok(
            {
                "fullsync": True,
                "counter": 1,
                "magnets": [
                    {"id": 1, "filename": "a", "uploadDate": 10, "downloaded": 5},
                    {"id": 2, "filename": "b", "uploadDate": 20},
                ],
            }
        ),
        ok({"counter": 2, "magnets": {"id": 1, "downloaded": 9}}),
        ok({"counter": 3, "magnets": [{"id": 2, "deleted": True}, {"id": 3, "filename": "c", "uploadDate": 30}]}),
        ok({"counter": 4, "magnets": []}),
        ok({"fullsync": True, "counter": 1, "magnets": [{"id": 9, "uploadDate": 1}]}),
    )
    live = client.live_magnets()
    assert [m["id"] for m in live.poll()] == [2, 1]
    assert rec.form()["counter"] == ["0"]
    merged = live.poll()
    assert rec.form()["counter"] == ["1"]
    assert rec.form(rec.requests[0])["session"] == rec.form()["session"]
    assert next(m for m in merged if m["id"] == 1) == {"id": 1, "filename": "a", "uploadDate": 10, "downloaded": 9}
    assert [m["id"] for m in live.poll()] == [3, 1]
    assert [m["id"] for m in live.poll()] == [3, 1]
    assert [m["id"] for m in live.poll()] == [9]


async def test_async_client() -> None:
    client, rec = make_async(ok({"magnets": [{"magnet": "a", "id": 1}]}), api_error("LINK_DOWN"))
    async with client:
        result = await client.upload_magnets(["a"])
        assert result[0].get("id") == 1
        assert rec.form() == {"magnets[]": ["a"]}
        with pytest.raises(AllDebridError):
            await client.unlock("L")


async def test_async_live_magnets() -> None:
    client, _ = make_async(
        ok({"fullsync": True, "counter": 1, "magnets": [{"id": 1}]}), ok({"counter": 2, "magnets": []})
    )
    live = client.live_magnets()
    assert [m["id"] for m in await live.poll()] == [1]
    assert [m["id"] for m in await live.poll()] == [1]
