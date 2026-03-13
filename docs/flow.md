# Triangle of Love Coach — Wireframe-ready Flow (MVP)

This document is a **wireframe-ready interaction flow** for generating low-fidelity screens and thinking through navigation + states.

---

## 1) Activity Diagram (MVP)

```mermaid
flowchart TD
  %% Triangle of Love Coach — Activity Diagram (MVP)

  A([Start / Open App]) --> B{Authenticated?}

  %% Auth + Onboarding
  B -- No --> C[Sign up / Log in]
  C --> D[Onboarding: Rules of engagement\n(no blame, private-by-default, not therapy)]
  D --> E[Pairing: Invite link/code\n(share code OR enter code)]
  E --> F{Partner paired?}
  F -- No --> E
  F -- Yes --> G[Set cadence\n- daily reminder window\n- weekly review day/time]
  G --> H[Home / Dashboard]

  %% Main hub
  B -- Yes --> H[Home / Dashboard]

  %% Home entry points
  H --> I{What is due?}
  I -->|Daily check-in due| J[Daily Check-in\n(<=5 star questions)]
  I -->|Weekly review due| K[Weekly Review]
  I -->|Nothing due| L[View Insights / Trends]
  I -->|Manage reminders/settings| M[Settings / Account Controls]

  %% Daily loop
  J --> J1[Answer 5 star questions\n(required)]
  J1 --> J2[Optional: one-line private note]
  J2 --> J3[Submit daily check-in]
  J3 --> J4[Completion confirmation\n+ optional micro-tip]
  J4 --> H

  %% Weekly loop
  K --> K1[Rate Triangle\nIntimacy / Passion / Commitment (1–5)]
  K1 --> K2[Private reflections (required)\n- What was nice?\n- What was annoying/hard?]
  K2 --> K3{Choose focus area?\n(optional)}
  K3 -->|Yes| K4[Pick focus\n(Intimacy/Passion/Commitment/Communication)]
  K3 -->|No| K5[System infers focus\n(weakest side/trend dip)]
  K4 --> K6[Generate Impulses\n1 primary + 1 bonus\n(1-week scope)]
  K5 --> K6
  K6 --> K7[User selects impulse focus\n(primary or bonus)]
  K7 --> K8{Schedule mid-week reminder?\n(optional)}
  K8 -->|Yes| K9[Set reminder for impulse]
  K8 -->|No| K10[Skip reminder]
  K9 --> K11[Weekly complete confirmation]
  K10 --> K11
  K11 --> H

  %% Insights & dashboard behavior
  L --> L1[Show personal trends\n(4–8 weeks)]
  L1 --> L2[Show couple aggregate trend\n(no partner-vs-you exact numbers)]
  L2 --> L3[Alignment summary (gentle)\nRecommend actions\n(no blame UI)]
  L3 --> H

  %% Achievements (can be accessed from home/insights)
  H --> N[Achievements]
  N --> N1[View achievements\n(secret/private by default)]
  N1 --> N2[Self-claim achievement\n(default measurement)]
  N2 --> N3{Request partner confirmation?\n(optional, per achievement)}
  N3 -->|No| N6[Mark as secret/private progress]
  N3 -->|Yes| N4[Send confirmation request\n(user-triggered only)]
  N4 --> N5[Partner responds\n(confirm/ignore/decline)]
  N5 --> N7[Update achievement status\n(confirmed or not)]
  N6 --> H
  N7 --> H

  %% Notifications
  subgraph R[Notifications & Reminders]
    R1[Daily reminder\n(in chosen time window)]
    R2[Weekly reminder\n(on chosen day/time)]
    R3[Gentle nudge if weekly pending\n(no guilt language)]
    R4[Optional impulse mid-week reminder]
  end

  %% Settings / safety / controls
  M --> M1[Edit cadence/reminders]
  M --> M2[Privacy & disclaimers\n(coaching tool, not therapy)]
  M --> M3[Unpair / leave couple]
  M --> M4[Delete account & data]
  M1 --> H
  M2 --> H
  M3 --> H
  M4 --> Z([End])
```

---

## 2) Screen Map (use these IDs as Figma frame names)

### Onboarding / Account
- **A0** Splash / Loading
- **A1** Sign up / Log in
- **A2** Rules of Engagement
- **A3** Pairing method (Invite / Enter code)
- **A4** Share invite link/code
- **A5** Enter invite code
- **A6** Pairing pending
- **A7** Set cadence (daily window, weekly day/time)
- **A8** Onboarding complete → Home

### Home + Check-ins
- **H1** Home / Dashboard (today status + CTAs)
- **D1** Daily check-in (5 star questions + optional note)
- **D2** Daily complete (confirmation + micro-tip)

### Weekly review + impulses
- **W1** Weekly review intro (privacy reminder)
- **W2** Triangle ratings (Intimacy/Passion/Commitment)
- **W3** Private reflections (Nice / Hard)
- **W4** Choose focus (optional)
- **W5** Impulse recommendations (Primary + Bonus)
- **W6** Impulse detail
- **W7** Weekly plan confirm + reminder toggle
- **W8** Weekly complete (summary)

### Insights
- **I1** Insights overview (personal trends + couple aggregate)
- **I2** Trend detail (dimension drilldown)
- **I3** Alignment summary (gentle recommendations)

### Achievements
- **G1** Achievements home
- **G2** Achievement detail
- **G3** Claim achievement (self-claim)
- **G4** Partner confirmation opt-in
- **G5** Confirmation sent / pending
- **G6** Partner confirmation inbox
- **G7** Confirmation result

### Settings / Safety
- **S1** Settings home
- **S2** Reminders & cadence
- **S3** Privacy & disclaimers
- **S4** Unpair / leave couple
- **S5** Delete account & data

---

## 3) Wireframe Requirements / UX Rules (from PRD)

- **Reflections are private-by-default** and **not shared in-app**.
- Avoid “blame UI” (no “your partner rated you low” language).
- Insights should show **personal trends** + **couple aggregate** without exact partner-vs-you numbers.
- Achievements default to **secret/private**; partner-confirmation is **opt-in** and **user-triggered** only.
- Reminders: daily + weekly + gentle weekly pending nudge + optional mid-week impulse reminder.

---