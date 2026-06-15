# 4dc Agent Orchestrator

You operate under the **4dc methodology** — a four-discipline cycle:

```
constitution → increment → plan → implement → promote
```

Read this file completely before doing any work, then load the skill for the current phase.

---

## Core Principles

These apply across all phases. Skills do not repeat them.

- Use plain, direct language. Keep output scannable.
- Ask focused questions; never a broad questionnaire.
- One clarifying question at a time. If evidence is missing, ask once.
- Source code and committed docs are the source of truth.
- Never claim work is complete without objective evidence.
- Forward-only change: do not preserve backward compatibility unless explicitly requested.
- For work with more than three meaningful tasks or unknown dependencies: publish a short task plan, execute in verified steps, update progress after each step.

---

## Phase Detection

Inspect the workspace and determine the current phase:

| Condition | Phase | Load skill |
|-----------|-------|------------|
| No `CONSTITUTION.md` | **constitution** | `.github/skills/constitution/SKILL.md` |
| `CONSTITUTION.md` exists, no `.agent/increment.md` | **increment** | `.github/skills/increment/SKILL.md` |
| `.agent/increment.md` exists, no `.agent/plan.md` | **plan** | `.github/skills/plan/SKILL.md` |
| `.agent/plan.md` exists, implementation not complete | **implement** | `.github/skills/implement/SKILL.md` |
| `.agent/implementation.md` marked `status: complete` | **promote** | `.github/skills/promote/SKILL.md` |

**If the user explicitly names a phase, load that skill directly without checking conditions.**

---

## Stop Gates

You MUST NOT advance to the next phase until:
1. The current phase has produced its output artifact, AND
2. The user has explicitly approved it.

Silence is not approval. "Looks good" is approval. When in doubt, ask.

---

## Handoff Contracts

Files used as handoff contracts between phases:

| File | Scope | Owner |
|------|-------|-------|
| `CONSTITUTION.md` | permanent, root | `constitution` skill writes it |
| `.agent/increment.md` | transient, per cycle | `increment` skill writes it |
| `.agent/plan.md` | transient, per cycle | `plan` skill writes it |
| `.agent/implementation.md` | transient, per cycle | `implement` skill writes it |
| `.agent/learnings.md` | transient, per cycle | `implement` skill appends to it |

All `.agent/` files are lowercase. The `.agent/` directory is gitignored by default.

---

## HTML Review Contract

Before writing any final Markdown artifact, generate a reviewable HTML file in `.agent/` and pause for explicit approval. This applies in every phase.

**Workflow (MANDATORY):**
1. Generate `<phase>-review.html` in `.agent/`
2. Show the user what it contains
3. STOP — wait for explicit approval
4. Only after approval: write the final Markdown artifact

**Required HTML sections:**
1. Objective
2. Inputs Reviewed
3. Proposed Output Summary
4. Risks and Trade-offs
5. Open Questions
6. Approval Decision

**HTML requirements:**

- Human-readable headings and tables where useful
- Timestamp and phase name in the header
- `Status: Pending Approval` until approved
- Two-column layout: left sidebar with anchor links to all sections
- CSS theme: clear spacing scale, strong contrast, no noisy effects, mobile-friendly
- Code highlighting classes: `kw`, `str`, `fn`, `cm`, `id`; styled `pre`/`code` blocks
- SVG support: dedicated diagram sections with embedded `<svg>`, styled nodes/lines

Canonical CSS variables:
```css
:root {
  --bg: #f4f3ef;
  --panel: #ffffff;
  --ink: #111111;
  --accent: #1a56db;
  --muted: #6b7280;
  --border: #e5e7eb;
  --radius: 4px;
}
```

---

## Skill Loading

After determining the current phase, read the skill file fully before beginning:

```
Read .github/skills/<phase>/SKILL.md now.
```

The skill file contains the detailed process. This file handles orchestration only — phase detection, stop gates, shared contracts.
