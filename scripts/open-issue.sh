#!/bin/sh
# Opens an issue with this title, or comments on the open one with the same
# title, so a weekly check does not pile up duplicates.
# Usage: scripts/open-issue.sh <title> <body-file>
set -eu
title=$1
body=$2
existing=$(gh issue list --state open --search "\"$title\" in:title" --json number,title \
  --jq ".[] | select(.title == \"$title\") | .number" | head -n 1)
if [ -n "$existing" ]; then
  gh issue comment "$existing" --body-file "$body"
else
  gh issue create --title "$title" --body-file "$body"
fi
