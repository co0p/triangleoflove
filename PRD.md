# PRD — Triangle of Love Coach (Mobile PWA)

## Document Role

This document defines product intent, user-facing behavior, scope, and success criteria.

---

## 1. Summary

Triangle of Love Coach is a relationship improvement and tracking app for couples. Each partner privately logs short **daily check-ins** (max 5 star questions) and a **weekly review** (triangle ratings + private reflections). The app visualizes trends and gently identifies areas for alignment, then provides **coach-style impulses** (micro-actions) and **gamified achievements** that help couples build mature relationship habits over time.

Reflections are **private by default** and are not shared in-app; partners can choose to discuss them verbally.

---

## 2. Goals & Non-goals

### Goals (MVP)
1. Help couples understand where they are in the **love triangle** and how it changes over time.
2. Encourage consistent engagement with minimal effort:
   - daily: 30–60 seconds
   - weekly: ~5 minutes
3. Provide actionable, personalized **impulses** that feel like a wingman:
   - specific, time-bounded, doable
4. Reinforce growth behaviors via **achievements** (habit gamification), including:
   - secret/private achievements (default)
   - optional partner-confirmation achievements
5. Improve relationship maturity skills:
   - appreciation
   - listening and non-interruption
   - conflict hygiene and repair
   - planning and reliability

### Non-goals (MVP)
- Not a substitute for therapy; no diagnosis, no mental health claims.
- No public social feed or community.
- No in-app partner messaging replacement.
- No “leaderboards” or competitive scoring between partners.

---

## 3. Target Users & Personas

### Primary
- **Growth-minded couples** who want structure and accountability.
- **Busy couples** who want quick, focused prompts and a “what should we do this week?” plan.

### Secondary (handled carefully)
- Couples with mild friction who want a gentle, non-blaming way to improve communication.

---

## 4. Product Principles

1. **Private-by-default reflections:** text prompts are personal; sharing happens offline by choice.
2. **No weaponization:** avoid UIs that encourage blame or “scorekeeping.”
3. **Coach voice:** warm, direct, playful where appropriate; never clinical or shaming.
4. **Small actions win:** always turn insight into a next step (impulse).
5. **Two perspectives are data:** differences indicate “alignment opportunity,” not right/wrong.
6. **Low friction:** every loop must feel lightweight and sustainable.

---

## 5. Core Loop

### Daily (30–60 seconds)
1. User completes daily check-in (max 5 star questions).
2. Optional one-line note (private).
3. App confirms completion and optionally shows a micro-tip.

### Weekly (~5 minutes)
1. User rates Intimacy/Passion/Commitment (1–5 each).
2. User writes two private reflections:
   - “What was nice last week?”
   - “What was annoying or hard last week?”
3. App suggests:
   - 1 primary impulse for the week
   - 1 optional bonus impulse
4. User selects a focus and optionally schedules a reminder.

---

## 6. Core Features (MVP)

### 6.1 Onboarding & Couple Pairing
- Create account
- Pair with partner via invite link/code
- Set cadence:
  - daily reminder time window
  - weekly review day/time
- Short “rules of engagement” screen (growth mindset, no blame).

### 6.2 Daily Check-in (Max 5 Questions)
Format: **1–5 stars** each (required), optional note.

Default question set (can be iterated):
1. “I felt close to my partner today.” *(Intimacy proxy)*
2. “We had positive energy / fun today.” *(Passion proxy)*
3. “I felt supported / we were a team today.” *(Commitment proxy)*
4. “Communication felt healthy today.” *(Skill proxy)*
5. “My stress level today.” *(Context; personal)*

### 6.3 Weekly Review
Required:
- Triangle ratings (1–5): **Intimacy**, **Passion**, **Commitment**
- Private prompts:
  - “What was nice last week?”
  - “What was annoying or hard last week?”

Optional:
- Choose one focus area for next week: Intimacy / Passion / Commitment / Communication.

### 6.4 Dashboard & Insights
- **Triangle trend view** over time (4–8 weeks)
- **Health metrics trend** (daily roll-ups, weekly averages)
- **Alignment summary** (gentle):
  - highlight dimension(s) trending lower and recommend actions
  - avoid presenting partner vs you as a “competition”

**Recommended MVP default:** show *personal trends* and a *couple aggregate* without revealing exact partner-vs-you numbers.

### 6.5 Impulses (“Wingman moves”)
Impulses are micro-interventions:
- tailored to the weakest triangle side, recent trend dips, or chosen focus
- phrased as an actionable suggestion (“Here’s your move”)
- scoped to *one week*

Examples:
- **Intimacy:** “10-minute no-phone ‘High/Low’ talk + one appreciation.”
- **Passion:** “Plan one novelty date (new place or new activity).”
- **Commitment:** “15-minute team planning: calendar + logistics + one shared goal.”
- **Communication:** “Try the ‘reflect back’ rule: repeat their point before responding.”

### 6.6 Achievements (Gamification)
Achievements reinforce habits without creating competition.

Properties:
- Category: Romance / Excitement / Intimacy / Commitment / Communication
- Visibility:
  - **Secret/private** (default)
  - **Partner-confirmed** (optional per achievement)
- Measurement:
  - self-claim (default)
  - optional confirmation request sent to partner (user-triggered only)

Seed examples:
- **Romantic — “Flower Man”**: weekly fresh flowers, 4-week streak.
- **Excitement — “Date Night Architect”**: plan 3 different date nights in a month.
- **Communication — “The Listener”**: 4 serious talks in a month without interrupting (best as partner-confirmed, but can remain secret).

---

## 7. Notifications & Reminders (MVP)
- Daily reminder within chosen time window.
- Weekly reminder on chosen day/time.
- Gentle nudge when a weekly review is pending (no guilt language).
- Optional mid-week reminder for the selected impulse.

---

## 8. Privacy, Safety, and Consent Requirements (Critical)
1. Reflection text is **private** and never shared to partner endpoints.
2. Avoid “blame UI”:
   - no phrasing like “your partner rated you low”
3. Partner-confirmation achievements:
   - must be **opt-in** per achievement
   - confirmation request is user-triggered
   - never expose secret achievement intent unless user chooses
4. Account controls:
   - unpair/leave couple
   - delete account and data
5. Clear disclaimers:
   - coaching tool, not therapy

---

## 9. Success Metrics (KPIs)

### Activation
- Invite acceptance rate within 48 hours
- % of couples completing first weekly review within 7 days

### Engagement
- Daily check-in completion rate per user
- Weekly review completion rate per couple
- Impulse selection rate after weekly review

### Retention
- Week 4 and Week 8 couple retention

### Outcome proxies (non-medical)
- Improvements in triangle ratings over 8 weeks
- Reduced “alignment variance” (measured carefully; aggregate-first)
- Achievement completion rate

---

## 10. MVP Scope & Milestones (Suggested)
**Milestone 1:** Auth + pairing + daily check-ins  
**Milestone 2:** Weekly review + dashboard trends  
**Milestone 3:** Impulses + achievements (private)  
**Milestone 4:** Optional partner-confirmation flow + PWA polish (offline + reminders)

---

## 11. Open Questions / Decisions
1. Exact “alignment” presentation:
   - aggregate-only vs showing partner-vs-you differences
2. Achievement confirmation UX:
   - how to request confirmation without creating pressure
3. Content strategy:
   - how many impulses ship in v1 (recommend 30–60 as a starter pack)
