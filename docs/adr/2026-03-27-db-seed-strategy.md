# ADR: Database Seed Strategy

## Context

The application requires a pre-seeded account to exist at startup for the
single-user login increment. A mechanism is needed to create the schema and
seed data when the database container starts.

## Decision

Bake an `init.sql` file into the `services/db` Dockerfile by copying it to
`/docker-entrypoint-initdb.d/`. The Postgres container executes scripts in that
directory automatically on first initialisation. The seed account credentials
are stored with a bcrypt-hashed password.

## Consequences

- **Benefits:** Zero external tooling; works identically in local development,
  CI, and any fresh container start.
- **Drawbacks:** Schema changes require rebuilding the image; no incremental
  migration history.
- **Trade-offs:** Acceptable while the schema is small and a single seed account
  is sufficient. When schema evolution or multiple environments are needed,
  a migration tool (e.g. golang-migrate) should replace this approach.

## Alternatives Considered

- **Migration tool (golang-migrate)**: Enables incremental schema changes and
  rollbacks, but adds tooling and an initial migration setup cost not warranted
  for a single-table schema.
- **Manual seed via application startup**: Would run on every start, requiring
  idempotency guards; more complex than a one-time init script.
