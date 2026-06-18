# Triangle of Love — Implemented User Flow (Current)

This document reflects the **currently implemented** navigation and interaction flows in the codebase.

---

## 1) Activity Diagram (Current)

```mermaid
flowchart TD
  A([Open App]) --> B{Authenticated?}

  B -- No --> C[Login]
  C --> D{Valid credentials?}
  D -- No --> C1[Show auth error]
  C1 --> C
  D -- Yes --> E[Dashboard]

  B -- Yes --> E[Dashboard]

  %% Primary authenticated routes
  E --> F[Check-in page]
  E --> G[Insights weekly matrix]
  E --> H[Pairing page]
  E --> I[Profile page]

  %% Check-in flow
  F --> F1[Load today's check-in]
  F1 --> F2[Edit 7 ratings + optional note]
  F2 --> F3[Save today check-in]
  F3 --> F4[Show save confirmation state]
  F4 --> E

  %% Insights flow
  G --> G1[Load 7-day insights window]
  G1 --> G2[Select day]
  G2 --> G3[Load daily insight for date]
  G3 --> G

  %% Pairing flow
  H --> H1[Load invite code and couple status]
  H1 --> H2[Regenerate invite code]
  H1 --> H3[Connect with partner code]
  H1 --> H4[Unpair current couple]
  H2 --> H
  H3 --> H
  H4 --> H

  %% Profile flow
  I --> I1[Load current user profile]
  I1 --> I2[Change password]
  I2 --> I3{Current password valid?}
  I3 -- No --> I4[Show validation/auth error]
  I4 --> I
  I3 -- Yes --> I5[Password changed confirmation]
  I5 --> I

  %% Admin
  E --> J{Role is admin?}
  J -- Yes --> K[Admin users page]
  K --> K1[List users]
  K1 --> K2[Activate/deactivate user]
  K2 --> K

  %% Exit
  I --> L[Logout]
  L --> C
```

---

## 2) Screen Map (Current Routes)

### Public
- **P1** Login (`/login`)
- **P2** Register (`/register`)

### Authenticated
- **A1** Dashboard (`/dashboard`)
- **A2** Session (`/session`)
- **A3** Insights weekly (`/insights`)
- **A4** Pairing (`/pairing`)
- **A5** Profile (`/profile`)

### Admin
- **ADM1** Admin users (`/admin/users`)

---

## 3) Implemented UX Rules

- Sessions are **daily** and saved at `GET/PUT /api/v1/sessions/today`.
- Check-in inputs are seven ratings (six relationship metrics + mood) plus an optional private note.
- Insights include:
  - weekly matrix view from `GET /api/v1/insights`
  - per-date insight view from `GET /api/v1/insights/{date}`
- Pairing supports invite code retrieval/regeneration, connect, and unpair.
- Profile supports viewing own account and changing password.
- Admin users can list accounts and toggle activation state.

---

## 4) Not Implemented Yet

The following concepts appear in product planning docs but are not currently present in the app flow:
- Weekly session flow
- Impulse generation and selection
- Monthly shared session trigger questions
- Achievements flow
