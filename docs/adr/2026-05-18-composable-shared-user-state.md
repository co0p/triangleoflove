# ADR: Composable-Based Shared User State for Session Identity

## Context

The frontend needed one authenticated-user identity value, `firstName`, to appear consistently in both the NavBar and dashboard without duplicate `GET /api/v1/users/me` calls or greeting flicker during route changes.

The project already favors a shared app state layer for server state, but this increment only required a single session-scoped value and a small amount of loading coordination. Introducing a full store library for that narrow need would add setup and indirection without solving a more complex state problem.

## Decision

Use a Vue composable with module-scoped reactive state for authenticated user identity in the current browser session.

The composable owns:
- a shared `firstName` ref
- a singleton in-flight `load()` promise to deduplicate concurrent fetches
- token-aware cache invalidation when the session changes
- a `reset()` method for logout and test isolation

Components that need the authenticated user's display name read it through this composable rather than fetching profile data independently.

## Consequences

- **Benefits:** Removes duplicate profile fetches, prevents greeting flicker across authenticated route transitions, and keeps the solution smaller than a full store.
- **Drawbacks:** The pattern is limited to the current runtime session and relies on module scope, which is less explicit than a dedicated store when state needs grow.
- **Trade-offs:** The team accepts a lightweight composable now and will revisit a store library only if cross-feature shared state expands materially.

## Alternatives Considered

- **Pinia store**: Rejected for this increment because one shared identity value did not justify store setup and additional abstraction.
- **Per-component `getMe()` calls**: Rejected because it duplicated server state, increased API calls, and caused visible UI instability during navigation.
