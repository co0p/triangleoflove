# ADR: Railway as Deployment Target for All Services

## Context

The project deploys three runtime services: frontend, backend, and database. A consistent deployment target is needed to avoid split-configuration drift and to keep CI/delivery pipelines simple.

## Decision

Railway is the canonical deployment target for all production services in this project:

- Each service (`services/frontend`, `services/backend`, `services/db`) has a `railway.toml` configuration.
- Merge to `main` triggers automatic Railway deployment for all services.
- Railway manages the production Postgres instance for the database service.
- Environment variables are configured through Railway's environment system, not committed to the repository.

## Consequences

- **Benefits:** Single platform for all services simplifies observability, secrets management, and rollback. No multi-platform coordination overhead.
- **Drawbacks:** Platform lock-in; migrating away would require re-implementing deployment config for all services.
- **Trade-offs:** Accept platform dependency in exchange for low operational overhead during early product development.

## Alternatives Considered

- Split platforms (e.g., Vercel for frontend, fly.io for backend): more flexibility per service but increases deployment coordination complexity.
- Self-hosted (e.g., VPS + Docker): full control but significantly higher operational burden.
