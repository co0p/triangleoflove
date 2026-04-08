# CSS Library-First Pattern

**Date**: 2026-04-08
**Status**: Accepted

## Context

CSS is split between `library.css` and per-component `<style scoped>` blocks without a governing rule. Raw hex and rgba values have drifted into component style blocks, undermining the token system. The design-system visual reference (`GalleryView.vue`) lives inside Vue routing, coupling it to app auth and runtime.

## Decision

1. **Single source of truth**: All shared CSS classes live in `services/frontend/src/assets/library.css`. No class that appears in two or more components may remain in a component `<style>` block.

2. **Token compliance**: All color and spacing values in style rules must reference `var(--token)`. Raw hex (`#...`) or rgba values are permitted only inside the `:root` block of `library.css`.

3. **Extraction threshold**: A class is extracted to `library.css` when it is used in two or more components. Single-use layout rules remain in the component's `<style scoped>` block.

4. **Design-system reference**: A standalone `design-system.html` lives in `docs/design-system.html`. It references `library.css` via a relative path (`../services/frontend/src/assets/library.css`) and showcases every shared class without requiring the app to run. It is updated in the same commit as any library change.

## Alternatives Considered

- **Extracting by semantic generality** (e.g. a class that "looks reusable"): Rejected — introduces subjective judgment with no enforceable signal. The duplication threshold is objective.
- **Co-locating `design-system.html` in `docs/`**: Accepted — visible alongside other design documentation; relative path to `library.css` is stable.
- **Hot-reload support for the standalone HTML**: Deferred — out of scope; convention over tooling is sufficient at current team size.

## Consequences

- Component `<style>` blocks are audited at PR review for raw values and duplicate classes.
- `design-system.html` must be updated whenever `library.css` changes; enforced by code review.
- `GalleryView.vue` and the `/gallery` route are removed; `docs/design-system.html` is the new visual reference.
