---
name: 4dc-implement
description: "Use after plan.md is approved. Executes the plan in Red→Green→Refactor order, tracking progress and capturing learnings after each subtask."
---

# Implement Skill

## One Responsibility

Execute `.agent/plan.md` in controlled steps, maintain continuous progress state, and record every decision and deviation as it happens.

---

## Expected Input

- `CONSTITUTION.md`
- `.agent/increment.md` (approved)
- `.agent/plan.md` (approved)
- Current test baseline (run tests before touching anything)

---

## Concrete Output

- **`.agent/implementation.md`** — live progress log, updated after each subtask; final status either `status: complete` or `status: blocked`
- **`.agent/learnings.md`** — decisions made, deviations from plan, and lessons for `promote`

---

<HARD-GATE>
Do NOT write production code for a `[behavior]` subtask before a failing test exists for it (Red→Green→Refactor).
Do NOT mark a subtask complete without objective evidence (test output, command result, or observable behavior).
Do NOT skip or reorder subtasks without documenting the reason in learnings.md.
Do NOT proceed to the next subtask if the current one fails verification.
Do NOT mix `[tidy]` and `[behavior]` changes in the same commit.
Do NOT give the AI more context than it needs for the current subtask — constrain input to preserve output quality.
</HARD-GATE>

---

## Process

1. **Establish baseline** — run existing tests; record pass/fail state in `implementation.md`
2. **For each `[tidy]` subtask in plan.md (in order):**
   a. Make the structural change (no behavior change)
   b. Run all tests — confirm they still pass
   c. Commit separately; do not mix with behavior changes
3. **For each `[behavior]` subtask in plan.md (in order):**
   a. Constrain AI context — share only what the AI needs for this subtask
   b. Write the failing test (Red)
   c. Run it — confirm it fails for the right reason
   d. Write minimal production code (Green)
   e. Run all tests — confirm they pass
   f. Refactor structure if needed — run tests again; commit separately from behavior change
   g. Update `implementation.md`: mark subtask complete with evidence
   h. Append any decisions or surprises to `learnings.md`
4. **Final verification** — run full test suite; confirm all acceptance criteria from `increment.md` are met
5. **Mark complete** — set `status: complete` in `implementation.md`

---

## implementation.md Structure

```markdown
# Implementation: <increment goal>

status: in-progress  <!-- or: complete | blocked -->
started: <ISO date>

## Baseline
Tests before: X passing, Y failing

## Subtasks

### 1. <subtask name>
status: complete
evidence: `npm test` — 12 passing, 0 failing (added test: <name>)

### 2. <subtask name>
status: in-progress
```

---

## learnings.md Structure

```markdown
# Learnings: <increment goal>

## Decisions
- <decision>: <rationale>

## Deviations
- Subtask N: <what changed and why>

## Surprises
- <unexpected finding>

## Promote Candidates
- <ADR, doc update, or test pattern worth keeping>
```

---

## Checklist

- [ ] Baseline test run recorded
- [ ] All `[tidy]` subtasks applied and committed before any `[behavior]` subtask
- [ ] Each `[behavior]` subtask follows Red→Green→Refactor
- [ ] `[tidy]` and `[behavior]` changes in separate commits
- [ ] Each subtask has objective completion evidence
- [ ] All acceptance criteria met
- [ ] `implementation.md` status set to `complete`
- [ ] `learnings.md` has promote candidates listed

---

## Handoff

Terminal artifacts: `.agent/implementation.md` (status: complete) + `.agent/learnings.md`
Next skill: `4dc-promote` — load `skills/promote/SKILL.md`
