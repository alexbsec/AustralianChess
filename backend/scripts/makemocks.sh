#!/bin/bash
set -euo pipefail

FILE="$(realpath "$1")"
DIR="$(dirname "$FILE")"
PKG="$(grep -E '^\s*package\s+' "$FILE" | head -n 1 | awk '{print $2}')"

mapfile -t INTERFACES < <(grep -E 'type\s+[A-Za-z_][A-Za-z0-9_]*\s+interface\s*\{' "$FILE" | awk '{print $2}')

# find module root and import path
MOD="$(cd "$DIR" && go env GOMOD)"
[ "$MOD" != "/dev/null" ] || { echo "go.mod not found (run inside a Go module)"; exit 1; }
MODDIR="$(dirname "$MOD")"
IMPORT_PATH="$(cd "$DIR" && go list -f '{{.ImportPath}}' .)"

for I in "${INTERFACES[@]}"; do
  OUT="$DIR/mock_$(echo "$I" | tr '[:upper:]' '[:lower:]').go"
  if [ -f "$OUT" ]; then
    echo "Removing existing $OUT"
    rm "$OUT"
  fi
  echo "Generating $OUT (only $I)..."
  ( cd "$MODDIR" && mockgen -destination="$OUT" -package="$PKG" "$IMPORT_PATH" "$I" )
done
