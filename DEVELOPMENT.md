# Development Guide

This document explains how to run the local stack and API tests.

## Prerequisites

- Docker Desktop (or Docker Engine + Docker Compose plugin)
- A shell with `docker` and `docker compose` available

## Start the Stack Locally

Build and start the application services:

```bash
docker compose up -d --build frontend backend db
```

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
