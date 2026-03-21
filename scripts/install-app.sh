#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
DEFAULT_TARGET="/Applications"
FALLBACK_TARGET="$HOME/Applications"
APP_NAME="Focusbar.app"

if [[ $# -gt 0 ]]; then
  TARGET_DIR="$1"
elif [[ -w "$DEFAULT_TARGET" ]]; then
  TARGET_DIR="$DEFAULT_TARGET"
else
  TARGET_DIR="$FALLBACK_TARGET"
fi

mkdir -p "$TARGET_DIR"

"$ROOT_DIR/scripts/build-app-bundle.sh"

rm -rf "$TARGET_DIR/$APP_NAME"
ditto "$ROOT_DIR/dist/$APP_NAME" "$TARGET_DIR/$APP_NAME"

echo "Installed Focusbar to $TARGET_DIR/$APP_NAME"
echo "Launch with: open \"$TARGET_DIR/$APP_NAME\""
