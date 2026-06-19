# Roadmap

## Document Role

Tracks feature delivery against the product tiers defined in `PRD.md`. Updated at the end of each 4dc increment cycle during the promote phase. Status moves from **Partial** → **Done** when acceptance tests pass and the feature is verified on the deployed stack.

---

## Personal Tier

| Feature | Status | Acceptance Tests |
|---------|--------|-----------------|
| Registration | Done | [auth-backend.spec.js](../testing/tests/auth/auth-backend.spec.js) |
| Login / Logout | Done | [login.spec.js](../testing/tests/login/login.spec.js) |
| Dashboard | Done | [dashboard-matrix.spec.js](../testing/tests/dashboard/dashboard-matrix.spec.js) |
| Navbar | Done | [navbar.spec.js](../testing/tests/navbar/navbar.spec.js) |
| Profile (view + password change) | Done | [profile.spec.js](../testing/tests/profile/profile.spec.js) |
| Daily session (record ratings + note) | Done | [checkin.spec.js](../testing/tests/checkin/checkin.spec.js) |
| Weekly Insights matrix | Done | [insights.spec.js](../testing/tests/insights/insights.spec.js) |
| Users / me endpoint | Done | [users-me-backend.spec.js](../testing/tests/users/users-me-backend.spec.js) |
| Weekly session flow (backend and endpoint refactor) | Partial | — |
| Dashboard weekly backlog indicator and copy reframing | Done | [WeeklyBacklogDots.spec.js](../services/frontend/src/components/WeeklyBacklogDots.spec.js), [DashboardView.spec.js](../services/frontend/src/views/DashboardView.spec.js) |
| Impulses | Not Started | — |
| Private reflections | Not Started | — |

## Couple Tier

| Feature | Status | Acceptance Tests |
|---------|--------|-----------------|
| Pairing (invite code, connect, unpair) | Done | [pairing.spec.js](../testing/tests/pairing/pairing.spec.js) |
| Couple aggregate trend view | Not Started | — |
| Monthly shared session | Not Started | — |

---

## Infrastructure

| Item | Status | Notes |
|------|--------|-------|
| Docker Compose local stack | Done | `docker-compose.yml` at repo root |
| Railway deployment | Done | `railway.toml` per service |
| Playwright acceptance test harness | Done | `testing/` folder |
| Health checks (backend, db, frontend) | Done | `testing/tests/health/` |
| Admin account management | Done | `docs/api/admin.yml` |
| Rate limiting (registration + login) | Done | `docs/adr/2026-05-04-xff-client-ip-extraction.md` |
