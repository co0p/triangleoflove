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
- Updated `docs/roadmap.md` — feature moved from In Progress to Done, acceptance test link added
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
2. **Identify promotion candidates** from `learnings.md`’s "Promote Candidates" section plus your own review
3. **For each candidate:**
   - State: what it is, where it goes, why it’s durable
   - Generate `.agent/promotion-review.html` covering all candidates
4. **STOP** — present HTML review, await per-candidate approval
5. **On approval** — write each approved permanent artifact
6. **Clean up** — archive or delete `.agent/` files for this cycle

## HTML Review Contract

Before writing final Markdown artifacts, generate a reviewable HTML file in `.agent/` and pause for approval.

**Workflow Order (MANDATORY):**
1. Generate HTML review file in `.agent/`
2. Present HTML to user for review
3. STOP and wait for explicit approval
4. Only after approval: write final Markdown artifacts

**File Naming Convention:**
- All `.agent/` artifacts MUST use lowercase filenames
- Examples: `increment.md`, `plan.md`, `implementation.md`, `promote.md`
- Review files: `increment-review.html`, `plan-review.html`, `implementation-review.html`, `promotion-report.html`

Required report sections:
1. Objective
2. Inputs Reviewed
3. Proposed Output Summary
4. Risks and Trade-offs
5. Open Questions
6. Approval Decision

HTML requirements:
- Human-readable headings and table(s) where useful.
- Include a timestamp and phase name.
- Include a clear line: `Status: Pending Approval` until approved.
- Use a two-column layout with a left sidebar for quick section navigation.
- Sidebar must contain anchor links to all required report sections.
- Apply a pleasing no-fuzz CSS theme:
	- clear spacing scale and typography hierarchy
	- strong contrast and legible colors
	- simple surfaces and borders without noisy effects
	- mobile-friendly responsive behavior
- Include code highlighting support for snippets:
	- provide semantic classes for tokens (`kw`, `str`, `fn`, `cm`, `id`)
	- style `pre` and `code` blocks for readability
- Support inline SVG rendering for diagrams when useful:
	- allow dedicated diagram sections with embedded `<svg>`
	- style SVG text, lines, and nodes for visual consistency with the theme
- Do not write final Markdown artifacts until approval is explicit.

Canonical HTML skeleton (recommended):

```html
<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Phase Review</title>
	<style>
		:root {
			--bg: #f4f3ef;
			--panel: #ffffff;
			--ink: #111111;
			--muted: #4d4c45;
			--line: #1b1b1b;
			--accent: #234642;
			--codeBg: #151515;
			--codeInk: #ececec;
			--kw: #f4c95d;
			--str: #8dd694;
			--fn: #8cb4ff;
			--cm: #9aa0a6;
			--id: #ff9f7a;
		}
		* { box-sizing: border-box; }
		body {
			margin: 0;
			background: var(--bg);
			color: var(--ink);
			font-family: "IBM Plex Sans", "Segoe UI", -apple-system, sans-serif;
			line-height: 1.45;
		}
		.layout {
			max-width: 1180px;
			margin: 0 auto;
			display: grid;
			grid-template-columns: 250px minmax(0, 1fr);
			gap: 14px;
			padding: 16px;
		}
		.sidebar {
			position: sticky;
			top: 12px;
			align-self: start;
			border: 2px solid var(--line);
			background: var(--panel);
			padding: 12px;
		}
		.sidebar h2 {
			margin: 0 0 8px;
			font-size: 0.98rem;
			letter-spacing: 0.02em;
		}
		.sidebar a {
			display: block;
			color: var(--accent);
			text-decoration: none;
			margin: 7px 0;
			font-weight: 700;
		}
		.content .card {
			border: 2px solid var(--line);
			background: var(--panel);
			padding: 14px;
			margin-bottom: 12px;
		}
		.status {
			display: inline-block;
			padding: 6px 10px;
			border: 2px solid var(--line);
			font-weight: 800;
			background: #fff6dd;
		}
		table {
			width: 100%;
			border-collapse: collapse;
			margin-top: 10px;
		}
		th, td {
			border: 1px solid #2f2f2f;
			text-align: left;
			padding: 8px;
			vertical-align: top;
		}
		pre {
			margin: 10px 0 0;
			padding: 12px;
			border: 1px solid #2d2d2d;
			background: var(--codeBg);
			color: var(--codeInk);
			overflow-x: auto;
		}
		code {
			font-family: "JetBrains Mono", Menlo, Monaco, monospace;
			font-size: 0.92rem;
		}
		.kw { color: var(--kw); }
		.str { color: var(--str); }
		.fn { color: var(--fn); }
		.cm { color: var(--cm); }
		.id { color: var(--id); }
		.diagram svg {
			width: 100%;
			height: auto;
			border: 1px solid #2f2f2f;
			background: #f9f8f4;
		}
		.diagram text {
			fill: #111111;
			font-size: 12px;
			font-family: "IBM Plex Sans", "Segoe UI", sans-serif;
		}
		.diagram .node {
			fill: #ffffff;
			stroke: #1b1b1b;
			stroke-width: 2;
		}
		.diagram .edge {
			stroke: #1b1b1b;
			stroke-width: 2;
			fill: none;
			marker-end: url(#arrow);
		}
		@media (max-width: 920px) {
			.layout { grid-template-columns: 1fr; }
			.sidebar {
				position: static;
				margin-bottom: 8px;
			}
		}
	</style>
</head>
<body>
	<div class="layout">
		<nav class="sidebar" aria-label="Review Sections">
			<h2>Navigate</h2>
			<a href="#objective">Objective</a>
			<a href="#inputs-reviewed">Inputs Reviewed</a>
			<a href="#proposed-output-summary">Proposed Output Summary</a>
			<a href="#risks-and-trade-offs">Risks and Trade-offs</a>
			<a href="#open-questions">Open Questions</a>
			<a href="#approval-decision">Approval Decision</a>
		</nav>

		<main class="content">
			<section class="card" id="objective">
				<h1>Phase Review</h1>
				<p>Phase: Implement | Generated: YYYY-MM-DD HH:MM UTC</p>
				<p><span class="status">Status: Pending Approval</span></p>
			</section>

			<section class="card" id="inputs-reviewed">
				<h2>Inputs Reviewed</h2>
				<ul>
					<li>CONSTITUTION.md</li>
					<li>.agent/increment.md</li>
					<li>.agent/plan.md</li>
					<li>.agent/implementation.md</li>
					<li>.agent/learnings.md</li>
				</ul>
			</section>

			<section class="card" id="proposed-output-summary">
				<h2>Proposed Output Summary</h2>
				<pre><code><span class="kw">const</span> <span class="id">status</span> = <span class="str">"Pending Approval"</span>;
<span class="cm">// Example token classes for code highlighting</span></code></pre>
			</section>

			<section class="card" id="risks-and-trade-offs">
				<h2>Risks and Trade-offs</h2>
				<table>
					<thead>
						<tr><th>Risk</th><th>Trade-off</th><th>Mitigation</th></tr>
					</thead>
					<tbody>
						<tr><td>Example</td><td>Example</td><td>Example</td></tr>
					</tbody>
				</table>
			</section>

			<section class="card diagram" id="open-questions">
				<h2>Open Questions</h2>
				<svg viewBox="0 0 520 140" role="img" aria-label="Optional flow diagram">
					<defs>
						<marker id="arrow" markerWidth="10" markerHeight="7" refX="9" refY="3.5" orient="auto">
							<polygon points="0 0, 10 3.5, 0 7" fill="#1b1b1b"></polygon>
						</marker>
					</defs>
					<rect class="node" x="20" y="35" width="130" height="54"></rect>
					<text x="45" y="67">Input</text>
					<path class="edge" d="M150,62 L250,62"></path>
					<rect class="node" x="250" y="35" width="130" height="54"></rect>
					<text x="273" y="67">Review</text>
					<path class="edge" d="M380,62 L500,62"></path>
					<rect class="node" x="500" y="35" width="0" height="0" style="display:none"></rect>
				</svg>
				<p>Replace this sample SVG with any relevant architecture or flow graph when needed.</p>
			</section>

			<section class="card" id="approval-decision">
				<h2>Approval Decision</h2>
				<p>Do not write final Markdown artifacts until explicit approval.</p>
			</section>
		</main>
	</div>
</body>
</html>
```

---

## Promotion Categories

| Type | Trigger | Destination |
|------|---------|-------------|
| Architecture decision | Non-obvious choice with lasting impact | `docs/adr/ADR-<date>-<slug>.md` |
| Guardrail update | Constitution rule violated, needs clarification | `CONSTITUTION.md` |
| Behavior change | Public API, CLI, or user-facing behavior changed | `README.md` |
| Feature shipped | Acceptance tests pass; feature complete | `docs/roadmap.md` — move to Done, add acceptance test link |
| Test pattern | New testing approach worth standardizing | `CONSTITUTION.md` testing section |
| Known issue | Found but not fixed this cycle | `docs/known-issues.md` |
| New domain concept | An aggregate, event, value object, or term used in code/tests that has no shared definition | `docs/domain-model.md` (create using the template in the Appendix if absent) — place in the correct section: **Aggregates**, **Events**, **Value Objects**, or **Terms** |
| Structural change | A container added, removed, or re-wired | `docs/architecture.md` (create using the template in the Appendix if absent) |

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

---

## Appendix: Document Templates

Use these verbatim as the starting content when creating a new document for the first time.

### Template: docs/domain-model.md

```markdown
# Domain Model

Living record of this project's domain concepts. Updated during the `promote` phase each cycle.

> Source of truth for shared language. Entries are typed — each concept is classified as an **Aggregate**, **Event**, **Value Object**, or **Term**. If a concept is not here, its meaning is ambiguous — add it.

---

## Format

### Aggregates

Aggregates are consistency boundaries — a cluster of entities and value objects with a single root that enforces all invariants.

**AggregateName** — One-sentence description of what this aggregate governs.
- *Root entity:* The entity through which all access must go.
- *Invariants:* Rules enforced within this boundary.
- *Raises:* Events this aggregate produces.
- *Added:* `[increment slug]`

### Events

Domain events record something that happened in the domain — past tense, immutable facts.

**EventName** — One-sentence description of what occurred and why it matters.
- *Payload:* Key data carried by this event.
- *Trigger:* What causes this event to be raised.
- *Consumers:* Who reacts to this event.
- *Added:* `[increment slug]`

### Value Objects

Value objects represent a domain concept defined entirely by its attributes — no identity, immutable.

**ValueObjectName** — One-sentence definition of what it represents.
- *Attributes:* Fields that define equality.
- *Invariants:* Rules that must always hold.
- *Related:* Other concepts this depends on.
- *Added:* `[increment slug]`

### Terms

Ubiquitous language terms that appear in code, tests, and conversations but are not events or value objects.

**Term** — One-sentence definition in domain context.
- *Example:* How it appears in code, CLI output, or tests.
- *Related:* Other terms this depends on or contrasts with.
- *Added:* `[increment slug]`

---

## Aggregates

### [ExampleAggregate]

**[ExampleAggregate]** — [One-sentence description of what consistency boundary this aggregate owns].
- *Root entity:* `[RootEntityName]`
- *Invariants:* [e.g. total must never exceed limit, state transitions must follow X]
- *Raises:* [ExampleEventName]
- *Added:* `[increment slug]`

---

## Events

### [ExampleEventName]

**[ExampleEventName]** — [One-sentence description of the fact that occurred].
- *Payload:* `[field1]`, `[field2]`
- *Trigger:* [What user action or system condition raises this]
- *Consumers:* [Who or what reacts]
- *Added:* `[increment slug]`

---

## Value Objects

### [ExampleValueObject]

**[ExampleValueObject]** — [One-sentence definition].
- *Attributes:* `[attribute1]`, `[attribute2]`
- *Invariants:* [e.g. must be non-negative, must match pattern X]
- *Related:* [ExampleEventName]
- *Added:* `[increment slug]`

---

## Terms

### [ExampleTerm]

**[ExampleTerm]** — [One-sentence definition grounded in this project's domain, not a generic definition].
- *Example:* `[concrete example from code, CLI, or test]`
- *Related:* [ExampleValueObject]
- *Added:* `[increment slug]`

---

## Rules

- Entries must be placed in the correct section: Aggregates, Events, Value Objects, or Terms.
- Definitions are written from the perspective of this project's domain, not general software engineering.
- When a term's meaning changes, update the definition in place — do not add a new entry.
- If two entries turn out to mean the same thing, pick one and add a redirect: `**Alias** — See [Canonical Name].`
- Remove entries that are no longer part of the codebase.
```

### Template: docs/architecture.md

```markdown
# Architecture — [Project Name]

C4 Level 2: Container diagram. Updated when structural boundaries change.

> This is not a design doc. It answers one question: what are the runtime containers, what do they do, and how do they communicate?

---

## Context (C4 Level 1 summary)

**System:** [Project Name]
**Users:** [Who uses it — one line each]
**Purpose:** [One sentence: what problem does this system solve?]

---

## Containers

A container is any separately runnable or deployable unit: a process, a script, a service, a database.

| Container | Technology | Responsibility |
|-----------|-----------|----------------|
| [Container A] | [e.g. Go binary, bash script, Node.js process] | [What it does — one line] |
| [Container B] | [technology] | [responsibility] |
| [Container C] | [technology] | [responsibility] |

---

## Container Diagram

```
┌─────────────────────────────────────────────────────────┐
│  [System Name]                                          │
│                                                         │
│  ┌──────────────────┐         ┌──────────────────────┐  │
│  │  [Container A]   │──────▶  │  [Container B]       │  │
│  │                  │  [how]  │                      │  │
│  │  [technology]    │         │  [technology]        │  │
│  └──────────────────┘         └──────────────────────┘  │
│           │                             │                │
│           ▼                             ▼                │
│  ┌──────────────────┐         ┌──────────────────────┐  │
│  │  [Container C]   │         │  [External System]   │  │
│  │  [technology]    │         │  (out of scope)      │  │
│  └──────────────────┘         └──────────────────────┘  │
└─────────────────────────────────────────────────────────┘

         [User / Actor]
              │
              ▼ [interaction description]
         [Container A]
```

---

## Communication

| From | To | Protocol / Mechanism | Notes |
|------|----|----------------------|-------|
| [Container A] | [Container B] | [e.g. stdout pipe, HTTP, file read] | [any constraint] |
| [User] | [Container A] | [e.g. CLI args, browser] | |

---

## Data Stores

| Store | Type | Owned by | Schema / Format |
|-------|------|----------|-----------------|
| [Store A] | [e.g. SQLite file, CSV, in-memory] | [Container A] | [brief description] |

If no persistent data store exists, state that explicitly:
> This system is stateless. No persistent data store.

---

## Key Constraints

Constraints that affect all containers and must not be violated:

- [e.g. "No network calls — runs entirely offline"]
- [e.g. "Single binary, no install step"]
- [e.g. "All state is held in browser memory; nothing is written to a server"]

---

## Out of Scope

Explicitly name what this diagram does NOT cover:

- Internal component structure of each container (C4 Level 3 — not written unless needed)
- Deployment topology
- CI/CD pipeline

---

## Update Policy

Update this file when:
- A container is added, removed, or its technology changes
- A communication path between containers changes
- A new external system dependency is added

Do NOT update for internal refactors, new features within an existing container, or test changes.

**Last updated:** [YYYY-MM-DD] — [brief reason]
```