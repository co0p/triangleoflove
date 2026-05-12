# DESIGN — Technical Architecture & Implementation Suggestions

## Document Role

This document defines technical implementation decisions that satisfy the product goals in `PRD.md`.

---

## 1. High-Level Architecture

### Components
1. **Web app (PWA)** — Vue 3 + Vite
2. **API service** — Go (REST JSON)
3. **Database** — Postgres (Railway managed)
4. **(Optional later)** Background worker — for scheduled nudges, email, etc.

### Key design goals
- Mobile-first, fast UX
- Private-by-default reflections
- Simple deployment and iteration
- LLM-friendly: mainstream tools, clear conventions, predictable code

---

## 2. Frontend (Vue PWA)

### Recommended stack
- **Vue 3 + TypeScript**
- **Vite**
- **vite-plugin-pwa**
- **Custom CSS library** (`services/frontend/src/assets/library.css`) — design tokens via CSS custom properties, no external framework. See `docs/adr/2026-03-30-custom-css-no-framework.md`.
- State:
  - Start with **Pinia** (or Vue composables) if needed
- Networking:
  - `fetch` wrapper + typed API client (small, explicit)
- Charts:
  - MVP: simple sparklines / minimal SVG charts
  - Optional: a Vue chart library later
### CSS Library

The visual system is a single file — `services/frontend/src/assets/library.css` — imported
globally in `main.js`. It follows a three-layer structure:

1. **Raw palette tokens** — named colour stops on `:root` (e.g. `--color-sage-500`)
2. **Semantic tokens** — purpose-driven aliases that components reference (e.g. `--color-primary`)
3. **Component classes** — selectors that reference only semantic tokens via `var()`; never raw values. A class used in two or more components must be extracted here from the component's `<style scoped>` block.

A theme is a named block of semantic token overrides applied to a container element. Components
never reference a theme directly — they inherit overrides through the CSS cascade.

The developer reference is `docs/design-system.html` — a standalone HTML file that references `library.css` directly and can be opened in any browser without running the app. It is updated in the same commit as any library change.

### Score-band presentation pattern

Feature-specific state colors (for example score band cues) follow the same three-layer discipline:

1. Add raw palette stops to `:root` only when no existing stop is suitable.
2. Add semantic tokens per state (for example `--color-insight-very-low`) that reference palette stops via `var()`.
3. Add state modifier classes in the component `<style scoped>` block that reference semantic tokens only.

Views derive score-band classes from numeric values. The API returns numeric scores only; band thresholds and color mapping remain frontend concerns.

### Per-dimension color token pattern

When multiple named dimensions each need their own color family, extend the
score-band pattern with dimension-scoped semantic tokens:

1. Add raw palette stops per family to `:root` only when no existing stop is
   suitable (e.g. `--color-sage-100`, `--color-rose-100`, `--color-gold-100`
   added for the very-low shade of each dimension family).
2. Add semantic tokens scoped to the dimension and band:
   `--color-{dimension}-{band}` (e.g. `--color-intimacy-very-low`).
   These reference palette stops via `var()`.
3. In the component, derive a CSS class per cell that encodes both dimension
   and band (e.g. `intimacy-very-low`), and bind one class per cell in the
   template. Component `<style scoped>` maps each class to its semantic token.

This keeps color concerns fully in CSS, allows per-dimension theming without
JavaScript, and is tree-shaken at build time per component.

*Emerged during CheckinMatrix (check-in history grid), May 2026.*

### Component conventions

- **NavBar** (`services/frontend/src/components/NavBar.vue`) imports `logo.svg` directly
  rather than accepting it as a prop. There is one logo asset and no story for swapping it;
  a prop would be an unused abstraction. Revisit if the logo becomes configurable.
### PWA requirements
- Offline-first behavior for check-ins:
  - Store pending submissions in **IndexedDB** (or localForage)
  - Sync queue when online
- Add-to-home-screen friendly
- Lightweight caching strategy:
  - cache shell + latest dashboard data
  - avoid caching private reflection text beyond what’s necessary

---

## 3. Backend (Go API)

### Go choices
- Router: **chi** (clean net/http style)
- JSON: standard `encoding/json`
- Validation: keep simple initially (manual validation or small validator lib)
- Auth:
  - Email + password
  - Password hashing: **bcrypt** or **argon2id**
  - JWT:
    - access token (short-lived)
    - refresh token (longer-lived, stored securely)

### REST API design
- Versioned routes: `/api/v1/...`
- Strict separation of:
  - **private user data** (reflections)
  - **shared aggregate couple data**
- Idempotency:
  - daily check-in should be upserted by (user_id, date)
  - weekly review upserted by (user_id, week_start_date)

### Layering note: shared domain sentinel errors

- Shared sentinel errors that must be imported by both repositories and services live in `internal/domain`.
- Example: `ErrDuplicateEmail` belongs in `domain` rather than `service` to avoid a circular import (`repository` must not import `service`).
- This keeps layering compliant with `repository -> domain` and `service -> domain`.

### Suggested endpoint list (MVP)
- Auth:
  - `POST /api/v1/auth/register`
  - `POST /api/v1/auth/login`
  - `POST /api/v1/auth/refresh`
  - `POST /api/v1/auth/logout`
- Couple pairing:
  - `POST /api/v1/couple/invite`
  - `POST /api/v1/couple/join`
  - `POST /api/v1/couple/leave`
- Check-ins:
  - `PUT /api/v1/checkins/daily` (upsert)
  - `PUT /api/v1/checkins/weekly` (upsert)
  - `GET /api/v1/insights/me` (personal trends)
  - `GET /api/v1/insights/couple` (aggregate-only)
- Impulses:
  - `GET /api/v1/impulses/weekly` (generated suggestion)
  - `POST /api/v1/impulses/plan` (select impulse)
- Achievements:
  - `GET /api/v1/achievements`
  - `POST /api/v1/achievements/claim`
  - `POST /api/v1/achievements/request-confirmation` (optional)

---

## 4. Data Model (Suggested Tables)

### Core tables
- `users`
- `couples`
- `couple_members` (user_id, couple_id, role, joined_at)
- `daily_checkins`
  - user_id, couple_id, date
  - q1..q5 (smallints 1–5)
  - optional_note (private)
- `weekly_reviews`
  - user_id, couple_id, week_start_date
  - intimacy, passion, commitment (1–5)
  - nice_text (private)
  - annoying_text (private)
- `impulses_catalog`
- `impulse_deliveries` (what was suggested when)
- `impulse_plans` (what the user selected)
- `achievements_catalog`
- `achievement_progress`
- `achievement_confirmations` (optional)

### Privacy rule enforcement
- Reflections (`nice_text`, `annoying_text`, notes) must never be returned via couple endpoints.
- Couple insights endpoints should compute aggregates without returning user-level raw text.

---

## 5. Impulse Engine (MVP: Rule-Based)
Start rule-based for reliability and safety.

Inputs:
- last weekly triangle ratings
- recent daily trend dips (e.g., 7-day average)
- chosen focus

Output:
- primary impulse + optional bonus
- difficulty level (easy/medium)
- estimated time cost (5/10/20 min)

Implementation suggestion:
- Store impulses as content cards with tags:
  - `dimension=intimacy|passion|commitment|communication`
  - `effort=low|medium`
  - `format=talk|plan|surprise|ritual`
- Pick top candidate based on:
  - lowest dimension OR biggest downward trend
  - avoid repeating the same impulse too frequently (cooldown)

---

## 6. Achievements System Design
### Achievement types
- **Secret (default):** user can track and claim without partner involvement.
- **Partner-confirmed (optional):** user triggers a confirmation request.

Key UX requirement:
- Partner should not automatically learn what the other is “working on.”

Backend notes:
- Confirmation requests should be generic unless user chooses details:
  - e.g., “Can you confirm we did the thing we talked about?” with optional label.

---

## 7. Railway Deployment (Recommended)
### Services
1. **web**: Vue app build + serve
   - Option A: Node service that serves `dist/`
   - Option B: Railway static hosting (if used in your setup)
2. **api**: Go service
3. **postgres**: Railway Postgres plugin

### Environment variables
**API**
- `DATABASE_URL`
- `JWT_SECRET`
- `CORS_ORIGINS` (comma-separated)
- `APP_BASE_URL`

**Web**
- `VITE_API_BASE_URL`

### CI/CD
- GitHub → Railway auto-deploy per branch (recommended)
- Use migrations on deploy (careful):
  - either run migrations in the API startup (simple MVP)
  - or a separate migration job step (cleaner later)
- CI/CD pipeline separation pattern:
  - `ci-orchestrator.yml` owns test/build orchestration and deployment trigger conditions.
  - `deployment.yml` owns deployment execution only.
  - Orchestration hands off to deployment only after CI/API gates pass, keeping release operations isolated and easier to operate.

---

## 8. LLM-Friendly Repo Conventions (Strongly Recommended)
- Monorepo structure:
  - `/web` (Vue)
  - `/api` (Go)
- Consistent formatting:
  - web: ESLint + Prettier
  - api: gofmt + golangci-lint
- Typed contracts:
  - Maintain a shared `openapi.yaml` (or a simple `/contracts` JSON schema set)
  - Generate TS types (optional but helpful)

---

## 9. Security & Compliance Notes
- Use HTTPS only (Railway provides)
- Store passwords hashed (bcrypt/argon2id)
- Implement rate limits on auth endpoints
- Ensure data deletion is feasible (hard delete or anonymize)
- Add clear disclaimers: not therapy; if users are in crisis, seek professional help

---

## Emergent Architecture

> Patterns and structures that emerged through TDD. Updated after each increment. Not a planning document.

### Handler Struct with Per-Route Methods

- **What**: Each handler is a struct (`PairingHandler`) with one method per route — `GetCode`, `Regenerate`, `Connect`, `GetCoupleStatus`. Methods are registered individually on the mux in `main.go` via `http.HandlerFunc(ph.Method)`.
- **Why it emerged**: A single handler switching on `r.URL.Path` made each route hard to test in isolation and mixed routing concerns into handler logic. Registering discrete methods on the mux keeps routing in `main.go` and makes each endpoint independently testable.
- **Where used**: `services/backend/internal/web/pairing_handler.go`

### Repository Split by Aggregate

- **What**: When a handler serves two distinct aggregates, each aggregate gets its own repository type and its own service interface. `PairingRepository` owns `invite_code` operations on the `accounts` table; `CoupleRepository` owns the `couples` table. The service takes two named interfaces (`InviteCodeRepo`, `CoupleRepo`) rather than one combined interface.
- **Why it emerged**: A single `PairingRepo` interface with six methods made the service's dependencies opaque and its mocks unwieldy. Splitting by aggregate clarifies ownership and keeps each mock minimal.
- **Where used**: `services/backend/internal/repository/pairing_repository.go`, `couple_repository.go`, `services/backend/internal/service/pairing_service.go`

### Page-Level Layout Utility Class

- **What**: A `.page` class in `library.css` provides `min-height: 100vh` and `background-color: var(--color-bg)` for full-height views. Views set `class="page"` on their root element.
- **Why it emerged**: Each full-height view was defining a scoped per-view class with identical rules. Extracting to `library.css` removes the duplication and follows the existing convention that classes used in two or more components belong in the library.
- **Where used**: `services/frontend/src/assets/library.css`, `PairingView.vue`, `DashboardView.vue`

### Seed Data Helper Module

- **What**: Seed user data lives in `testing/tests/helpers/users.js` as a `USERS` array of objects (`{ email, password, firstName }`). Other helpers derive their constants from this single source.
- **Why it emerged**: Email and password strings were duplicated across `auth.js` and test files. Centralising into `users.js` gives tests a single place to update when seed data changes and makes the full user list available for multi-user flows (e.g. pairing connect).
- **Where used**: `testing/tests/helpers/users.js`, `testing/tests/helpers/auth.js`

### Dual-Layer Registration Credential Policy

- **What**: Registration evaluates the same credential rules twice: immediate signals and submit-time guards in the registration view, then final enforcement in the backend auth service. The auth API client preserves backend validation messages instead of replacing them with a generic failure.
- **Why it emerged**: The increment required real-time guidance during form entry and also guaranteed rejection if invalid data still reached account creation. Keeping the same rule set on both sides reduced drift between UI guidance and backend enforcement.
- **Where used**: `services/frontend/src/views/RegisterView.vue`, `services/frontend/src/api/auth.js`, `services/backend/internal/service/auth_service.go`, `services/backend/internal/web/registration_handler.go`

| Date | Increment | Changes |
|------|-----------|---------|
| 2026-04-09 | Couple Pairing | Handler struct pattern, repository split by aggregate, `.page` layout utility, seed data helper module |
| 2026-04-12 | Apply the Test Pyramid | Vue component test patterns for router-using views, mock reset discipline |
| 2026-04-13 | Unpair from couple | Destructive-action confirmation modal pattern |
| 2026-04-27 | Weekly Insights Matrix | RouterLink stub `v-bind="$attrs"` for `to` prop assertability in component tests |
| 2026-05-11 | Registration Flow Clarity and Validation | Dual-layer registration credential policy and backend validation-message passthrough at the auth API boundary |

### Destructive-Action Confirmation Modal

- **What**: A full-viewport overlay (`.modal-overlay`) containing a card (`.modal-card`) with a question and two action buttons (`.modal-actions`). Used before any irreversible user action.
- **Why it emerged**: The Unpair action is unilateral and permanent; TDD of `PairingView.vue` drove a dedicated confirmation step before the API call.
- **Where used**: `services/frontend/src/views/PairingView.vue`; classes in `services/frontend/src/assets/library.css`.

### Vue Component Tests — Router and Mock Reset

- **What**: Views that call `useRouter()` require `vi.mock('vue-router', ...)` at the module level and a `RouterLink` stub via `global.stubs`. `vi.clearAllMocks()` clears call history but does **not** reset `mockResolvedValue` — default return values must always be re-applied in `beforeEach`. When a test needs to assert on the `to` prop of a `RouterLink`, use `{ template: '<a v-bind="$attrs"><slot />' }` so the prop is forwarded to the rendered anchor and assertable via `wrapper.find('a').attributes('href')` or equivalent.
- **Why it emerged**: `DashboardView.spec.js` — mock state from the `GivenPartnerName` test leaked into the `GivenNoPartner` test when relying on `clearAllMocks()` alone. The `v-bind="$attrs"` addition came from the Weekly Insights increment when asserting that the insights entry link pointed to `/insights`.
- **Where used**: `services/frontend/src/views/DashboardView.spec.js`
- **Pattern**:
  ```js
  vi.mock('vue-router', () => ({
    useRouter: vi.fn().mockReturnValue({ push: vi.fn() }),
  }));

  const stubs = {
    NavBar: true,
    RouterLink: { template: '<a v-bind="$attrs"><slot /></a>' },
  };

  beforeEach(() => {
    vi.clearAllMocks();
    someApi.fn.mockResolvedValue(defaultValue); // explicit reset required
  });

  mount(MyView, { global: { stubs } });
  ```

---
