# ADR: GHCR-Backed API Test Guardrail Before Production Deploy

## Context

The delivery pipeline builds service images and publishes release artifacts. To reduce deployment risk, acceptance tests should validate deployable artifacts, not only source-built local containers.

The repository now uses top-level `testing/` for API acceptance tests.

## Decision

We use API acceptance tests from `testing/` as a production deployment guardrail in CI:

- Build and publish deployable service images.
- Run API acceptance tests against those release images in pipeline orchestration.
- Treat passing `testing/` acceptance tests as a required gate before production deployment progression.

## Consequences

- **Benefits:** Higher confidence that released artifacts are deployable and behaviorally correct.
- **Drawbacks:** Longer CI duration and more registry/network dependency.
- **Trade-offs:** Accept slower pipelines in exchange for stronger pre-deploy validation.

## Alternatives Considered

- Source-only acceptance tests as final gate: faster but weaker artifact-level assurance.
- Deploy first, test after deploy: faster release path but higher production risk.
