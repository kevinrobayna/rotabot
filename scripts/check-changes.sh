#!/usr/bin/env bash

set -eu -o pipefail

changes=$(git status --porcelain)

if [ -z "$changes" ]; then
  echo "PASS: No changes found."
else
  echo "FAIL: Uncommitted changes found:"
  echo "$changes"
  exit 1
fi