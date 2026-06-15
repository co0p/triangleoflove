---
name: 4dc-plan
description: "Use after increment.md is approved. Converts increment intent into an ordered, verifiable technical execution plan with concrete subtasks."
---

# Plan Skill

## One Responsibility

Define HOW to deliver `.agent/increment.md` — an ordered sequence of actionable subtasks with explicit verification points.

---

## Expected Input

- `CONSTITUTION.md`
- `.agent/increment.md` (must be approved)
- Current codebase structure (read relevant files)

---

## Concrete Output

`.agent/plan.md` containing:
- **Goal**: copied from `increment.md` — one sentence
- **Tidy First?**: yes/no — whether existing code should be structurally tidied before implementing; if yes, list the specific tidyings
- **Approach**: 2–3 sentences on strategy; no code yet
- **Subtasks**: ordered list, each tagged with a type:
  - `[tidy]` — structural change only, no behavior change; must leave all existing tests passing
  - `[behavior]` — new or changed observable outcome; requires a failing test first
  - Each subtask also has: description (what, not how), verification step, and dependencies on prior subtasks
- **Risks**: known unknowns that could block execution

> Tidy subtasks always come before behavior subtasks. They exist solely to make the behavior change easier or safer to implement.

---

<HARD-GATE>
Do NOT start implementation during this phase — no code, no file edits.
Do NOT write plan.md until the HTML review is approved.
Do NOT list subtasks without a type tag and a verification step.
Every `[behavior]` subtask must map to at least one acceptance criterion from increment.md.
Do NOT add a `[tidy]` subtask that does not directly serve the behavior change that follows it.
</HARD-GATE>

---

## Process

1. **Read inputs** — `CONSTITUTION.md`, `.agent/increment.md`, relevant source files
2. **Tidy First?** — assess whether structural tidyings (guard clauses, dead code removal, cohesion reordering, etc.) would reduce the cost or risk of the behavior subtasks. If yes, list them as `[tidy]` subtasks at the top of the plan. A tidy subtask must have a clear answer to: *"does this make the behavior change easier?"* — if not, skip it.
3. **Identify risks** — unknown dependencies, test gaps, ambiguous requirements
4. **Draft subtasks** — ordered, each tagged `[tidy]` or `[behavior]`, independently verifiable, sized for one focused work session
5. **Generate `.agent/plan-review.html`** — present the plan with full traceability to acceptance criteria
6. **STOP** — wait for explicit approval or revision
7. **On approval** — write `.agent/plan.md`

---

## Checklist

- [ ] `increment.md` acceptance criteria read
- [ ] Tidy First? assessed and documented
- [ ] All `[tidy]` subtasks justified by the behavior change that follows
- [ ] Relevant source files scanned
- [ ] Every subtask has a type tag and a verification step
- [ ] Every `[behavior]` subtask maps to an acceptance criterion
- [ ] Risks documented
- [ ] HTML review generated and shown
- [ ] User approval received
- [ ] `.agent/plan.md` written

---

## Handoff

Terminal artifact: `.agent/plan.md`
Next skill: `4dc-implement` — load `skills/implement/SKILL.md`
