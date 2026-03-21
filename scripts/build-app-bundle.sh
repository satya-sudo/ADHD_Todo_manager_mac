#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
BUNDLE_ROOT="$ROOT_DIR/dist/Focusbar.app"
CONTENTS_DIR="$BUNDLE_ROOT/Contents"
MACOS_DIR="$CONTENTS_DIR/MacOS"

rm -rf "$BUNDLE_ROOT"
mkdir -p "$MACOS_DIR"

(
  cd "$ROOT_DIR"
  go build -o "$MACOS_DIR/focusbar" ./cmd/app
)
cp "$ROOT_DIR/macos/Info.plist" "$CONTENTS_DIR/Info.plist"

chmod +x "$MACOS_DIR/focusbar"
codesign --force --deep --sign - "$BUNDLE_ROOT"

echo "Built app bundle at $BUNDLE_ROOT"
echo "Signed app bundle at $BUNDLE_ROOT"
echo "Run with: open \"$BUNDLE_ROOT\""
