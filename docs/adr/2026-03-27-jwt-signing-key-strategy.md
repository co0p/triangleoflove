# ADR: JWT Signing Key Strategy

## Context

The backend signs JWTs for authenticated sessions. A signing secret is required.
In the first increment, with a single pre-seeded account and no production deployment,
a hardcoded development constant was used to keep the implementation simple.
The constant is clearly marked as not for production use.

## Decision

For local development and CI, use a hardcoded constant (`devSecret`) in
`internal/auth/jwt.go`. Before any production deployment, replace this with a
required `JWT_SECRET` environment variable, sourced from a secret manager (e.g.
Railway's environment variable injection). The backend must refuse to start if
the variable is absent in non-development environments.

## Consequences

- **Benefits:** Zero infrastructure overhead in development; secret rotation is
  straightforward when the env var approach is in place.
- **Drawbacks:** The hardcoded constant is a security risk if accidentally
  deployed to production before the env var path is wired.
- **Trade-offs:** Acceptable for a single-user development environment; blocks
  production readiness until the env var path is implemented.

## Alternatives Considered

- **Required env var from day one**: Safer but adds friction to local setup before
  any secrets infrastructure exists.
- **Auto-generated secret per boot**: Simple but invalidates all sessions on every
  restart, which is unusable in development.
