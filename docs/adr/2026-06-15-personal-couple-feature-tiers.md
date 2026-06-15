# ADR: Personal / Couple Feature Tiers

Date: 2026-06-15

## Status

Accepted

## Context

The app targets individuals in relationships. A user may sign up before their partner does, or use the app without ever inviting them. The original PRD assumed pairing was required for core functionality, which made solo use feel like a broken demo state.

## Decision

All product features are assigned to exactly one of two tiers:

**Personal tier** — available immediately to any authenticated user, no active pairing required:
- Weekly session (5 questions → triangle scores → reflections → impulse)
- Personal dashboard with health state and trend
- Impulses
- Private reflections

**Couple tier** — available only when the user has an active pairing:
- Monthly shared session
- Couple aggregate trend view
- Partner invite management

Unpaired users see Couple-tier features with a single, calm prompt ("Invite your partner to unlock this"). No empty states, no pressure.

## Consequences

- Solo use is a first-class, complete product state.
- Auth and data access logic must check pairing status before serving any Couple-tier endpoint.
- UI components that render Couple-tier features must handle the unpaired state gracefully — a locked/invite state, not an error or empty state.
- The Personal tier must never degrade in usefulness as Couple features are added.
