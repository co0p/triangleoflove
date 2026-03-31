# ADR: -1 as the Unset Rating Sentinel

## Context

Each of the five Rating columns in the `checkins` table must represent three states: a deliberate value of 1–10, and "the user has not yet set this". The page loads with sliders that show a mid-point position but signal "no deliberate choice made"—this distinction is meaningful at the domain level and must survive a round-trip through the API and database.

## Decision

Store `-1` in each Rating column to represent an unset value. The DB schema enforces `CHECK (value = -1 OR value BETWEEN 1 AND 10)`. The API returns `-1` as-is; no normalisation to `null`. The frontend treats any value of `-1` as Unset Rating and renders the slider with a distinct visual style (`checkin-slider--unset`).

## Consequences

- **Benefits:** Single representation of "unset" from DB through API to UI — no translation step. `-1` is explicit and unambiguous in a field that otherwise holds 1–10. CHECK constraint enforces the invariant at the database level.
- **Drawbacks:** Consumers that do not read this decision may treat `-1` as a low score rather than "not set".
- **Trade-offs:** The meaning of `-1` must be documented in the API contract and respected by all clients — currently only the first-party frontend.

## Alternatives Considered

- **`null` in the DB column and API response**: Requires nullable SMALLINT columns and null-handling at every layer; `-1` already names the concept in the ubiquitous language, making `null` a redundant second representation.
- **Separate boolean column `ratings_unset`**: Adds a column that must stay in sync with the Rating values; creates two representations of the same invariant.
