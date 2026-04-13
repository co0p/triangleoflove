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
| Dashboard | The home screen shown to an authenticated Account, displaying a personalised greeting |
| Check-in | A daily record of an Account's relational wellbeing, comprising five Ratings and an optional Note, scoped to a single calendar date. One per Account per date. |
| Dimension | One of the five named relational aspects being rated: "Felt close today", "Positive energy / fun", "Supported / team", "Communication healthy", "My stress level". Fixed enumeration; not user-configurable. |
| Rating | An integer value representing how strongly an Account felt a Dimension on a given date. Valid saved values are 1–10; -1 means the Rating has not been deliberately set. |
| Unset Rating | A Rating with value -1; displayed at slider position 5 with a distinct visual style to signal no deliberate choice has been made. Replaced by the actual slider value on first Save Check-in. |
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

## Bounded Contexts

### Auth
- **Responsibility**: Verifies Credentials, issues Tokens, enforces that only valid Token holders reach Protected Resources, serves Profile data, and allows authenticated Accounts to change their own password.
- **Key concepts**: Account, Credentials, Token, Identity, Profile, Login, Protected Resource, Password Change, Logout
- **Relationships**: Sole owner of the `accounts` table; upstream identity provider to the Home context.

### Home
- **Responsibility**: Displays the Dashboard, holds the Token in browser storage, routes to login when no Token exists, and hosts the Profile page — the account management entry point navigable from the NavBar.
- **Key concepts**: Dashboard, Token, Profile, Logout
- **Relationships**: Downstream consumer of Auth — calls `POST /api/v1/auth/login` to obtain a Token and `GET /api/v1/users/me` to retrieve Profile data. Hosts the Profile page which calls `PUT /api/v1/auth/password`.

### Check-in
- **Responsibility**: Owns the lifecycle of a Check-in — recording, loading, and updating Ratings and Notes for an Account on a given date.
- **Key concepts**: Check-in, Dimension, Rating, Unset Rating, Note, Check-in Date, Save Check-in
- **Relationships**: Downstream of Auth — requires a valid Token (Account identity) to record or retrieve a Check-in. Entered from the Home context.

### Pairing
- **Responsibility**: Manages Invite Code generation and regeneration; forms and soft-ends Couples between Accounts; serves Couple status (partner first name, paired-since date, and active/ended state) to authenticated Accounts.
- **Key concepts**: Invite Code, Couple, Active Couple, Ended Couple, Pairing, Paired, Unpair
- **Relationships**: Reads Account identity from Auth. Owns the `couples` table (including the `ended_on` column) and the `invite_code` column on `accounts`.

## Aggregates

### Account
Protects two invariants — (1) a login attempt only produces a Token when the supplied Credentials match the stored identity; (2) `changePassword` only succeeds when the supplied current password matches the stored `hashedPassword`. Attempts with a non-matching current password are rejected with `ErrInvalidCredentials`. Owns email, hashed password, and first name; exposes `SaveHashedPassword` as the sole mutation for password updates.

Stored in the `accounts` table (id UUID, email, hashed_password, first_name).

### CheckIn
Protects the invariant that at most one Check-in exists per Account per date. Owns all five Ratings and the Note. Natural identity key is `(accountId, date)`.

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
- **Rating**: An integer. Valid saved values are 1–10. `-1` means the user has not yet set a value (Unset Rating). Compared by value; no identity.
- **Note**: A string, may be empty. No identity. Compared by value.
- **InviteCode**: A 6-character uppercase alphanumeric string. Implemented as `type InviteCode string` — a typed string to prevent raw string substitution at repository boundaries. Generated randomly; no identity. Compared by value. Replaced (not mutated) on pairing or explicit regeneration.
- **CoupleSummary**: A read model produced by the Couple aggregate. Holds `PartnerFirstName` (string) and `FormedOn` (UTC timestamp). No identity; compared by value. Never mutated — always replaced by a fresh query.

## Domain Events

- **AccountAuthenticated**: Emitted when Credentials are valid. Carries the Token.
- **PasswordChanged**: Emitted when `changePassword` succeeds. Carries `accountID` and `changedAt`.
- **AuthenticationFailed**: Emitted when Credentials do not match. Carries no Token.
- **CheckInSaved**: Emitted when a Check-in is created or updated. Carries `accountId`, `date`, and `savedAt`.
- **CoupleFormed**: Emitted when two Accounts successfully pair. Carries `accountIdA`, `accountIdB`, `formedOn`.
- **CoupleEnded**: Emitted when `Unpair` sets `ended_on`. Carries `coupleID`, `accountIDA`, `accountIDB`, and `endedOn`.
