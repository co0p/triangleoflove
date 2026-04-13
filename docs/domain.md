# Domain Model

## Ubiquitous Language

| Term | Definition |
|------|------------|
| Account | A person with an account in the system, identified by email and password |
| Credentials | The email and password pair an Account submits to prove their identity. Password stored hashed; never persisted in plain text |
| Token | A signed JWT issued by the backend after successful login, stored in the browser to authenticate subsequent requests |
| Identity | The verified claim that a request comes from a known Account, carried by a Token |
| Profile | The displayable attributes of an Account — first name for the initial increment |
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
| Couple | A permanent bond between exactly two Accounts, formed when one Account submits the other's Invite Code. Owns the formation date and future pair-scoped data. |
| Pairing | The act of two Accounts forming a Couple via Invite Code exchange. One-time per Couple formation. |
| Paired | The state of an Account that belongs to a Couple. Derived from Couple membership. |
| ErrNotFound | A single shared sentinel error (`domain.ErrNotFound`) returned by any repository method when the requested record does not exist. Replaces per-entity variants. Callers use `errors.Is(err, domain.ErrNotFound)` to detect absence. |

## Bounded Contexts

### Auth
- **Responsibility**: Verifies Credentials, issues Tokens, enforces that only valid Token holders reach Protected Resources, and serves Profile data.
- **Key concepts**: Account, Credentials, Token, Identity, Profile, Login, Protected Resource
- **Relationships**: Sole owner of the `accounts` table; upstream identity provider to the Home context.

### Home
- **Responsibility**: Displays the Dashboard to an authenticated Account, holds the Token in browser storage, and routes the Account to the login screen when no valid Token exists.
- **Key concepts**: Dashboard, Token, Profile
- **Relationships**: Downstream consumer of Auth — calls `POST /api/v1/auth/login` to obtain a Token and `GET /api/v1/users/me` to retrieve Profile data for the greeting. Links to the Check-in context via navigation.

### Check-in
- **Responsibility**: Owns the lifecycle of a Check-in — recording, loading, and updating Ratings and Notes for an Account on a given date.
- **Key concepts**: Check-in, Dimension, Rating, Unset Rating, Note, Check-in Date, Save Check-in
- **Relationships**: Downstream of Auth — requires a valid Token (Account identity) to record or retrieve a Check-in. Entered from the Home context.

### Pairing
- **Responsibility**: Manages Invite Code generation and regeneration; forms Couples between Accounts; serves Couple status (partner first name, paired-since date) to authenticated Accounts.
- **Key concepts**: Invite Code, Couple, Pairing, Paired
- **Relationships**: Reads Account identity from Auth. Owns the `couples` table and the `invite_code` column on `accounts`.

## Aggregates

### Account
Protects the invariant that a login attempt only produces a Token when the supplied Credentials match the stored identity. Owns email, hashed password, and first name.

Stored in the `accounts` table (id UUID, email, hashed_password, first_name).

### CheckIn
Protects the invariant that at most one Check-in exists per Account per date. Owns all five Ratings and the Note. Natural identity key is `(accountId, date)`.

Stored in the `checkins` table. Written lazily — no row exists until Save Check-in is triggered.

### Couple
Protects the invariant that exactly two distinct Accounts form a Couple and that each Account belongs to at most one Couple. Owns the formation date and will own future pair-scoped fields.

Stored in the `couples` table (id UUID, account_id_a, account_id_b, formed_on TIMESTAMPTZ). The Invite Code is stored as `invite_code` on the `accounts` table — it is ephemeral and belongs to the Account, not the Couple.

## Value Objects

- **Credentials**: Email + plain-text password supplied at login. Never persisted.
- **Email**: A validated email address. Compared by value; used as the lookup key for an Account.
- **HashedPassword**: The stored, hashed form of the Account's password. Never returned outside the Auth boundary.
- **Token**: A signed JWT string with an issuance timestamp. The only Auth artifact that crosses into the Home context.
- **Rating**: An integer. Valid saved values are 1–10. `-1` means the user has not yet set a value (Unset Rating). Compared by value; no identity.
- **Note**: A string, may be empty. No identity. Compared by value.
- **InviteCode**: A 6-character uppercase alphanumeric string. Implemented as `type InviteCode string` — a typed string to prevent raw string substitution at repository boundaries. Generated randomly; no identity. Compared by value. Replaced (not mutated) on pairing or explicit regeneration.
- **CoupleSummary**: A read model produced by the Couple aggregate. Holds `PartnerFirstName` (string) and `FormedOn` (UTC timestamp). No identity; compared by value. Never mutated — always replaced by a fresh query.

## Domain Events

- **AccountAuthenticated**: Emitted when Credentials are valid. Carries the Token.
- **AuthenticationFailed**: Emitted when Credentials do not match. Carries no Token.
- **CheckInSaved**: Emitted when a Check-in is created or updated. Carries `accountId`, `date`, and `savedAt`.
- **CoupleFormed**: Emitted when two Accounts successfully pair. Carries `accountIdA`, `accountIdB`, `formedOn`.
