# Domain Model

## Ubiquitous Language

| Term | Definition |
|------|------------|
| Account | A person with an account in the system, identified by email and password |
| Credentials | The email and password pair an Account submits to prove identity. Password is stored hashed and never persisted in plain text |
| Token | A signed JWT issued by the backend after successful login, stored in the browser to authenticate subsequent requests |
| Identity | The verified claim that a request comes from a known Account, carried by a Token |
| Profile | The displayable attributes of an Account (first name, email) and the entry point for account management actions (logout and password change) |
| Logout | The act of removing the Token from browser storage, ending the local session. The backend issues no token invalidation (JWTs are stateless) |
| Password Change | The act of an authenticated Account replacing their password by supplying current and new password. Endpoint: `PUT /api/v1/auth/password` |
| Protected Resource | Any backend endpoint that requires a valid Token to respond |
| Login | The act of submitting Credentials and receiving a Token |
| Visitor | A person who is not signed in and is trying to create an Account or reach sign-in |
| Registration | The self-service act of a Visitor creating an Account by supplying first name, email, password, and password confirmation |
| Password Confirmation | The second password entry used during Registration to confirm intended password |
| Password Rule | A requirement a registration password must satisfy. Current rules: minimum length and at least one non-alphanumeric character |
| Validation Signal | Immediate feedback shown while a Visitor types, indicating whether a registration rule is currently satisfied |
| Validation Error | Clear blocking feedback shown when Registration cannot proceed because a required rule failed |
| Registration Success Message | The visible confirmation shown on the sign-in page after successful Registration |
| Sign-in Entry Point | The visible path from the sign-in page to the registration form |
| Role | A fixed label assigned to an Account at creation time, governing permissions (`user` or `admin`) |
| Active Account | An Account whose `is_active` flag is `true`; can log in and receive a Token |
| Inactive Account | An Account whose `is_active` flag is `false`; login is rejected regardless of credential correctness |
| Activation | The admin action of setting an Account's `is_active` to `true` |
| Deactivation | The admin action of setting an Account's `is_active` to `false`, immediately blocking login |
| Rate Limit | A per-IP throttle on Registration and Login endpoints; exceeded requests return HTTP 429 until the window clears |
| Dashboard | The authenticated home screen that displays a personalized greeting and entry points to check-ins, pairing, and insights |
| Check-in | A daily record of relational wellbeing comprising seven Ratings (six relationship metrics and one Mood) and an optional Note, scoped to one calendar date. One per Account per date |
| Dimension | One of seven named aspects being rated: six Relationship Metrics (two per triangle side) plus one personal Mood context |
| Relationship Metric | One of six proxy questions corresponding to an Intimacy, Commitment, or Passion dimension |
| Intimacy Metric | Either `felt_understood` or `meaningful_sharing` |
| Commitment Metric | Either `could_count_on_them` or `effort_for_us` |
| Passion Metric | Either `desire` or `spark` |
| Mood | A personal context rating recorded alongside the six Relationship Metrics. Not a Relationship Metric |
| Rating | An integer representing how strongly an Account felt a Dimension on a given date. Valid set values are 1-5; 0 means Unset |
| Unset Rating | A Rating with value 0, signaling no deliberate choice has been made |
| Note | An optional free-text observation attached to a Check-in. Private to the Account and never shared |
| Check-in Date | The UTC calendar date for which a Check-in is recorded; defaults to today when opening the page |
| Save Check-in | The action of persisting a Check-in (create or update) for the current Check-in Date using upsert semantics |
| Invite Code | A 6-character uppercase alphanumeric code generated for an Account to initiate pairing |
| Couple | A bond between exactly two Accounts, carrying a formation date and optionally an end date |
| Pairing | The act of two Accounts forming a Couple via Invite Code exchange |
| Paired | The state of an Account that belongs to an Active Couple |
| Active Couple | A Couple whose `ended_on` is null |
| Ended Couple | A Couple whose `ended_on` is set; both members are considered unpaired but the record is retained |
| Unpair | The act of one Account ending an Active Couple by setting `ended_on` |
| ErrNotFound | Shared sentinel error (`domain.ErrNotFound`) returned by repository methods when a record does not exist |
| Daily Insight | Normalized dimension scores derived from a single Check-in. Intimacy, Commitment, and Passion are 0-100, or -1 when unavailable |
| Weekly Insight | A single day's insight scores within the 7-day window. Date is `YYYYMMDD`; scores are 0-100 or -1 if no check-in exists |
| Weekly Insights Window | Rolling 7-day period from 6 days ago through yesterday (UTC), used for the weekly matrix view |

## Bounded Contexts

### Auth
- **Responsibility**: Verifies Credentials, issues Tokens carrying role claims, handles Registration, evaluates credential policy, enforces that Inactive Accounts cannot log in, serves Profile data, and allows authenticated Accounts to change their own password.
- **Key concepts**: Visitor, Account, Registration, Password Confirmation, Password Rule, Validation Signal, Validation Error, Registration Success Message, Sign-in Entry Point, Role, Active Account, Inactive Account, Credentials, Token, Identity, Profile, Login, Protected Resource, Password Change, Rate Limit, Logout
- **Relationships**: Sole owner of the `accounts` table; upstream identity provider to other contexts.

### Admin
- **Responsibility**: Enables Accounts with role `admin` to list all Accounts and toggle activation state.
- **Key concepts**: Role, Active Account, Inactive Account, Activation, Deactivation, AccountSummary
- **Relationships**: Downstream of Auth for identity and role claims. Permitted to write only `is_active` on `accounts`.

### Home
- **Responsibility**: Displays the Dashboard, stores the Token in browser storage, routes to login when no Token exists, and hosts the Profile page.
- **Key concepts**: Dashboard, Token, Profile, Logout
- **Relationships**: Downstream of Auth. Calls `POST /api/v1/auth/login`, `GET /api/v1/users/me`, and `PUT /api/v1/auth/password`.

### Check-in
- **Responsibility**: Owns the lifecycle of a Check-in (recording, loading, updating Ratings and Note for an Account on a date).
- **Key concepts**: Check-in, Dimension, Relationship Metric, Intimacy Metric, Commitment Metric, Passion Metric, Mood, Rating, Unset Rating, Note, Check-in Date, Save Check-in
- **Relationships**: Downstream of Auth. Requires a valid Token to record or retrieve a Check-in.

### Pairing
- **Responsibility**: Manages Invite Code generation/regeneration, forms and soft-ends Couples, serves couple status to authenticated Accounts.
- **Key concepts**: Invite Code, Couple, Active Couple, Ended Couple, Pairing, Paired, Unpair
- **Relationships**: Reads Account identity from Auth. Owns `couples` and `accounts.invite_code`.

### Insights
- **Responsibility**: Derives and serves normalized insight scores from Check-in data for daily and weekly views.
- **Key concepts**: Daily Insight, Weekly Insight, Weekly Insights Window
- **Relationships**: Downstream of Auth. Reads check-ins via `InsightsRepository`; does not depend on `CheckinService`.

## Aggregates

### Account
Protects three invariants: (1) login produces a Token only when Credentials match and Account is Active; (2) password change succeeds only when current password matches the stored hash; (3) registration rejects duplicate email, invalid email format, short passwords, and passwords lacking a non-alphanumeric character.

Stored in `accounts` (`id`, `email`, `hashed_password`, `first_name`, `role`, `is_active`, `created_at`).

### CheckIn
Protects the invariant that at most one Check-in exists per Account per date. Owns seven Ratings (`felt_understood`, `meaningful_sharing`, `could_count_on_them`, `effort_for_us`, `desire`, `spark`, `mood`) and optional `note`. Natural identity key is `(accountId, date)`.

Stored in `checkins`. Written lazily: no row exists until Save Check-in is triggered.

### Couple
Protects two invariants: (1) a Couple can only form between two distinct Accounts not already in an Active Couple; (2) Unpair applies only to an Active Couple (`ended_on IS NULL`).

Stored in `couples` (`id`, `account_id_a`, `account_id_b`, `formed_on`, `ended_on`). Invite Code is stored on `accounts.invite_code`.

## Value Objects

- **Credentials**: Email + plain-text password supplied at login. Never persisted.
- **Email**: Validated email address used as the Account lookup key.
- **HashedPassword**: Stored hashed form of password. Never returned outside Auth.
- **PlainPassword**: Raw password in `PUT /api/v1/auth/password` request. Never persisted.
- **Token**: Signed JWT string with issuance timestamp.
- **Role**: Constrained string value with valid values `user` and `admin`.
- **Rating**: Integer value object with allowed values 0-5 (`0` means Unset).
- **Note**: Optional string attached to a Check-in.
- **InviteCode**: Typed string (`type InviteCode string`) holding a 6-character uppercase alphanumeric code.
- **AccountSummary**: Admin read model containing `id`, `email`, `firstName`, `role`, `isActive`, `createdAt`.
- **CoupleSummary**: Pairing read model containing `PartnerFirstName` and `FormedOn`.
- **DailyInsight**: Read model with Intimacy/Commitment/Passion (0-100 or -1) from one Check-in.
- **WeeklyInsight**: Read model with date (`YYYYMMDD`) and three scores (0-100 or -1) for one day in the weekly window.

## Domain Events

- **AccountAuthenticated**: Emitted when Credentials are valid; carries Token.
- **AccountRegistered**: Emitted when Registration succeeds.
- **PasswordChanged**: Emitted when password change succeeds.
- **AuthenticationFailed**: Emitted when Credentials do not match.
- **AccountActivated**: Emitted when admin sets `is_active` to `true`.
- **AccountDeactivated**: Emitted when admin sets `is_active` to `false`.
- **CheckInSaved**: Emitted when a Check-in is created or updated; carries `accountId`, `date`, `savedAt`.
- **CoupleFormed**: Emitted when two Accounts successfully pair.
- **CoupleEnded**: Emitted when Unpair sets `ended_on`.
