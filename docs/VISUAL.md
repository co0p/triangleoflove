# Visual Design — Triangle of Love

## Overview

The visual system is a custom CSS library built on design tokens. All component classes
reference named tokens, never hard-coded values. This single constraint is what makes the
whole UI re-themeable — swapping a theme changes every component at once, without touching
component code.

The live component gallery is the canonical reference for everything described here:
`http://localhost:5173/gallery` in development, or append `/gallery` to the deployed frontend
URL. Browse it at 375 px to see the mobile-first baseline.

---

## Design Language

### Warmth and softness

The palette comes directly from the logo mark — a sage green triangle, a dusty rose heart, and
a warm gold upward arrow. These are not loud or saturated colours. They sit in the mid-tone
range: natural, approachable, a little earthy. The overall effect should feel like something
handmade and personal, not a polished tech product.

Backgrounds and surfaces stay close to off-white. Borders and dividers are barely-there. The
colours themselves carry meaning only when they need to — a primary action in sage, an error
in muted rose, a focus ring in gold — so the screen never feels busy.

### Hierarchy through weight, not size

The type scale is conservative. Headings step up in weight and size to establish hierarchy, but
gaps between levels are small enough that the page reads as calm, not dramatic. Body copy sits
at a comfortable line-height. Helper text and captions pull back in colour rather than
shrinking aggressively.

The font is the system stack — whatever the device considers its default sans-serif. This keeps
the app feeling native on mobile without introducing a web font load.

### Rounded and generous

Corners are generously rounded, especially on buttons (pill-shaped by default) and cards. This
softness echoes the organic quality of the logo. Sharp corners or tight radii would fight the
brand.

Spacing follows a consistent scale. Padding inside components is never cramped.

### Mobile is the primary context

The primary user is on a phone. Every layout and interaction decision starts from that
assumption.

- **Tap targets are never smaller than 44 px.** This is the minimum recommended by Apple and
  Google for comfortable, error-free tapping. Buttons carry a `min-height` to enforce it.
- **The soft keyboard is an adversary.** When an input is focused on mobile, the soft keyboard
  can reduce the visible viewport by 300 px or more. Vertical layouts must not rely on
  `align-items: center` without also accounting for this collapse — use top-anchored
  alignment with generous top padding instead.
- **Text truncates, never wraps, in fixed-height bars.** Navigation bars and headers have
  limited height. Any text that could vary in length (names, titles) must truncate with
  an ellipsis rather than overflowing or reflowing the surrounding layout.
- **375 px is the baseline.** Design and review at this width first. Wider behaviour is
  additive.

### The focus ring as a brand moment

The keyboard focus ring is warm gold — the same colour as the arrow in the logo. It should
never be removed or suppressed. It is the brand's way of saying "you are here" for keyboard
and switch-access users. Every interactive element shows it.

### Themes without component changes

A theme is a named override of the token layer, not a separate stylesheet. Applying a theme
class to any container element re-colours everything inside it. Components never know which
theme is active — they only see tokens.

---

## Component Principles

1. **Token-only references.** Component classes must use `var(--token)` for every value.
   Raw hex, raw sizes, and raw radii belong only in the token definitions in `library.css`.

2. **States before variants.** Before adding a new component variant, ensure the existing
   states (default, hover, focus, disabled, error) are fully specified. Incomplete states
   leave gaps that become bugs at runtime.

3. **Mobile-first by default.** The container is narrow and centred. Every component must
   hold at 375 px first. Wider layouts are additive, never assumed.

4. **No external dependencies.** The library is one `.css` file. Adding a component means
   adding class definitions to that file. No build plugins, no utility generators, no
   framework defaults to override.

---

## Gallery

The gallery renders every component class in every documented state with a live theme switcher.
If a class is not in the gallery, it is effectively undiscoverable by future contributors.

Source: `services/frontend/src/views/GalleryView.vue`
