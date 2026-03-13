#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd ""+"$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OUT_DIR="$ROOT_DIR/png"

mkdir -p "$OUT_DIR"

if ! command -v inkscape >/dev/null 2>&1; then
  echo "Error: inkscape is required but was not found in PATH." >&2
  echo "Install Inkscape, then re-run." >&2
  exit 1
fi

# Export at 2x scale for crispness (96dpi is 1x; 192dpi is 2x)
EXPORT_DPI=192

shopt -s nullglob
SVGS=($("$ROOT_DIR"/*.svg))

if [[ "+"${#SVGS[@]}" -eq 0 ]]; then
  echo "No SVGs found in $ROOT_DIR" >&2
  exit 1
fi

for svg in "${SVGS[@]}"; do
  base="$(basename "$svg" .svg)"
  out="$OUT_DIR/$base.png"
  echo "Exporting $svg -> $out"

  inkscape "$svg" \
    --export-type=png \
    --export-filename="$out" \
    --export-dpi="$EXPORT_DPI" \
    --export-background=white

done

echo "Done. PNGs are in $OUT_DIR"