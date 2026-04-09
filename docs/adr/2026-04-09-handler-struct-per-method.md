# Handler Struct with Per-Route Methods

**Date**: 2026-04-09  
**Status**: Accepted

## Context

The initial pairing handler was implemented as a single function switching on `r.URL.Path` and `r.Method`. This worked for the first route but became unwieldy as four routes were added (`GET /api/v1/pairing`, `POST /api/v1/pairing/regenerate`, `POST /api/v1/pairing/connect`, `GET /api/v1/couples/me`). The single-function approach:

- Mixed routing logic into handler code
- Made it impossible to test individual routes without constructing a correctly-pathed request
- Forced all routes to share the same dependency injection point

## Decision

Each handler is a struct with one exported method per route. The struct is constructed once in `main.go` with its dependencies injected. Each method is registered individually on the mux via `http.HandlerFunc(ph.Method)`.

```go
type PairingHandler struct { svc PairingService }

func (h *PairingHandler) GetCode(w http.ResponseWriter, r *http.Request)    { ... }
func (h *PairingHandler) Regenerate(w http.ResponseWriter, r *http.Request) { ... }
```

Routing responsibility lives entirely in `main.go`:

```go
mux.Handle("GET /api/v1/pairing",             mid(http.HandlerFunc(ph.GetCode)))
mux.Handle("POST /api/v1/pairing/regenerate", mid(http.HandlerFunc(ph.Regenerate)))
```

## Consequences

- Each route is independently unit-testable by calling the method directly with a test `ResponseRecorder`
- Dependency injection is explicit at construction time via the struct
- Adding a new route is a two-step operation: add a method, register it in `main.go`
- The mux pattern string (e.g. `"GET /api/v1/pairing"`) is the authoritative record of method + path — no branching inside the handler
