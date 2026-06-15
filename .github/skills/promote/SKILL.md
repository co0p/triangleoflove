---
name: 4dc-promote
description: "Use after implementation.md is marked complete. Reviews all .agent/ artifacts, proposes promotions to permanent docs, and closes the cycle."
---

# Promote Skill

## One Responsibility

Merge durable outcomes from the `.agent/` working set into permanent project artifacts before the branch merges.

---

## Expected Input

- `CONSTITUTION.md`
- `.agent/increment.md`
- `.agent/plan.md`
- `.agent/implementation.md` (status: complete)
- `.agent/learnings.md`
- Existing `docs/`, ADR log, and any project documentation

---

## Concrete Output

One or more of the following, per approval:
- Updated `CONSTITUTION.md` (if guardrails need revision)
- New ADR in `docs/adr/` (for significant architectural decisions)
- Updated `README.md` or other docs (for changed behavior or usage)
- Deleted or archived `.agent/` files after promotion (keeping `.agent/` clean for next cycle)

---

<HARD-GATE>
Do NOT write permanent docs until each promotion candidate has been individually approved.
Do NOT promote guesses or plans — only promote what was actually built and verified.
Do NOT delete .agent/ files until all promotions are written and confirmed.
Present each candidate separately with destination path and rationale.
</HARD-GATE>

---

## Process

1. **Read all `.agent/` artifacts** — full review of increment, plan, implementation, and learnings
2. **Identify promotion candidates** from `learnings.md`'s "Promote Candidates" section plus your own review
3. **For each candidate:**
   - State: what it is, where it goes, why it's durable
   - Generate `.agent/promotion-review.html` covering all candidates
4. **STOP** — present HTML review, await per-candidate approval
5. **On approval** — write each approved permanent artifact
6. **Clean up** — archive or delete `.agent/` files for this cycle

---

## Promotion Categories

| Type | Trigger | Destination |
|------|---------|-------------|
| Architecture decision | Non-obvious choice with lasting impact | `docs/adr/ADR-<date>-<slug>.md` |
| Design optionality | A decision that bought or sold future options — worth naming for the next cycle | `docs/adr/ADR-<date>-<slug>.md` |
| Guardrail update | Constitution rule violated, needs clarification | `CONSTITUTION.md` |
| Behavior change | Public API, CLI, or user-facing behavior changed | `README.md` |
| Test pattern | New testing approach worth standardizing | `CONSTITUTION.md` testing section |
| Known issue | Found but not fixed this cycle | `docs/known-issues.md` |

---

## Checklist

- [ ] All `.agent/` artifacts read
- [ ] Promotion candidates identified and categorized
- [ ] HTML review generated covering all candidates
- [ ] User approval received per candidate
- [ ] Each approved artifact written to permanent location
- [ ] `.agent/` files cleaned up

---

## Handoff

Terminal artifacts: permanent docs updated, `.agent/` clean
Cycle complete. Next action: `4dc-increment` for the next cycle — load `skills/increment/SKILL.md`
