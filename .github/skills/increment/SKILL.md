---
name: 4dc-increment
argument-hint: "short increment intent, e.g. 'feature: export weekly summary'"
description: "Use after CONSTITUTION.md exists and before any implementation. Defines the next narrow, testable behavior change — WHAT and WHY only, no technical detail."
---

# Increment Skill

## One Responsibility

Define one small, outcome-focused increment with measurable acceptance criteria and explicit out-of-scope boundaries.

---

## Expected Input

- `CONSTITUTION.md`
- User intent (one sentence or short phrase describing the desired outcome)

---

## Concrete Output

`.agent/increment.md` containing:
- **Goal**: one sentence — the user-observable outcome
- **Acceptance criteria**: 2–5 binary, verifiable conditions
- **Out of scope**: explicit exclusions that prevent scope creep
- **Constitution constraints**: which guardrails apply to this increment

---

<HARD-GATE>
Do NOT include technical design, file names, implementation approaches, or coding detail in increment.md.
Do NOT start a plan or any implementation work during this phase.
Do NOT approve an increment with vague acceptance criteria ("works correctly", "feels right").
One increment per cycle — if scope expands, split into separate increments.
</HARD-GATE>

---

## Process

1. **Read context** — `CONSTITUTION.md`, any prior `.agent/` files from the last cycle
2. **Clarify intent** — ask 1–3 focused questions to sharpen WHAT and WHY
3. **Generate `.agent/increment-review.html`** — show proposed increment definition
4. **STOP** — wait for explicit approval or revision requests
5. **On approval** — write `.agent/increment.md`

---

## Checklist

- [ ] `CONSTITUTION.md` read
- [ ] User intent clarified with ≤3 questions
- [ ] Acceptance criteria are binary and verifiable
- [ ] Out-of-scope list is non-empty
- [ ] HTML review generated and shown
- [ ] User approval received
- [ ] `.agent/increment.md` written

---

## Handoff

Terminal artifact: `.agent/increment.md`
Next skill: `4dc-plan` — load `skills/plan/SKILL.md`
