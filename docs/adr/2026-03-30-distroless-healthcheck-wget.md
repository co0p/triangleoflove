# ADR: Distroless Healthcheck via Busybox wget

## Context

The backend runs on `gcr.io/distroless/static-debian12`, which contains no shell, no package manager, and no utility binaries. Docker's `HEALTHCHECK` with `CMD-SHELL` is therefore unavailable. The image needs a probe binary capable of making an HTTP request to `/api/v1/health`.

## Decision

Copy `/bin/wget` from `busybox:musl` into the distroless final stage as a dedicated third build stage. The `HEALTHCHECK` uses exec form: `CMD ["/wget", "-q", "-O", "-", "http://localhost:8080/api/v1/health"]`.

## Consequences

- **Benefits:** Minimal binary footprint; no shell attack surface added; probe works inside distroless without switching to a debug image.
- **Drawbacks:** Adds a third build stage (`busybox:musl`); wget binary is not managed by the distroless base image's update cycle.
- **Trade-offs:** Accepted for the operational benefit of a working healthcheck without compromising the distroless security posture.

## Alternatives Considered

- **`distroless/static-debian12:debug`**: Includes busybox but is not recommended for production; adds shell and other tools.
- **Custom Go healthcheck binary**: Correct but heavier; adds compile step and binary to maintain.
- **Switch base image to `alpine`**: wget available natively, but loses the distroless security guarantees.
