# Deployment Guide

This document describes the production deployment model for Triangle of Love using GitHub Actions and Railway.

## Overview

Production deployments follow a GHCR-backed Railway redeploy model:

1. GitHub Actions runs tests and builds Docker images on every push to `main`.
2. Backend and frontend images are published to the GitHub Container Registry (GHCR).
3. The deployment workflow triggers Railway to redeploy the image-backed services.

Railway services must be configured to pull images from GHCR. GitHub Actions is the source of truth for build, test, and publish.

## Published Images

On every push to `main`, the CI pipeline publishes two images to GHCR:

| Service  | Image                                                   |
| -------- | ------------------------------------------------------- |
| Backend  | `ghcr.io/co0p/triangleoflove-backend:latest`            |
| Backend  | `ghcr.io/co0p/triangleoflove-backend:<commit-sha>`      |
| Frontend | `ghcr.io/co0p/triangleoflove-frontend:latest`           |
| Frontend | `ghcr.io/co0p/triangleoflove-frontend:<commit-sha>`     |

Both images carry OCI metadata labels including `org.opencontainers.image.source` linking back to this repository.

Pull requests run tests and build images but do not publish to GHCR.

## CI/CD Pipeline

The pipeline is orchestrated by `.github/workflows/ci-orchestrator.yml`:

```
push to main
  ├── frontend-ci.yml   → tests, build, publish frontend image
  ├── backend-ci.yml    → tests, format check, publish backend image
  └── api-tests-ci.yml  → integration tests against published GHCR images
        └── deployment.yml → railway redeploy (backend + frontend)
```

Deployment only runs after all tests pass. PRs run the same tests but skip publishing and deployment.

## Railway Setup

### Required service configuration

Configure each Railway service as an **image-backed service** pointing at the GHCR `latest` tag:

- **Backend service**: `ghcr.io/co0p/triangleoflove-backend:latest`
- **Frontend service**: `ghcr.io/co0p/triangleoflove-frontend:latest`

Do not use Railway's source-based deployment for these services. The GitHub Actions pipeline controls what is built and when it is deployed.

### GHCR image visibility

If the GHCR packages are private, Railway needs credentials to pull them:

1. Create a GitHub Personal Access Token (classic) with `read:packages` scope.
2. In your Railway service image settings, provide:
   - Registry: `ghcr.io`
   - Username: your GitHub username (e.g. `co0p`)
   - Password: the PAT created above

Making the GHCR packages public simplifies setup by removing the credential requirement. You can change package visibility in the GitHub repository under **Settings → Packages**.

### Database

Use Railway's managed **Postgres** plugin for production. Do not deploy the repo's `services/db` image to Railway. The `services/db/` directory contains initialization SQL and a local dev Dockerfile that uses `ghcr.io/railwayapp-templates/postgres-ssl:latest`; on Railway, the managed Postgres service provides the same SSL-enabled Postgres without a custom image.

## Required Secrets and Variables

Configure the following in your GitHub repository under **Settings → Secrets and variables → Actions**:

### Secrets

| Name            | Description                                   |
| --------------- | --------------------------------------------- |
| `RAILWAY_TOKEN` | Railway API token with deploy access          |

### Variables

| Name                   | Description                                                        | Default        |
| ---------------------- | ------------------------------------------------------------------ | -------------- |
| `RAILWAY_PROJECT_ID`   | Railway project ID (visible in your Railway project URL or settings) | _(required)_   |
| `RAILWAY_ENVIRONMENT`  | Railway environment name to deploy to                              | `production`   |

### Railway environment variables

Set these in each Railway service's **Variables** tab:

**Backend:**

| Variable       | Description                                             |
| -------------- | ------------------------------------------------------- |
| `DATABASE_URL` | Postgres connection string from Railway managed Postgres |
| `JWT_SECRET`   | Secret key for signing JWT tokens                       |
| `CORS_ORIGINS` | Allowed CORS origins (e.g. `https://your-frontend.up.railway.app`) |
| `APP_BASE_URL` | Public base URL of the backend service                  |
| `PORT`         | Port the backend listens on (Railway sets this automatically) |

**Frontend:**

| Variable            | Description                                     |
| ------------------- | ----------------------------------------------- |
| `VITE_API_BASE_URL` | Public URL of the deployed backend API service  |

## Rollback

To roll back to a previous release, redeploy a commit-specific image tag. Every push to `main` publishes a `<commit-sha>`-tagged image alongside `latest`. You can manually trigger a Railway redeploy after updating the Railway service image tag to the desired commit SHA.
