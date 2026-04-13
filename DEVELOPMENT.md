# Development Guide

This document explains how to run the local stack and API tests.

## Prerequisites

- Docker Desktop (or Docker Engine + Docker Compose plugin)
- A shell with `docker` and `docker compose` available

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

## Run API Tests

Run the Dockerized Playwright API test suite using the same orchestration pattern as CI:

```bash
docker compose --profile tests up --build --abort-on-container-exit --exit-code-from api-tests api-tests
```

After tests finish (pass or fail), clean up containers and volumes:

```bash
docker compose --profile tests down -v
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
docker compose --profile tests up --build --abort-on-container-exit --exit-code-from api-tests api-tests
docker compose --profile tests down -v
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
