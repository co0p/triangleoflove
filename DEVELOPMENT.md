# Development Guide

This document explains how to run the local stack and API tests.

## Canonical Task Entrypoint

For the supported repository-root development tasks in this repo, contributors and LLMs should prefer the root `Makefile` instead of reconstructing the underlying commands.

Supported targets:

| Task | Command |
| --- | --- |
| Show supported tasks | `make help` |
| Run backend tests | `make backend-test` |
| Run frontend tests | `make frontend-test` |
| Build the frontend bundle | `make frontend-build` |
| Run Docker acceptance tests | `make acceptance-test` |
| Build local Docker images | `make docker-build` |

Tasks outside that supported set, such as starting or stopping the local stack, still use the direct commands documented below.

## Prerequisites

- Docker Desktop (or Docker Engine + Docker Compose plugin)
- A shell with `docker` and `docker compose` available
- `make`
- Go toolchain for local backend tests
- Node.js 20+ for local frontend tests and builds

## Local Development (Low Effort)

This is the simplest local setup for working on frontend and backend with minimal tooling.

1. Start the full stack once:

```bash
docker compose up -d --build frontend backend db
```

2. Make code changes in `services/frontend` or `services/backend`.

3. Rebuild only the service you changed:

```bash
# after backend changes
docker compose up -d --build backend

# after frontend changes
docker compose up -d --build frontend
```

4. Watch service logs when needed:

```bash
docker compose logs -f backend frontend db
```

Notes:
- Hot reload/autoload is not enabled in this setup.
- This keeps local development simple and close to the containerized runtime used by CI.
- Backend database config uses a single `DATABASE_URL` variable in both local Docker and Railway.

## Local Environment Variables

For shell-based local runs (outside Docker), load the local variable contract:

```bash
source ./env-local.sh
```

This sets:
- `DATABASE_URL`
- `PORT`

## Start the Stack Locally

Build and start the application services:

```bash
docker compose up -d --build frontend backend db
```

The `db` service uses Railway's SSL-enabled Postgres base image (`ghcr.io/railwayapp-templates/postgres-ssl:latest`) so local and Railway behavior stay aligned.

The project uses ephemeral host ports for frontend and backend, so Docker chooses free ports automatically.

Find the assigned host ports:

```bash
docker compose port frontend 5173
docker compose port backend 8080
```

Example output:

- `0.0.0.0:52557` for frontend
- `0.0.0.0:52458` for backend

Use those ports in your browser or curl.

Verify backend endpoints:

```bash
BACKEND_PORT=$(docker compose port backend 8080 | awk -F: '{print $2}')
curl --fail "http://localhost:${BACKEND_PORT}/health"
curl --fail "http://localhost:${BACKEND_PORT}/status"
```

## Run Backend Tests

Use the canonical root target:

```bash
make backend-test
```

## Run Frontend Tests

Use the canonical root target:

```bash
make frontend-test
```

The target installs frontend dependencies into `services/frontend/node_modules` on first run via `npm ci`.

## Build Frontend

Use the canonical root target:

```bash
make frontend-build
```

## Run API Tests

Use the canonical root target:

```bash
make acceptance-test
```

Underlying command used by the make target:

```bash
docker compose --profile tests up --build --abort-on-container-exit --exit-code-from api-tests api-tests
```

The make target cleans up containers and volumes automatically after the run.

## Build Docker Images

Use the canonical root target:

```bash
make docker-build
```

Underlying command used by the make target:

```bash
docker compose --profile tests build frontend backend db api-tests
```

## Stop and Clean Up

Stop running services:

```bash
docker compose down
```

Stop services and remove volumes:

```bash
docker compose down -v
```

## Troubleshooting

If tests fail after changing dependencies, force image rebuild:

```bash
make acceptance-test
```

If services look stale, restart the stack:

```bash
docker compose down
docker compose up -d --build frontend backend db
```

## Playwright Testing Conventions

### `loginViaUI` — await navigation before leaving the dashboard

After calling `loginViaUI`, the token is stored asynchronously. If you navigate away from the dashboard immediately (e.g. `page.goto('/profile')`), the router auth guard may fire before the token is in `localStorage` and redirect back to `/login`.

Always add this assertion after `loginViaUI` before navigating away:

```js
await expect(page).toHaveURL(/\/dashboard/);
await page.goto('/profile');
```
