# CONSTITUTION.md

## Document Role

This document defines binding engineering and delivery constraints for the repository.

## Architectural Decisions

### Layering
- Default: Service code lives under `services/` with one subfolder per runtime service (for example `services/frontend`, `services/backend`, and `services/db`).
- Default: Each service has its own `Dockerfile`.
- Default: Local development uses one root `docker-compose` setup to run all services together.
- Default: Frontend UI components call an API client layer only. Direct HTTP calls from UI components are not allowed.
- Default: All shared CSS classes live in `library.css`. When a class appears in two or more components it must be extracted to `library.css`. Component `<style>` blocks contain only single-use layout rules. All color and spacing values in style rules must reference `var(--token)`; raw hex or rgba values are permitted only inside the `:root` block of `library.css`.
- Default: Backend flow is `handlers/routes -> service/use-case -> repository/data`.
- Default: All backend routes use the `/api/v1/` prefix.
- Default: The `services/db` folder contains migrations and seed scripts only.
- Exceptions: Layering can be bypassed only in short-lived spike branches.
- Enforcement signal: PR review blocks feature code that crosses these boundaries.

### Error Handling
- Default: Backend services return typed domain errors. All model types and the shared `ErrNotFound` sentinel live in `internal/domain`. Repository methods signal record absence by returning `domain.ErrNotFound`; they do not return per-entity error variants or boolean found-flags.
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
- Default: The frontend Vite dev server proxies `/api` to the backend via the `BACKEND_URL` environment variable (defaults to `http://backend:8080`). This is required because the frontend Dockerfile runs the Vite dev server, not a static build.
- Default: Repository method names follow Spring Data conventions — `FindBy<Field>` for queries, `Save` for create/update, `ExistsBy<Field>` for existence checks. No `Get*`, `Set*`, or `Create*` prefixes. See `docs/adr/2026-04-13-standardize-repository-pattern.md`.
- Exceptions: Direct dependency calls are allowed only inside dedicated adapter modules.
- Enforcement signal: Service and domain layers must not import transport or driver packages directly.

### Development Workflow
- Default: When the repository root `Makefile` exposes a supported local development task, contributors and LLMs should invoke the make target from the repository root instead of reconstructing the underlying command.
- Default: The `Makefile` is the canonical entrypoint for supported local development tasks only; stack lifecycle, deployment, and other unsupported workflows may continue to use direct commands until a make target exists for them.
- Default: Root make targets remain thin wrappers over the existing Go, npm, Docker Compose, and Docker workflows rather than redefining those workflows.
- Exceptions: If no root make target exists for a task, direct tool-specific commands are allowed.
- Enforcement signal: `DEVELOPMENT.md` and automation guidance reference supported make targets first.

## Testing Expectations

- Test location: Unit tests are colocated with source. API acceptance tests live in `testing/`. Frontend component tests use Vitest and @vue/test-utils, are colocated in `services/frontend/src/`, and run via `npm test` inside that directory. The `testing/` folder is for API and browser (Playwright) acceptance tests only.
- Coverage: Critical backend API flows require acceptance tests. Core business logic paths require unit tests. Browser-driven acceptance tests are used for critical auth and routing flows. Full cross-service user-journey E2E tests are out of scope.
- Runtime: Acceptance tests must run in CI and must also run locally against the `docker-compose` stack.
- Mocking: Unit tests may mock ports/adapters. Acceptance tests should hit real API endpoints with real service wiring and controlled test data setup/teardown.
- Browser-driven acceptance tests cover user-journey happy paths only. Edge cases (invalid inputs, error states, authorization boundaries) are covered by colocated Go unit tests in the service and handler layers, not by Playwright.
- Playwright acceptance tests must simulate real user flows through the UI (navigate, click, fill, assert visible output). They must not call API endpoints directly via `request`. Test data preconditions are established through the seed server, not by calling the API.
- Playwright spec files are grouped by feature in subdirectories under `testing/tests/` (e.g. `testing/tests/pairing/pairing.spec.js`). The file name provides feature context; test names use Given-When-Then format without a feature prefix.
- Vue component specs: When a component template uses `<router-link>`, the test must pass a `RouterLink` stub via `global.stubs` in the `mount` call. A `vi.mock('vue-router')` export does not auto-register the component globally and will produce a "Failed to resolve component: router-link" warning without the stub.

### Layer Preference Guide

Place each test at the **lowest layer** that can meaningfully cover it:

| Layer | Own | Example |
|-------|-----|---------|
| Go unit (service / handler) | Error paths, auth boundaries, domain logic, input validation | `TestCheckinHandler_GivenNoAuth_WhenGET_ThenReturns401` — no DB, uses `httptest` + mock repo |
| Vue component (Vitest) | View render states (empty, loaded, error) driven by mock API data | `TestDashboard_GivenNoPartner_WhenRendered_ThenNotConnectedVisible` — no backend, mocks `api/*.js` |
| API acceptance (Playwright, Docker stack) | Happy-path user journeys through the real UI with real backend wiring | `GivenValidCredentials_WhenLogin_ThenRedirectsToDashboard` — browser, real DB |

**Rule**: if a mock suffices, don't reach for Docker. If a unit suffices, don't reach for a component test. Service constructors accept repository interfaces (not concrete structs) to enable unit mocking without a database.

## Artifact Layout

- **CONSTITUTION.md**: Project root
- **docs/DESIGN.md**: `docs/` (emergent architecture, updated after increments)
- **ADRs**: `docs/adr/` using `YYYY-MM-DD-short-title.md` naming (unnumbered)
- **API contracts**: `docs/api/` with OpenAPI as the canonical contract when present
- **Other docs**: `docs/`

## Target Audience

- Default: The primary user is on a mobile phone. All UI decisions start from mobile constraints.
- Default: The 375 px viewport is the design baseline. Wider breakpoints are additive.
- Default: Interactive elements must meet a 44 px minimum tap target height.
- Default: Layouts must remain usable when the soft keyboard is open and reduces the visible
  viewport by up to 300 px.
- Default: Any PRD for this product must include a Mobile Interaction Model section that defines binding touch-first interaction constraints before any feature descriptions. At minimum it must specify: how required inputs are handled without a keyboard; and the 44 px tap target enforcement.
- Enforcement signal: Visual review at 375 px is required before marking any frontend deliverable Done.
- Enforcement signal: PRDs without a Mobile Interaction Model section are not ready for the plan phase.

## Delivery Practices

- PR size: One coherent feature change per PR. Prefer small PRs reviewable in a single focused pass.
- CI requirements: Full gate before merge to `main`: build, lint, unit tests, and API acceptance tests.
- Branching: Use individual branches per feature/change and merge into `main`.
- Deployment: Merge to `main` triggers Railway auto-deploy for all services. See `docs/adr/2026-03-27-railway-deployment-target.md`.