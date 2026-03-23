#!/usr/bin/env bash
set -euo pipefail

# Local runtime contract: match production variable names.
export DATABASE_URL="postgres://triangle:triangle@localhost:5432/triangleoflove?sslmode=disable"
export PORT="8080"

echo "Exported DATABASE_URL and PORT for local backend runtime."
