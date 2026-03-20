# ADR: Docker Runtime Conventions for Local and CI

## Context

The project needed a stable local and CI runtime baseline for frontend, backend, database, and API acceptance tests. During implementation, three patterns proved important:

- Health/status endpoints for explicit service verification.
- Ephemeral host port mapping to reduce developer port conflicts.
- Readiness checks before dependent services and tests execute.

These patterns improved reproducibility and reduced flaky startup/test behavior.

## Decision

We standardize the following Docker runtime conventions:

- Local and CI orchestration uses the root compose setup as the primary runtime contract.
- Services expose explicit health/status behavior for runtime verification.
- Host-facing service ports are mapped ephemerally in local usage to avoid machine-specific collisions.
- Database and dependent services use readiness checks/health-gated startup to avoid race conditions.
- API acceptance tests run only after stack readiness conditions are satisfied.

## Consequences

- **Benefits:** More reliable startup, fewer flaky acceptance runs, and lower onboarding friction.
- **Drawbacks:** Slightly more compose and startup configuration complexity.
- **Trade-offs:** Prefer deterministic runtime behavior over minimal configuration.

## Alternatives Considered

- Static host port mapping everywhere: simpler but causes frequent local collisions.
- Startup without readiness checks: simpler but leads to intermittent failures and race conditions.
