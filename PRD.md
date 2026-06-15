# PRD — Triangle of Love (Mobile PWA)

## Document Role

This document defines product intent, user-facing behaviour, feature scope, and success signals. It is the authoritative reference for what the product does and why.

---

## 1. Summary

Triangle of Love is a relationship health app for adults who want to do the work on their relationship — individually and together. Each person completes a short **weekly session** (under 5 minutes, fully doable on a phone during a commute) that rates how the relationship felt across three dimensions, captures two private reflections, and surfaces one actionable impulse for the week ahead.

Once enough data exists, the app suggests an optional **monthly shared conversation** — a set of trigger questions drawn from the month's patterns for partners to discuss offline.

There is no success state. Relationships are healthy, worth watching, or need attention. The app helps you see which, and gives you something concrete to do about it.

---

## 2. Goals & Non-goals

### Goals (MVP)

1. Help each person understand where their relationship stands across the three dimensions of the love triangle (Intimacy, Passion, Commitment) and how that changes over time.
2. Make reflection sustainable: one weekly session, completable in under 5 minutes on a phone, with no typing required.
3. Provide one concrete, actionable impulse per session — something specific and doable this week.
4. Surface a monthly shared conversation prompt when there is enough data, so partners can talk about what the app has noticed — on their own terms, offline.
5. Show relationship health as a continuous signal, not a score.

### Non-goals (MVP)

- No daily check-ins or daily habit tracking.
- No gamification — no achievements, streaks, badges, points, or leaderboards.
- No success/failure framing of any kind.
- Not a substitute for therapy; no diagnosis, no clinical claims.
- No in-app partner messaging or chat.
- No public social feed or community.

---

## 3. Product Principles

1. **Health, not success.** The product frame is medical, not competitive. A relationship is healthy or it needs attention. Nothing about the product should imply winning or losing.
2. **Touch first.** Every required interaction is a tap. Writing is always optional and always brief.
3. **Private by default.** Reflection text belongs to the person who wrote it and is never shared automatically.
4. **No blame UI.** Nothing in the product should make it easy to turn data into a weapon. Differences between partners are "alignment opportunities", not right/wrong.
5. **Coach voice.** Warm, direct, specific. Never clinical, never shaming, never vague.
6. **Small actions win.** Every insight must connect to a next step.
7. **Personal tier is complete.** The app works fully without a partner. Couple features are an additive unlock, not a required path.

---

## 4. Mobile Interaction Model

The primary use context is a phone in one hand — on the subway, during a coffee break, before bed. Every interaction is designed for this context.

### Binding constraints (apply to every feature)

| Rule | Detail |
|------|--------|
| Touch controls for all required input | Star tap, button press, or swipe — never a text field for required input |
| Text is always optional | Reflection notes, labels, and any free-text field can be skipped without penalty |
| Short text only | When text is offered, the field is capped at ≤140 characters |
| One screen per logical step | No pagination within a single step; scroll is preferred over page transitions |
| Session completable with no typing | A full weekly session from open to close must be achievable without opening the keyboard |
| 44 px minimum tap target | Every interactive element meets the minimum touch target height |
| Soft-keyboard safe | Layouts remain usable when the keyboard reduces the visible viewport by up to 300 px |

*The 44 px tap target, soft-keyboard safety, and 375 px baseline are set by `CONSTITUTION.md` (Target Audience section); the values here are derivative. If they diverge, CONSTITUTION is authoritative.*

### What each input type is used for

| Input | Used for |
|-------|---------|
| Star tap (1–5) | All 5 weekly questions |
| Large button / card tap | Impulse selection, monthly session opt-in, navigation |
| Optional short-text field | Reflection notes (two per session, both skippable) |
| Read-only card display | Monthly shared session trigger questions — no input required |

---

## 5. Feature Gate Model

Every feature belongs to exactly one tier. The Personal tier is a complete, self-contained product. The Couple tier adds shared features that require an active partner pairing.

### Personal tier — no partner required

| Feature | Description |
|---------|-------------|
| Weekly session | Answer 5 questions by tapping stars → triangle scores calculated → two optional reflection prompts → one impulse |
| Personal dashboard | Your own triangle trend over time, with health state label + trend line per dimension |
| Impulses | One actionable suggestion per session, drawn from your own data |
| Private reflections | Stored to your account only, never shared automatically |

### Couple tier — requires active pairing

| Feature | Description |
|---------|-------------|
| Monthly shared session | Offered when ≥2 data points exist in the calendar month; surfaces trigger questions for an offline conversation |
| Couple trend view | Aggregated triangle trend across both partners — alignment framing, no raw per-partner numbers |
| Partner invite | Send and accept a pairing invite; manage pairing state (active / dissolved) |

### Unpaired user experience

Couple features are visible to unpaired users with a single, calm prompt: *"Invite your partner to unlock this."* No empty states that look broken, no guilt, no pressure.

---

## 6. Core Loop

| Cadence | Tier | Time | What happens |
|---------|------|------|-------------|
| Weekly | Personal | ~3–5 min | 5 questions → triangle scores → reflections → impulse |
| Monthly | Couple (opt-in) | ~20 min offline | Trigger questions surfaced → partners discuss offline |

---

## 7. Weekly Session (Personal tier)

The weekly session is the only regular input. It is designed to be completed in one uninterrupted sitting of 3–5 minutes.

### Flow

1. **5 questions** — all presented on a single scrollable screen. Each answered by tapping a star (1–5). No pagination. No "next" button per question.
2. **Triangle scores calculated** — each question is designed to proxy one or more of the three dimensions (Intimacy, Passion, Commitment). The mapping is handled internally; the user sees dimension scores, not raw question averages.
3. **Two reflection prompts** — one at a time, each with an optional short-text field (≤140 chars, skippable):
   - *"What was nice about us last week?"*
   - *"What was hard or frustrating last week?"*
4. **One impulse** — a single actionable suggestion for the week ahead, displayed as a large card. The user taps to accept or skip.

### Default question set

| # | Question | Dimension(s) proxied |
|---|----------|---------------------|
| 1 | "I felt emotionally close to my partner this week." | Intimacy |
| 2 | "There was positive energy and warmth between us." | Intimacy, Passion |
| 3 | "We had fun or did something enjoyable together." | Passion |
| 4 | "I felt like we were on the same team this week." | Commitment |
| 5 | "We handled things together reliably and with care." | Commitment |

*The question set is iterable. Mapping weights are an implementation concern, not a PRD concern.*

### Impulses

Impulses are micro-interventions — specific, time-bounded, and doable this week.

- Tailored to the weakest triangle dimension or a notable trend dip
- Phrased as a direct suggestion, not a question
- Scoped to one week

Examples:
- **Intimacy:** "10-minute no-phone catch-up tonight — each person shares one high and one low from this week."
- **Passion:** "Plan something neither of you has done before. Book it before Sunday."
- **Commitment:** "15 minutes this week: calendar sync, one shared decision, one thing you're each handling."

*Full impulse library content and curation are deferred to a future increment.*

---

## 8. Monthly Shared Session (Couple tier)

### Trigger

The monthly shared session becomes available when ≥2 weekly data points exist within the current calendar month. Either partner's data counts — both completing sessions is not required.

When the trigger fires, the app surfaces a single opt-in prompt. No notification pressure. No automatic sharing of any data.

### What the session provides

The app selects 2–3 trigger questions from a curated library, matched to the month's triangle dimension patterns.

Questions are displayed as read-only cards. No input is required. The purpose is to give the couple something concrete and relevant to talk about offline — not to capture their conversation.

### Trigger question library (seed set — full library is a future increment)

Questions are mapped to triangle dimension + trend direction:

| Pattern | Example trigger question |
|---------|-------------------------|
| Intimacy trending down | "What would help you feel closer this month?" |
| Passion flat or low | "When did we last do something that genuinely surprised us both — what was it?" |
| Commitment dip | "What's one thing on our plate right now that feels like it belongs to both of us but doesn't?" |
| All dimensions healthy | "What's one thing you want to make sure we protect about us right now?" |
| Intimacy up, Passion flat | "What are we good at connecting on — and where do we feel a bit stuck?" |

*Full library size, curation process, and selection algorithm are a future increment.*

### Privacy

The trigger questions are derived from aggregated dimension trends, not from reflection text. Reflection notes are never read by the algorithm and never shown to the partner.

---

## 9. Relationship Health Model

Relationships are not scored. Each of the three triangle dimensions (Intimacy, Passion, Commitment) carries a **health state** derived from recent trend direction and magnitude.

### Health states

| State | Meaning |
|-------|---------|
| **Healthy** | Dimension is stable or improving over recent weeks |
| **Worth watching** | Dimension shows a mild or recent downward trend |
| **Needs attention** | Dimension shows a sustained or significant downward trend |

These are directional signals, not clinical thresholds. Exact transition logic is an implementation concern. No numeric thresholds are prescribed here.

### What this is not

- Not a diagnosis
- Not a relationship score
- Not a comparison between partners

---

## 10. Dashboard & Insights

### Personal dashboard (Personal tier)

Designed to be glanceable in 10 seconds on a 375 px screen.

- Triangle visualisation showing current scores for Intimacy, Passion, Commitment
- Per-dimension trend line (4–8 weeks of history)
- Per-dimension health state label alongside the trend line
- Most recent impulse, with the option to view history

### Couple trend view (Couple tier)

- Aggregated triangle trend across both partners
- Framed as an alignment view — not "you vs. them"
- Does not expose raw per-partner scores

---

## 11. Onboarding & Pairing

### Personal onboarding
1. Create account (email + password)
2. Set weekly reminder day and time window
3. Brief orientation: what the triangle means, what a session looks like

### Couple pairing (Couple tier unlock)
1. One partner generates an invite link or code
2. The other accepts
3. Pairing is active; Couple tier features unlock for both
4. Either partner can dissolve the pairing at any time — no confirmation required from the other party

---

## 12. Privacy, Safety, and Consent

1. **Reflection text is private.** It is stored only on the account that created it and is never sent to a partner endpoint.
2. **No blame UI.** No phrasing that invites comparison, e.g. "your partner rated you low on commitment."
3. **Monthly session is opt-in.** The trigger surfaces an invitation. No data is shared until the user taps to accept.
4. **Trigger questions use aggregated trends only.** Reflection text does not feed the monthly session question selection.
5. **Account controls:**
   - Dissolve pairing
   - Delete account and all associated data
6. **Disclaimer:** coaching tool, not therapy. No clinical claims.

---

## 13. Notifications & Reminders (MVP)

- Weekly reminder within the user's chosen day and time window
- Gentle nudge if the weekly session has not been completed by the end of the chosen window (once only, no guilt language)
- Opt-in notification when the monthly shared session becomes available

---

## 14. Product Metrics

These replace all gamification KPIs.

| Metric | What it signals |
|--------|----------------|
| Weekly session completion rate | Engagement and habit formation |
| Health state distribution over cohort (aggregate, anonymous) | Whether the product is reaching people who need it |
| Impulse follow-through (self-reported) | Whether impulses are useful and actionable |
| Monthly shared session opt-in rate | Whether the couple feature is valuable |
| Week 4 and Week 8 retention | Habit durability |

No streak counts, no badge completion rates, no daily active user targets.

---

## 15. MVP Scope

**Milestone 1:** Auth + onboarding + weekly session (Personal tier)
**Milestone 2:** Personal dashboard + health states + impulses
**Milestone 3:** Couple pairing + couple trend view + monthly shared session
**Milestone 4:** Reminders + PWA polish (offline resilience)

---

## 16. Open Questions / Deferred Decisions

1. Question-to-dimension weighting: which questions contribute to which dimensions, and at what weight. (Plan-phase concern for a future code increment.)
2. Health state transition thresholds: the exact trend conditions that move a dimension between states. (V2 tuning concern.)
3. Full trigger question library: content, curation process, and selection algorithm. (Separate increment.)
4. Impulse library: full content, categorisation, and personalisation logic. (Separate increment.)
