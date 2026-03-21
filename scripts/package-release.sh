#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
DIST_DIR="$ROOT_DIR/dist"
APP_NAME="Focusbar.app"
ARCHIVE_PATH="$DIST_DIR/Focusbar-macOS.zip"

"$ROOT_DIR/scripts/build-app-bundle.sh"

rm -f "$ARCHIVE_PATH"
ditto -c -k --keepParent "$DIST_DIR/$APP_NAME" "$ARCHIVE_PATH"

echo "Packaged release archive at $ARCHIVE_PATH"
