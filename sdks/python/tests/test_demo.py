"""Calls the real API with the public demo keys. Set OFFLINE=1 to skip."""

from __future__ import annotations

import os
from collections.abc import Iterator

import pytest

from open_alldebrid import DEMO_KEYS, AllDebrid, AllDebridError, AsyncAllDebrid, DelayedStatus, is_error_code

pytestmark = [
    pytest.mark.demo,
    pytest.mark.skipif(os.environ.get("OFFLINE") == "1", reason="OFFLINE=1"),
]


@pytest.fixture
def client() -> Iterator[AllDebrid]:
    with AllDebrid(DEMO_KEYS["premium"], user_agent="open-alldebrid-python-tests") as c:
        yield c


@pytest.fixture
def anonymous() -> Iterator[AllDebrid]:
    with AllDebrid(user_agent="open-alldebrid-python-tests") as c:
        yield c


def test_public(anonymous: AllDebrid) -> None:
    anonymous.ping()
    assert anonymous.hosts()["hosts"]
    assert anonymous.host_domains()["hosts"]
    assert anonymous.host_priority()


def test_pin(anonymous: AllDebrid) -> None:
    assert anonymous.get_pin()["user_url"].startswith("https://")
    assert anonymous.check_pin("123B", "664c3ca2635c99f291d28e11ea18e154750bd21a")["activated"]


def test_user(client: AllDebrid, anonymous: AllDebrid) -> None:
    assert client.user()["username"] == "demoUserPremium"
    assert client.user_hosts()
    client.clear_notification("NOTIF_CODE")
    assert anonymous.verification_status("verificationTokenAllowed3f9eb0db")["verif"] == "allowed"
    with pytest.raises(AllDebridError) as caught:
        anonymous.resend_verification("verificationTokenAlreadySentb0db")
    assert caught.value.code == "ALREADY_SENT"


def test_links(client: AllDebrid) -> None:
    infos = client.link_infos(["https://example.com/file2MB.txt", "https://example.com/down"])
    assert infos[0].get("filename")
    assert client.extract_redirector("https://example.com/redir/folder")
    assert client.unlock("https://example.com/streaming").get("streaming")
    assert client.unlock("https://example.com/delayed").get("delayed")
    with pytest.raises(AllDebridError) as caught:
        client.unlock("https://example.com/down")
    assert is_error_code(caught.value, "LINK_DOWN")
    assert client.streaming_link("demov1bc797", "hd-720").get("link")
    assert client.delayed_link("1234567")["status"] == DelayedStatus.READY


def test_magnets(client: AllDebrid) -> None:
    uploaded = client.upload_magnets(["magnet:?xt=urn:btih:3A842783e3005495d5d1637f5364b59343c7844707"])
    assert uploaded[0].get("id")
    client.upload_torrent_files([("demo.torrent", b"d4:infod4:name4:demoee")])
    assert client.magnets()
    assert client.magnet(56789).get("files")
    assert len(client.magnet_files([56789])) == 1
    client.delete_magnet(56789)
    client.restart_magnet(99999)
    assert client.live_magnets().poll()


def test_saved_links(client: AllDebrid) -> None:
    assert client.saved_links()
    client.save_links(["https://example.com/file2MB.txt"])
    client.delete_saved_links(["https://example.com/file2MB.txt"])
    client.history()
    client.purge_history()


def test_bad_key(client: AllDebrid) -> None:
    with pytest.raises(AllDebridError) as caught:
        client.with_api_key("not-a-real-key").user()
    assert caught.value.code == "AUTH_BAD_APIKEY"


async def test_async() -> None:
    async with AsyncAllDebrid(DEMO_KEYS["premium"], user_agent="open-alldebrid-python-tests") as client:
        assert (await client.user())["username"] == "demoUserPremium"
        assert await client.magnets()
