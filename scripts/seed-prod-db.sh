#!/usr/bin/env bash
set -euo pipefail

# Usage: PROD_DATABASE_URL=<url> ./scripts/seed-prod-db.sh
#
# Drops all tables in the production database and re-seeds from services/db/init.sql.
# Requires psql. Install with: brew install libpq

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INIT_SQL="$SCRIPT_DIR/../services/db/init.sql"

if [ -z "${PROD_DATABASE_URL:-}" ]; then
  echo "Error: PROD_DATABASE_URL is not set."
  echo "Get it from Railway: db service → Variables → DATABASE_URL"
  exit 1
fi

if ! command -v psql &>/dev/null; then
  echo "Error: psql not found. Install with: brew install libpq"
  echo "Then add to PATH: export PATH=\"/opt/homebrew/opt/libpq/bin:\$PATH\""
  exit 1
fi

echo "WARNING: This will destroy all data in the production database."
read -r -p "Type 'yes' to continue: " confirm
if [ "$confirm" != "yes" ]; then
  echo "Aborted."
  exit 1
fi

echo "Dropping existing tables..."
psql "$PROD_DATABASE_URL" -c "DROP TABLE IF EXISTS checkins CASCADE;"
psql "$PROD_DATABASE_URL" -c "DROP TABLE IF EXISTS couples CASCADE;"
psql "$PROD_DATABASE_URL" -c "DROP TABLE IF EXISTS accounts CASCADE;"

echo "Applying init.sql..."
psql "$PROD_DATABASE_URL" -f "$INIT_SQL"

echo "Done. Production database seeded."
