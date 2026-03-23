# CONSTITUTION.md

## Document Role

This document defines binding engineering and delivery constraints for the repository.

## Architectural Decisions

### Layering
- Default: Service code lives under `services/` with one subfolder per runtime service (for example `services/frontend`, `services/backend`, and `services/db`).
- Default: Each part has its own `Dockerfile`.
- Default: Local development uses one root `docker-compose` setup to run all parts together.
- Default: Frontend UI components call an API client layer only. Direct HTTP calls from UI components are not allowed.
- Default: Backend flow is `handlers/routes -> service/use-case -> repository/data`.
- Default: The `services/db` folder contains migrations and seed scripts only.
- Exceptions: Layering can be bypassed only in short-lived spike branches.
- Enforcement signal: PR review blocks feature code that crosses these boundaries.

### Error Handling
- Default: Backend services return typed domain errors.
- Default: HTTP handlers translate domain errors into HTTP responses.
- Default: Unexpected errors return a generic HTTP 500 response with no internal details.
- Exceptions: None for production paths.
- Enforcement signal: Handler tests must verify error-to-status translation.

### State Management
- Default: Frontend server state is loaded through the API client and stored in a shared app state layer.
- Default: Component-local state is for local presentation and interaction only.
- Exceptions: Local-only state is allowed for trivial single-screen behavior with no cross-screen reuse.
- Enforcement signal: Shared business state must not be duplicated across multiple UI components.

### Dependencies
- Default: External systems are accessed through backend adapters/repositories, not directly from handlers or domain logic.
- Default: Frontend uses typed API contracts/client wrappers, not ad-hoc request and response shapes.
- Exceptions: Direct dependency calls are allowed only inside dedicated adapter modules.
- Enforcement signal: Service and domain layers must not import transport or driver packages directly.

## Testing Expectations

- Test location: Unit tests are colocated with source. API acceptance tests live in `testing/`.
- Coverage: Critical backend API flows require acceptance tests. Core business logic paths require unit tests. UI end-to-end tests are out of scope.
- Runtime: Acceptance tests must run in CI and must also run locally against the `docker-compose` stack.
- Mocking: Unit tests may mock ports/adapters. Acceptance tests should hit real API endpoints with real service wiring and controlled test data setup/teardown.

## Artifact Layout

- **CONSTITUTION.md**: Project root
- **DESIGN.md**: Project root (emergent architecture, updated after increments)
- **ADRs**: `docs/adr/` using `YYYY-MM-DD-short-title.md` naming (unnumbered)
- **API contracts**: `docs/api/` with OpenAPI as the canonical contract when present
- **Other docs**: `docs/`
- **Working context**: `.4dc/current/` (temporary, gitignored)

## Delivery Practices

- PR size: Flexible. Keep PRs reviewable and focused on one coherent change.
- CI requirements: Full gate before merge to `main`: build, lint, unit tests, and API acceptance tests.
- Branching: Use individual branches per feature/change and merge into `main`.
- Deployment: Merge to `main` triggers pipeline and automatic production deployment.
- Deployment artifact policy (temporary): Use latest-style Railway deployment flow while delivery speed is the priority.
- Deployment artifact policy reason: Optimize for moving fast and fixing fast; exact-SHA rollout remains planned follow-up work.