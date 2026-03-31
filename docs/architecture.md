# Architecture

## System Context

```mermaid
C4Context
    title Triangle of Love — System Context
    Person(account, "Account", "A pre-seeded user of the coaching app")
    System(tol, "Triangle of Love", "Relationship coaching web app — login, dashboard, and daily check-in")
    Rel(account, tol, "Logs in, views dashboard, records daily check-in")
```

## Container Diagram

```mermaid
C4Container
    title Triangle of Love — Container Diagram
    Person(account, "Account")
    Container(frontend, "Frontend", "Vue 3 / Vite", "Serves login, dashboard, and check-in page; stores Token in browser")
    Container(backend, "Backend", "Go / net/http", "Issues Tokens, validates identity, serves profile and check-in data")
    ContainerDb(db, "Database", "PostgreSQL", "Stores accounts and check-ins")
    Rel(account, frontend, "Uses", "HTTPS")
    Rel(frontend, backend, "POST /api/v1/auth/login", "JSON/HTTPS")
    Rel(frontend, backend, "GET /api/v1/users/me", "JSON/HTTPS + Bearer Token")
    Rel(frontend, backend, "GET /api/v1/checkins/today", "JSON/HTTPS + Bearer Token")
    Rel(frontend, backend, "PUT /api/v1/checkins/today", "JSON/HTTPS + Bearer Token")
    Rel(backend, db, "Reads/writes accounts and checkins", "SQL")
```

## Backend Component Diagram

```mermaid
C4Component
    title Backend — Component Diagram
    Container(frontend, "Frontend", "Vue 3", "Calls API")
    Component(authHandler, "Auth Handler", "Go / net/http", "Handles POST /api/v1/auth/login")
    Component(meHandler, "Me Handler", "Go / net/http", "Handles GET /api/v1/users/me")
    Component(jwtMiddleware, "JWT Middleware", "Go", "Validates Bearer Token; rejects with 401 if absent or invalid")
    Component(authService, "Auth Service", "Go", "Orchestrates login: looks up Account, verifies Credentials, issues Token")
    Component(accountRepo, "Account Repository", "Go / sql", "Reads Account record from Postgres by email")
    Component(checkinHandler, "Check-in Handler", "Go / net/http", "Handles GET and PUT /api/v1/checkins/today")
    Component(checkinService, "Check-in Service", "Go", "Orchestrates load and upsert of a Check-in for today")
    Component(checkinRepo, "Check-in Repository", "Go / sql", "Reads and upserts rows in the checkins table")
    ContainerDb(db, "Database", "PostgreSQL", "accounts and checkins tables")
    Rel(frontend, authHandler, "POST /api/v1/auth/login")
    Rel(frontend, jwtMiddleware, "GET /api/v1/users/me + Bearer Token")
    Rel(frontend, jwtMiddleware, "GET|PUT /api/v1/checkins/today + Bearer Token")
    Rel(jwtMiddleware, meHandler, "Passes verified identity")
    Rel(jwtMiddleware, checkinHandler, "Passes verified accountId")
    Rel(authHandler, authService, "login(credentials)")
    Rel(meHandler, authService, "getProfile(accountId)")
    Rel(authService, accountRepo, "findByEmail / findById")
    Rel(accountRepo, db, "SQL query")
    Rel(checkinHandler, checkinService, "getToday(accountId) / save(accountId, checkIn)")
    Rel(checkinService, checkinRepo, "findByAccountAndDate / upsert")
    Rel(checkinRepo, db, "SQL select / upsert")
```
