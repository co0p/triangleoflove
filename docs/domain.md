# Domain Model

## Ubiquitous Language

| Term | Definition |
|------|------------|
| Account | A person with an account in the system, identified by email and password |
| Credentials | The email and password pair an Account submits to prove their identity. Password stored hashed; never persisted in plain text |
| Token | A signed JWT issued by the backend after successful login, stored in the browser to authenticate subsequent requests |
| Identity | The verified claim that a request comes from a known Account, carried by a Token |
| Profile | The displayable attributes of an Account (first name, email) and the entry point for account management actions (logout and password change). |
| Logout | The act of removing the Token from browser storage, ending the local session. The backend issues no token invalidation — JWTs are stateless. Client-side only. |
| Password Change | The act of an authenticated Account replacing their password by supplying their current password and a new one. Requires proof of current Credentials. Endpoint: `PUT /api/v1/auth/password`. |
| Protected Resource | Any backend endpoint that requires a valid Token to respond |
| Login | The act of submitting Credentials and receiving a Token |
| Visitor | A person who is not signed in and is trying to create an Account or reach sign-in. |
| Registration | The self-service act of a Visitor creating an Account by supplying first name, email address, password, and password confirmation, then returning to sign-in after success. |
| Password Confirmation | The second password entry used during Registration to confirm the Visitor intended the submitted password. |
| Password Rule | A single requirement a registration password must satisfy before account creation is accepted. Current rules: minimum length and at least one non-alphanumeric character. |
| Validation Signal | Immediate feedback shown while a Visitor types, indicating whether a registration rule is currently satisfied. |
| Validation Error | Clear blocking feedback shown when Registration cannot proceed because a required rule failed. |
| Registration Success Message | The visible confirmation shown on the sign-in page after successful Registration. |
| Sign-in Entry Point | The visible path from the sign-in page to the registration form. |
| Role | A fixed label assigned to an Account at creation time, governing permissions. Two values: `user` and `admin`. |
| Active Account | An Account whose `is_active` flag is `true`; can log in and receive a Token. |
| Inactive Account | An Account whose `is_active` flag is `false`; login is rejected regardless of credential correctness. |
| Activation | The admin action of setting an Account's `is_active` to `true`. |
| Deactivation | The admin action of setting an Account's `is_active` to `false`, immediately blocking login. |
| Rate Limit | A per-IP throttle on Registration and Login endpoints; exceeded requests return HTTP 429 until the window clears. |
| Dashboard | The home screen shown to an authenticated Account, displaying a personalised greeting |
| Check-in | A daily record of an Account's relational wellbeing, comprising seven Ratings (six Relationship Metrics and one Mood) and an optional Note, scoped to a single calendar date. One per Account per date. |
| Dimension | One of seven named aspects being rated: six Relationship Metrics (two per triangle side) and one personal Mood context. Fixed enumeration; not user-configurable. |
| Relationship Metric | One of the six research-grounded proxy questions corresponding to an Intimacy, Commitment, or Passion dimension. |
| Intimacy Metric | Either of the two Relationship Metrics measuring perceived emotional closeness: felt_understood and meaningful_sharing. |
| Commitment Metric | Either of the two Relationship Metrics measuring perceived reliability and investment: could_count_on_them and effort_for_us. |
| Passion Metric | Either of the two Relationship Metrics measuring desire and excitement: desire and spark. |
| Mood | A personal context rating — "How is your overall mood today?" — recorded alongside the six Relationship Metrics but kept visually separate. Not a Relationship Metric. |
| Rating | An integer value representing how strongly an Account felt a Dimension on a given date. Valid set values are 1–5; 0 means the Rating has not been deliberately set (Unset). Plain `int`; no pointer. |
| Unset Rating | A Rating with value 0; displayed with a distinct visual style to signal no deliberate choice has been made. |
| Note | An optional free-text observation attached to a Check-in. Always private to the Account; never shared. Stored as an empty string when not entered. |
| Check-in Date | The calendar date (UTC) for which a Check-in is recorded; defaults to today when opening the page. One Check-in per Account per date. |
| Save Check-in | The action of persisting a Check-in (create or update) for the current Check-in Date. Lazy: no record is written until this action is taken. Upsert semantics — no separate create vs update exposed in the UI or API. |
| Invite Code | A 6-character uppercase alphanumeric code generated for an Account, used to initiate pairing. Ephemeral — replaced after a successful pairing or explicit regeneration. Stored on the Account. |
| Couple | A bond between exactly two Accounts, carrying a formation date and optionally an end date. An Active Couple has no end date; an Ended Couple does. |
| Pairing | The act of two Accounts forming a Couple via Invite Code exchange. One-time per Couple formation. |
| Paired | The state of an Account that belongs to an Active Couple. Derived from Couple membership. |
| Active Couple | A Couple whose `ended_on` is null. Determines the Paired state for both members. Derived from the `couples` table; not a separate record type. |
| Ended Couple | A Couple whose `ended_on` is set. Both members are considered unpaired but the record is retained for relational history. |
| Unpair | The act of one Account ending an Active Couple. Unilateral — no partner consent required. Writes `ended_on` to the Couple record; does not delete it. |
| ErrNotFound | A single shared sentinel error (`domain.ErrNotFound`) returned by any repository method when the requested record does not exist. Replaces per-entity variants. Callers use `errors.Is(err, domain.ErrNotFound)` to detect absence. |
| Daily Insight | Normalized dimension scores derived from a single Check-in. Each of Intimacy, Commitment, and Passion is expressed as an integer 0–100, or -1 if both proxy metrics were unset. Computed server-side; never persisted. |
| Weekly Insight | A single day's insight scores within the Weekly Insights Window. Carries a date (YYYYMMDD) and Intimacy, Commitment, and Passion scores. A score of -1 indicates no Check-in exists for that day and dimension. |
| Weekly Insights Window | The rolling 7-day period from 6 days ago through yesterday (UTC), used for the weekly matrix view. Determined by server clock; index 0 = oldest day, index 6 = yesterday. |

## Bounded Contexts

### Auth
- **Responsibility**: Verifies Credentials, issues Tokens carrying role claims, handles Registration, evaluates credential policy for email format and password rules, enforces that Inactive Accounts cannot log in, serves Profile data, and allows authenticated Accounts to change their own password.
- **Key concepts**: Visitor, Account, Registration, Password Confirmation, Password Rule, Validation Signal, Validation Error, Registration Success Message, Sign-in Entry Point, Role, Active Account, Inactive Account, Credentials, Token, Identity, Profile, Login, Protected Resource, Password Change, Rate Limit, Logout
- **Relationships**: Sole owner of the `accounts` table; upstream identity provider to other contexts.

### Admin
- **Responsibility**: Enables Accounts with role `admin` to list all Accounts and toggle activation state. Authorization reads role from JWT claims.
- **Key concepts**: Role, Active Account, Inactive Account, Activation, Deactivation, AccountSummary
- **Relationships**: Downstream of Auth for identity and role claims. Permitted to write only `is_active` on the `accounts` table.

### Home
- **Responsibility**: Displays the Dashboard, holds the Token in browser storage, routes to login when no Token exists, and hosts the Profile page — the account management entry point navigable from the NavBar.
- **Key concepts**: Dashboard, Token, Profile, Logout
- **Relationships**: Downstream consumer of Auth — calls `POST /api/v1/auth/login` to obtain a Token and `GET /api/v1/users/me` to retrieve Profile data. Hosts the Profile page which calls `PUT /api/v1/auth/password`.

### Check-in
- **Responsibility**: Owns the lifecycle of a Check-in — recording, loading, and updating Ratings and Notes for an Account on a given date.
- **Key concepts**: Check-in, Dimension, Relationship Metric, Intimacy Metric, Commitment Metric, Passion Metric, Mood, Rating, Unset Rating, Note, Check-in Date, Save Check-in
- **Relationships**: Downstream of Auth — requires a valid Token (Account identity) to record or retrieve a Check-in. Entered from the Home context.

### Pairing
- **Responsibility**: Manages Invite Code generation and regeneration; forms and soft-ends Couples between Accounts; serves Couple status (partner first name, paired-since date, and active/ended state) to authenticated Accounts.
- **Key concepts**: Invite Code, Couple, Active Couple, Ended Couple, Pairing, Paired, Unpair
- **Relationships**: Reads Account identity from Auth. Owns the `couples` table (including the `ended_on` column) and the `invite_code` column on `accounts`.

### Insights
- **Responsibility**: Derives and serves normalized insight scores from Check-in data — both daily (single date) and weekly (7-day rolling window). Owns no persistent data of its own; reads from the Check-in aggregate's table.
- **Key concepts**: Daily Insight, Weekly Insight, Weekly Insights Window
- **Relationships**: Downstream of Auth — requires a valid Token. Reads check-in rows via InsightsRepository; does not depend on CheckinService.

## Aggregates

### Account
Protects three invariants — (1) a login attempt only produces a Token when supplied Credentials match the stored identity and the Account is Active; (2) `changePassword` only succeeds when the supplied current password matches the stored `hashedPassword`; (3) Registration is rejected when email is already taken, email format is invalid, password is shorter than 8 characters, or password lacks a non-alphanumeric character. Password confirmation remains a registration-time interaction rule and is not persisted on the Account. Attempts with a non-matching current password are rejected with `ErrInvalidCredentials`.

Stored in the `accounts` table (id UUID, email, hashed_password, first_name, role, is_active, created_at).

### CheckIn
Protects the invariant that at most one Check-in exists per Account per date. Owns seven Ratings (six Relationship Metrics — felt_understood, meaningful_sharing, could_count_on_them, effort_for_us, desire, spark — plus Mood) and the Note. Natural identity key is `(accountId, date)`.

Stored in the `checkins` table. Written lazily — no row exists until Save Check-in is triggered.

### Couple
Protects two invariants — (1) a Couple can only be formed between two distinct Accounts that are not already in an Active Couple; (2) `Unpair` can only be applied to an Active Couple (`ended_on IS NULL`). Attempting to unpair an Ended Couple is rejected.

Stored in the `couples` table (id UUID, account_id_a, account_id_b, formed_on TIMESTAMPTZ, ended_on TIMESTAMPTZ NULL). The Invite Code is stored as `invite_code` on the `accounts` table — it is ephemeral and belongs to the Account, not the Couple.

## Value Objects

- **Credentials**: Email + plain-text password supplied at login. Never persisted.
- **Email**: A validated email address. Compared by value; used as the lookup key for an Account.
- **HashedPassword**: The stored, hashed form of the Account's password. Never returned outside the Auth boundary.
- **PlainPassword**: The raw password string supplied in a `PUT /api/v1/auth/password` request — for current-password verification or as the new-password candidate. Never persisted; discarded after use.
- **Token**: A signed JWT string with an issuance timestamp. The only Auth artifact that crosses into the Home context.
- **Role**: A constrained string value object with two valid values: `user` and `admin`. Assigned at account creation.
- **Rating**: An integer. Valid set values are 1–5. `0` means the user has not yet set a value (Unset Rating). Plain `int`; no pointer. Compared by value; no identity.
- **Note**: A string, may be empty. No identity. Compared by value.
- **InviteCode**: A 6-character uppercase alphanumeric string. Implemented as `type InviteCode string` — a typed string to prevent raw string substitution at repository boundaries. Generated randomly; no identity. Compared by value. Replaced (not mutated) on pairing or explicit regeneration.
- **AccountSummary**: A read model for the admin user list. Holds `id`, `email`, `firstName`, `role`, `isActive`, and `createdAt`.
- **CoupleSummary**: A read model produced by the Couple aggregate. Holds `PartnerFirstName` (string) and `FormedOn` (UTC timestamp). No identity; compared by value. Never mutated — always replaced by a fresh query.
- **DailyInsight**: A read model computed from one Check-in. Holds Intimacy, Commitment, and Passion as 0–100 integers, or -1 when both proxy metrics are unset. Never persisted; produced by `domain.NewDailyInsight`.
- **WeeklyInsight**: A read model for one day in the Weekly Insights Window. Holds a date string (YYYYMMDD) and three dimension scores (0–100 or -1). Produced per day in the 7-day window; days with no Check-in yield all -1 scores.

## Domain Events

- **AccountAuthenticated**: Emitted when Credentials are valid. Carries the Token.
- **AccountRegistered**: Emitted when Registration succeeds.
- **PasswordChanged**: Emitted when `changePassword` succeeds. Carries `accountID` and `changedAt`.
- **AuthenticationFailed**: Emitted when Credentials do not match. Carries no Token.
- **AccountActivated**: Emitted when an admin sets `is_active` to `true`.
- **AccountDeactivated**: Emitted when an admin sets `is_active` to `false`.
- **CheckInSaved**: Emitted when a Check-in is created or updated. Carries `accountId`, `date`, and `savedAt`.
- **CoupleFormed**: Emitted when two Accounts successfully pair. Carries `accountIdA`, `accountIdB`, `formedOn`.
- **CoupleEnded**: Emitted when `Unpair` sets `ended_on`. Carries `coupleID`, `accountIDA`, `accountIDB`, and `endedOn`.
