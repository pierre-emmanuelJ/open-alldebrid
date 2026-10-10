"""Generates src/open_alldebrid/_sync_client.py from _async_client.py.

Usage: python scripts/unasync.py [--check]
"""

from __future__ import annotations

import re
import sys
from pathlib import Path

SRC = Path(__file__).resolve().parent.parent / "src" / "open_alldebrid"

HEADER = '"""The sync client, generated from _async_client.py by scripts/unasync.py. DO NOT EDIT."""\n'

# Applied in order.
REPLACEMENTS = [
    ("import asyncio\nimport random\n", "import random\nimport time\n"),
    ("asyncio.sleep", "time.sleep"),
    ("httpx.AsyncClient", "httpx.Client"),
    ("AsyncAllDebrid", "AllDebrid"),
    ("AsyncLiveMagnets", "LiveMagnets"),
    ("async def __aenter__", "def __enter__"),
    ("async def __aexit__", "def __exit__"),
    ("aclose", "close"),
    ("An async client", "A client"),
    ("an async context manager", "a context manager"),
    ("async with", "with"),
    ("async def ", "def "),
    ("await ", ""),
]


def unasync(source: str) -> str:
    body = re.sub(r'\A""".*?"""\n', "", source, count=1, flags=re.S)
    for old, new in REPLACEMENTS:
        body = body.replace(old, new)
    return HEADER + body


def main() -> int:
    expected = unasync((SRC / "_async_client.py").read_text())
    target = SRC / "_sync_client.py"
    if "--check" in sys.argv:
        if not target.exists() or target.read_text() != expected:
            print("_sync_client.py is out of date: run python scripts/unasync.py", file=sys.stderr)
            return 1
        return 0
    target.write_text(expected)
    return 0


if __name__ == "__main__":
    sys.exit(main())
