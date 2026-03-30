# ADR: Custom CSS Only — No CSS Framework

## Context

The frontend needed a visual system for the first styled increment. The team evaluated
whether to adopt an external CSS framework (Tailwind CSS was the original recommendation
in DESIGN.md) or build a lightweight custom library from scratch.

The brand palette is small and specific — three named colours derived directly from the
logo mark, with a supporting neutral scale. The component vocabulary for the initial
increment is narrow: buttons, inputs, cards, a navbar, and an error alert.

## Decision

We use a single custom CSS file (`services/frontend/src/assets/library.css`) built on
CSS custom properties. Design tokens are declared on `:root`. Component classes reference
only token variables. Themes are token overrides on a container class. No external CSS
framework or utility generator is introduced.

## Consequences

- **Benefits:** Zero build-step overhead; a single replaceable file; token-only component
  classes make the entire UI re-themeable by overriding a small set of variables; naming
  conventions are fully under team control.
- **Drawbacks:** No utility class system; layout patterns must be written by hand; no
  community component library to draw from.
- **Trade-offs:** Acceptable at this scale. The component surface is small enough that
  hand-written CSS is faster than learning and working around a framework's constraints.

## Alternatives Considered

- **Tailwind CSS**: Rejected — requires a build plugin, constrains class naming, and
  imposes a design system that conflicts with the brand-specific aesthetic.
- **Bootstrap**: Rejected — imposes a visual language that would require significant
  overrides to match the brand.
