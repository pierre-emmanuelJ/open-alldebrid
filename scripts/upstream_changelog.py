"""Detects changes to the official AllDebrid API docs.

Extracts the changelog of https://docs.alldebrid.com and compares it with
docs/upstream-changelog.txt, the version this spec was last checked against.

    python scripts/upstream_changelog.py           # exit 1 and print a diff when it changed
    python scripts/upstream_changelog.py --update  # record the current changelog
"""

from __future__ import annotations

import difflib
import re
import sys
import urllib.request
from html.parser import HTMLParser
from pathlib import Path

DOCS_URL = "https://docs.alldebrid.com/"
SNAPSHOT = Path(__file__).resolve().parent.parent / "docs" / "upstream-changelog.txt"


class _Text(HTMLParser):
    """Collects the text after the element with id="changelog"."""

    def __init__(self) -> None:
        super().__init__()
        self.inside = False
        self.skip = 0
        self.lines: list[str] = []
        self.current: list[str] = []

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        if dict(attrs).get("id") == "changelog":
            self.inside = True
        if tag in ("script", "style"):
            self.skip += 1
        if tag in ("li", "p", "br", "h1", "h2", "h3"):
            self._flush()

    def handle_endtag(self, tag: str) -> None:
        if tag in ("script", "style"):
            self.skip -= 1
        if tag in ("li", "p", "h1", "h2", "h3"):
            self._flush()

    def handle_data(self, data: str) -> None:
        if self.inside and not self.skip:
            self.current.append(data)

    def _flush(self) -> None:
        text = re.sub(r"\s+", " ", "".join(self.current)).strip()
        if text:
            self.lines.append(text)
        self.current = []


def fetch_changelog() -> str:
    request = urllib.request.Request(DOCS_URL, headers={"User-Agent": "open-alldebrid-drift-check"})
    with urllib.request.urlopen(request, timeout=30) as response:
        html = response.read().decode("utf-8", "replace")
    parser = _Text()
    parser.feed(html)
    parser._flush()
    if not parser.lines:
        raise SystemExit("no changelog found in the docs page")
    return "\n".join(parser.lines) + "\n"


def main() -> int:
    current = fetch_changelog()
    if "--update" in sys.argv:
        SNAPSHOT.write_text(current)
        return 0
    recorded = SNAPSHOT.read_text() if SNAPSHOT.exists() else ""
    if current == recorded:
        print("The upstream changelog has not changed.")
        return 0
    sys.stdout.writelines(
        difflib.unified_diff(recorded.splitlines(True), current.splitlines(True), "recorded", "docs.alldebrid.com")
    )
    return 1


if __name__ == "__main__":
    sys.exit(main())
