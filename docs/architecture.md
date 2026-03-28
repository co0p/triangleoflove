# Architecture

## System Context

```mermaid
C4Context
    title Triangle of Love — System Context
    Person(account, "Account", "A pre-seeded user of the coaching app")
    System(tol, "Triangle of Love", "Relationship coaching web app — login and personalised dashboard")
    Rel(account, tol, "Logs in and views dashboard")
```

## Container Diagram

```mermaid
C4Container
    title Triangle of Love — Container Diagram
    Person(account, "Account")
    Container(frontend, "Frontend", "Vue 3 / Vite", "Serves the login screen and dashboard; stores Token in browser")
    Container(backend, "Backend", "Go / net/http", "Issues Tokens, validates identity, serves profile data")
    ContainerDb(db, "Database", "PostgreSQL", "Stores accounts (credentials + first name)")
    Rel(account, frontend, "Uses", "HTTPS")
    Rel(frontend, backend, "POST /api/v1/auth/login", "JSON/HTTPS")
    Rel(frontend, backend, "GET /api/v1/users/me", "JSON/HTTPS + Bearer Token")
    Rel(backend, db, "Reads account record", "SQL")
```

## Backend Component Diagram

```mermaid
C4Component
    title Backend — Component Diagram (Auth increment)
    Container(frontend, "Frontend", "Vue 3", "Calls API")
    Component(authHandler, "Auth Handler", "Go / net/http", "Handles POST /api/v1/auth/login")
    Component(meHandler, "Me Handler", "Go / net/http", "Handles GET /api/v1/users/me")
    Component(jwtMiddleware, "JWT Middleware", "Go", "Validates Bearer Token; rejects with 401 if absent or invalid")
    Component(authService, "Auth Service", "Go", "Orchestrates login: looks up Account, verifies Credentials, issues Token")
    Component(accountRepo, "Account Repository", "Go / sql", "Reads Account record from Postgres by email")
    ContainerDb(db, "Database", "PostgreSQL", "accounts table")
    Rel(frontend, authHandler, "POST /api/v1/auth/login")
    Rel(frontend, jwtMiddleware, "GET /api/v1/users/me + Bearer Token")
    Rel(jwtMiddleware, meHandler, "Passes verified identity")
    Rel(authHandler, authService, "login(credentials)")
    Rel(meHandler, authService, "getProfile(accountId)")
    Rel(authService, accountRepo, "findByEmail / findById")
    Rel(accountRepo, db, "SQL query")
```
