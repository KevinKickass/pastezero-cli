#!/bin/bash
set -e

APP_NAME="pastezero"
BUILD_DIR="build"
RELEASE_DIR="release"
VERSION=$(git describe --tags 2>/dev/null || echo "v0.0.0-dev")
DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)

PLATFORMS=(
  "linux/amd64"
  "linux/arm64"
  "darwin/amd64"
  "darwin/arm64"
  "windows/amd64"
)

echo ">> E2E Tests"
#go test ./e2e -v

echo ">> Clean builds"
rm -rf "$BUILD_DIR" "$RELEASE_DIR"
mkdir -p "$BUILD_DIR" "$RELEASE_DIR"

for PLATFORM in "${PLATFORMS[@]}"
do
  IFS="/" read -r GOOS GOARCH <<< "$PLATFORM"
  EXT=""
  [ "$GOOS" == "windows" ] && EXT=".exe"

  BIN_NAME="${APP_NAME}-${GOOS}-${GOARCH}${EXT}"
  BIN_PATH="${BUILD_DIR}/${BIN_NAME}"

  echo "→ Building $BIN_NAME"
  env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 \
    go build -ldflags "-X main.version=${VERSION} -X main.buildTime=${DATE}" \
    -o "$BIN_PATH" main.go

  PKG_NAME="${APP_NAME}-${VERSION}-${GOOS}-${GOARCH}"
  PKG_PATH="${RELEASE_DIR}/${PKG_NAME}"

  mkdir -p "$PKG_PATH"
  cp "$BIN_PATH" "$PKG_PATH/$APP_NAME$EXT"

  # Archivieren
  if [ "$GOOS" == "windows" ]; then
    zip -j "${RELEASE_DIR}/${PKG_NAME}.zip" "$PKG_PATH/"*
  else
    tar -czf "${RELEASE_DIR}/${PKG_NAME}.tar.gz" -C "$PKG_PATH" .
  fi

  rm -rf "$PKG_PATH"
done

echo "[SUCCESS] Build + Packaging abgeschlossen: $RELEASE_DIR/"
