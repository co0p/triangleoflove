# DESIGN — Technical Architecture & Implementation Suggestions
**Date:** 2026-03-10  
**Stack choice:** Vue 3 PWA + Go API + Postgres, deployed on Railway

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
- **Tailwind CSS** (rapid mobile-first styling, easy iteration)
- State:
  - Start with **Pinia** (or Vue composables) if needed
- Networking:
  - `fetch` wrapper + typed API client (small, explicit)
- Charts:
  - MVP: simple sparklines / minimal SVG charts
  - Optional: a Vue chart library later

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

## 10. Next Implementation Deliverables (Suggested)
1. Repo skeleton (web/api) + Railway deploy config
2. Database schema + migrations
3. Auth + couple pairing endpoints
4. Daily + weekly check-in endpoints + UI
5. Insights aggregation endpoints + basic dashboard
6. Impulse engine v1 + achievements v1