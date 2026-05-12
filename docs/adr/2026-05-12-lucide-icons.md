# ADR: Lucide Icons as the Icon Library

## Context

The CheckinMatrix component needed small, recognisable icons to label the
three relationship dimensions (Intimacy, Commitment, Passion) in a compact
row header. An icon library was needed rather than hand-rolled SVGs to keep
assets maintainable and consistent with future icon use.

## Decision

Adopt `lucide-vue-next` as the project's icon library. The three dimension
icons are:

- **Heart** → Intimacy
- **Anchor** → Commitment
- **Flame** → Passion

Icons are imported individually from `lucide-vue-next` and rendered as Vue
components with explicit `size` and `stroke-width` props.

## Consequences

- **Benefits:** MIT licensed; SVG stroke-based (scales cleanly at any size);
  tree-shakeable (only imported icons reach the bundle); large icon set
  available for future use; first-class Vue 3 component API.
- **Drawbacks:** Adds a production dependency; icon style (stroke) must remain
  consistent if other icon sources are ever mixed in.
- **Trade-offs:** Stroke icons suit the app's earthy, minimal aesthetic; if a
  filled or duotone style were needed in future, a library swap or override
  would be required.

## Alternatives Considered

- **Heroicons**: Also MIT and stroke-based, but the Vue package is less
  actively maintained for Vue 3.
- **Hand-rolled SVGs**: Zero dependency cost, but duplicates work and
  introduces inconsistency as the icon set grows.
- **Material Icons**: Large and well-supported, but the visual style (filled,
  rounded) does not match the app's aesthetic.
