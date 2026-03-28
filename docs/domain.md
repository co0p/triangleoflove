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

## Bounded Contexts

### Auth
- **Responsibility**: Verifies Credentials, issues Tokens, enforces that only valid Token holders reach Protected Resources, and serves Profile data.
- **Key concepts**: Account, Credentials, Token, Identity, Profile, Login, Protected Resource
- **Relationships**: Sole owner of the `accounts` table; upstream identity provider to the Home context.

### Home
- **Responsibility**: Displays the Dashboard to an authenticated Account, holds the Token in browser storage, and routes the Account to the login screen when no valid Token exists.
- **Key concepts**: Dashboard, Token, Profile
- **Relationships**: Downstream consumer of Auth — calls `POST /api/v1/auth/login` to obtain a Token and `GET /api/v1/users/me` to retrieve Profile data for the greeting.

## Aggregates

### Account
Protects the invariant that a login attempt only produces a Token when the supplied Credentials match the stored identity. Owns email, hashed password, and first name.

Stored in the `accounts` table (id UUID, email, hashed_password, first_name).

## Value Objects

- **Credentials**: Email + plain-text password supplied at login. Never persisted.
- **Email**: A validated email address. Compared by value; used as the lookup key for an Account.
- **HashedPassword**: The stored, hashed form of the Account's password. Never returned outside the Auth boundary.
- **Token**: A signed JWT string with an issuance timestamp. The only Auth artifact that crosses into the Home context.

## Domain Events

- **AccountAuthenticated**: Emitted when Credentials are valid. Carries the Token.
- **AuthenticationFailed**: Emitted when Credentials do not match. Carries no Token.
