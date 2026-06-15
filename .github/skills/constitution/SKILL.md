---
name: 4dc-constitution
description: "Use when CONSTITUTION.md is missing or needs updating. Reads project context, asks focused questions, and produces project guardrails and SDLC standards."
---

# Constitution Skill

## One Responsibility

Create or update `CONSTITUTION.md` — the project's durable engineering guardrails.

---

## Expected Input

- Existing `CONSTITUTION.md` (if present)
- `README.md`
- Current project structure and any existing docs

---

## Concrete Output

`CONSTITUTION.md` containing:
1. Engineering principles grounded in XP, Tidy First, and lean software development
2. Architectural boundaries and dependency direction
3. Testing strategy and quality gates: Red→Green→Refactor with structural and behavioral changes separated
4. Tidy First policy: when to tidy before changing behavior, what constitutes a valid tidying
5. Documentation rules and ADR policy
6. SDLC artifact expectations: the `.agent/` contract (which files, what lifecycle)
7. AI collaboration rules: constrain context per task, preserve optionality, maintain human judgment over architectural decisions

---

<HARD-GATE>
Do NOT write CONSTITUTION.md until the HTML review in `.agent/constitution-review.html` has been explicitly approved.
Do NOT ask more than 5 questions per round.
Do NOT include implementation details — CONSTITUTION.md contains guardrails, not recipes.
Do NOT copy generic principles from the internet. Every rule must be justified by this project's specific context.
</HARD-GATE>

---

## Process

1. **Read project context** — scan `README.md`, existing `CONSTITUTION.md`, directory structure, any ADRs or docs
2. **Ask 3–5 focused questions** — surface constraints, pain points, and non-negotiables one round at a time
3. **Generate `.agent/constitution-review.html`** — present proposed guardrails in review format
4. **STOP** — wait for explicit approval or revision requests
5. **On approval** — write `CONSTITUTION.md`

---

## Checklist

- [ ] Existing docs read
- [ ] 3–5 questions asked and answered
- [ ] HTML review generated and shown
- [ ] User approval received
- [ ] `CONSTITUTION.md` written

---

## Handoff

Terminal artifact: `CONSTITUTION.md`
Next skill: `4dc-increment` — load `skills/increment/SKILL.md`
