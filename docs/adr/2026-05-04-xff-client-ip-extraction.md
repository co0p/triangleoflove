# ADR: Client IP Extraction via X-Forwarded-For for Rate Limiting

**Date:** 2026-05-04
**Status:** Accepted

## Context

The rate-limiting middleware introduced in the User Registration increment must
identify the originating client IP so that each real end-user gets an
independent request quota. The backend runs behind Railway's edge proxy, which
terminates TLS and forwards all traffic to the Go server over an internal
network. As a result, `r.RemoteAddr` in every inbound request reflects the
proxy's internal IP, not the user's IP — making it useless as a per-client
discriminator.

Railway's edge proxy, like most reverse proxies (nginx, AWS ALB, Cloudflare),
appends the real client IP as the leftmost entry in the `X-Forwarded-For`
header before the request reaches the origin. The header format is:

```
X-Forwarded-For: <client>, <proxy1>, <proxy2>, ...
```

## Decision

The rate-limiting middleware reads the **leftmost value** of the
`X-Forwarded-For` header as the canonical client IP. If the header is absent
(e.g. in direct connections during local development or unit tests), it falls
back to `r.RemoteAddr` with the port stripped.

Implementation (`internal/web/ratelimit.go`):

```go
func clientIP(r *http.Request) string {
    if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
        parts := strings.SplitN(xff, ",", 2)
        return strings.TrimSpace(parts[0])
    }
    addr := r.RemoteAddr
    if i := strings.LastIndex(addr, ":"); i != -1 {
        return addr[:i]
    }
    return addr
}
```

## Rationale

**Why the leftmost XFF entry?**
The leftmost entry is the IP the proxy received the connection from — i.e., the
original client. Rightmost entries are added by intermediate proxies and reflect
infrastructure addresses, not users.

**Why not `r.RemoteAddr`?**
In the Railway deployment topology, every request arrives from the internal
proxy IP. Using `RemoteAddr` would put all users in a single rate-limit bucket,
making the middleware useless.

**Why not a custom header?**
Railway sets `X-Forwarded-For` automatically. A custom header would require
proxy configuration changes and would not survive a platform change. The `XFF`
header is the de-facto standard understood by all major proxies and CDNs.

**Security consideration — header spoofing**
A client can send a forged `X-Forwarded-For: 1.2.3.4` header if they connect
directly (bypassing the proxy). In the Railway topology, direct access to the
origin port is blocked by the platform; all production traffic must pass through
the edge proxy, which appends the real IP. If the deployment topology changes
(e.g. multi-hop proxies), the chain should be validated and the rightmost
trusted-proxy IP extracted instead. This is noted as a backlog item.

## Consequences

- Rate limiting is effective per real client IP in the Railway production
  environment.
- Local development (`docker compose`) has no proxy layer; `r.RemoteAddr` is
  used as the fallback, which works correctly in that context.
- Unit tests set `X-Forwarded-For` directly to control the per-IP buckets
  without needing a real proxy.
- If the deployment topology adds a second hop (e.g. a CDN in front of
  Railway), the extraction logic must be revisited to walk the XFF chain from
  the rightmost trusted entry.
