#!/bin/bash
set -euo pipefail

FILE="$(realpath "$1")"
DIR="$(dirname "$FILE")"

MOCK_DIR="$DIR/mocks"
mkdir -p "$MOCK_DIR"

MOD="$(cd "$DIR" && go env GOMOD)"
[ "$MOD" != "/dev/null" ] || { echo "go.mod not found (run inside a Go module)"; exit 1; }

MODDIR="$(dirname "$MOD")"
IMPORT_PATH="$(cd "$DIR" && go list -f '{{.ImportPath}}' .)"

mapfile -t INTERFACES < <(
  grep -E 'type\s+[A-Za-z_][A-Za-z0-9_]*\s+interface\s*\{' "$FILE" | awk '{print $2}'
)

for I in "${INTERFACES[@]}"; do
  OUT="$MOCK_DIR/mock_$(echo "$I" | tr '[:upper:]' '[:lower:]').go"

  if [ -f "$OUT" ]; then
    echo "Removing existing $OUT"
    rm "$OUT"
  fi

  echo "Generating $OUT (only $I)..."

  (
    cd "$MODDIR"
    mockgen \
      -destination="$OUT" \
      -package="mocks" \
      "$IMPORT_PATH" \
      "$I"
  )
done
