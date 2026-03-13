# Wireframes (pencil-sketch SVG)

This folder contains **low-fidelity, black & white, pencil-sketch style** mobile wireframes for the Triangle of Love Coach MVP.

## Screens
The SVG files are numbered in the order of the primary user journey.

## Export SVG → PNG
These wireframes are stored as SVG so they stay editable. To generate PNGs:

### Prerequisites
- Install **Inkscape** (CLI `inkscape` must be available in your PATH).

macOS (Homebrew):
```bash
brew install inkscape
```

Ubuntu/Debian:
```bash
sudo apt-get update
sudo apt-get install -y inkscape
```

### Export
Run:
```bash
bash docs/wireframes/export-wireframes-to-png.sh
```

Outputs will be written to:
- `docs/wireframes/png/*.png`

### Notes
- Export scale defaults to **2x** for crisper review images.
- The export script is deterministic; re-run anytime after editing SVGs.

## Conventions
- SVG artboard: **390×844** (iPhone 13-ish)
- `viewBox`: `0 0 390 844`
- Monochrome strokes only; keep it low-fi.