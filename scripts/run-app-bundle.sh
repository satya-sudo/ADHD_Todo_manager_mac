#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
BUNDLE_PATH="$ROOT_DIR/dist/Focusbar.app"
LOG_PATH="${FOCUSBAR_LOG_PATH:-$HOME/Library/Logs/Focusbar/focusbar.log}"

"$ROOT_DIR/scripts/build-app-bundle.sh"

echo "Launching $BUNDLE_PATH"
echo "Logs: $LOG_PATH"

open "$BUNDLE_PATH"
