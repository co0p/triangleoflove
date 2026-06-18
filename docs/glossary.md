# Glossary

## Document Role

Defines terms used across code, tests, and documentation whose meaning is specific to this project. Authoritative source is `docs/domain.md` — entries here are promoted from that document when a term appears in code or tests without a shared definition. When definitions diverge, `docs/domain.md` wins.

---

## Terms

| Term | Definition |
|------|------------|
| Account | A person with an account in the system, identified by email and password. |
| Active Account | An Account whose `is_active` flag is `true`; can log in and receive a Token. |
| Active Couple | A Couple whose `ended_on` is null. Determines the Paired state for both members. |
| Activation | The admin action of setting an Account's `is_active` to `true`. |
| Check-in | A daily record of an Account's relational wellbeing, comprising seven Ratings and an optional Note, scoped to a single calendar date. One per Account per date. |
| Check-in Date | The calendar date (UTC) for which a Check-in is recorded; defaults to today. |
| Commitment Metric | Either of the two Relationship Metrics measuring perceived reliability and investment: `could_count_on_them` and `effort_for_us`. |
| Couple | A bond between exactly two Accounts, carrying a formation date and optionally an end date. |
| Credentials | The email and password pair an Account submits to prove their identity. |
| Daily Insight | Normalized dimension scores derived from a single Check-in. Each of Intimacy, Commitment, and Passion is expressed as an integer 0–100, or -1 if both proxy metrics were unset. |
| Dashboard | The home screen shown to an authenticated Account. |
| Deactivation | The admin action of setting an Account's `is_active` to `false`, immediately blocking login. |
| Dimension | One of seven named aspects being rated: six Relationship Metrics and one Mood. |
| Ended Couple | A Couple whose `ended_on` is set. Both members are considered unpaired; the record is retained for history. |
| ErrNotFound | The single shared sentinel error (`domain.ErrNotFound`) returned by any repository method when a requested record does not exist. |
| Identity | The verified claim that a request comes from a known Account, carried by a Token. |
| Inactive Account | An Account whose `is_active` flag is `false`; login is rejected regardless of credential correctness. |
| Intimacy Metric | Either of the two Relationship Metrics measuring perceived emotional closeness: `felt_understood` and `meaningful_sharing`. |
| Invite Code | A 6-character uppercase alphanumeric code generated for an Account, used to initiate pairing. |
| Login | The act of submitting Credentials and receiving a Token. |
| Logout | The act of removing the Token from browser storage, ending the local session. Client-side only. |
| Mood | A personal context rating alongside the six Relationship Metrics. Not a Relationship Metric. |
| Note | An optional free-text observation attached to a Check-in. Always private; never shared. |
| Paired | The state of an Account that belongs to an Active Couple. |
| Pairing | The act of two Accounts forming a Couple via Invite Code exchange. |
| Passion Metric | Either of the two Relationship Metrics measuring desire and excitement: `desire` and `spark`. |
| Password Change | The act of an authenticated Account replacing their password by supplying their current password and a new one. |
| Password Confirmation | The second password entry used during Registration to confirm the intended password. |
| Password Rule | A single requirement a registration password must satisfy. Current rules: minimum length and at least one non-alphanumeric character. |
| Profile | The displayable attributes of an Account (first name, email) and the entry point for account management actions. |
| Protected Resource | Any backend endpoint that requires a valid Token to respond. |
| Rate Limit | A per-IP throttle on Registration and Login endpoints; exceeded requests return HTTP 429. |
| Rating | An integer value (1–5) representing how strongly an Account felt a Dimension. 0 means Unset. |
| Registration | The self-service act of a Visitor creating an Account. |
| Relationship Metric | One of the six research-grounded proxy questions corresponding to an Intimacy, Commitment, or Passion dimension. |
| Role | A fixed label assigned to an Account at creation time: `user` or `admin`. |
| Save Check-in | The action of persisting a Check-in (create or update) for the current Check-in Date. Upsert semantics. |
| Token | A signed JWT issued by the backend after successful login. |
| Unpair | The act of one Account ending an Active Couple. Unilateral. Writes `ended_on` to the Couple record. |
| Unset Rating | A Rating with value 0; displayed with a distinct visual style to signal no deliberate choice was made. |
| Visitor | A person who is not signed in and is trying to create an Account or reach sign-in. |
| Weekly Insight | A single day's insight scores within the Weekly Insights Window. |
| Weekly Insights Window | The rolling 7-day period from 6 days ago through yesterday (UTC), used for the weekly matrix view. |
