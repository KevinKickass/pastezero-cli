#!/bin/bash
set -e

VERSION=$(git describe --tags 2>/dev/null || echo "v0.0.0-dev")
RELEASE_DIR="release"

echo "-> GitHub Release erstellen: $VERSION"
gh release create "$VERSION" "$RELEASE_DIR"/* --title "$VERSION" --notes "Automatischer Build"
