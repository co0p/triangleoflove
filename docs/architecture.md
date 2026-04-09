# Architecture

## System Context

```mermaid
C4Context
    title Triangle of Love — System Context
    Person(account, "Account", "A pre-seeded user of the coaching app")
    System(tol, "Triangle of Love", "Relationship coaching web app — login, dashboard, daily check-in, and partner pairing")
    Rel(account, tol, "Logs in, views dashboard, records daily check-in, and pairs with a partner")
```

## Container Diagram

```mermaid
C4Container
    title Triangle of Love — Container Diagram
    Person(account, "Account")
    Container(frontend, "Frontend", "Vue 3 / Vite", "Serves login, dashboard, check-in, and pairing pages; stores Token in browser")
    Container(backend, "Backend", "Go / net/http", "Issues Tokens, validates identity, serves profile, check-in, and pairing data")
    ContainerDb(db, "Database", "PostgreSQL", "Stores accounts, check-ins, and couples")
    Rel(account, frontend, "Uses", "HTTPS")
    Rel(frontend, backend, "POST /api/v1/auth/login", "JSON/HTTPS")
    Rel(frontend, backend, "GET /api/v1/users/me", "JSON/HTTPS + Bearer Token")
    Rel(frontend, backend, "GET /api/v1/checkins/today", "JSON/HTTPS + Bearer Token")
    Rel(frontend, backend, "PUT /api/v1/checkins/today", "JSON/HTTPS + Bearer Token")
    Rel(frontend, backend, "GET /api/v1/pairing", "JSON/HTTPS + Bearer Token")
    Rel(frontend, backend, "POST /api/v1/pairing/regenerate", "JSON/HTTPS + Bearer Token")
    Rel(frontend, backend, "POST /api/v1/pairing/connect", "JSON/HTTPS + Bearer Token")
    Rel(frontend, backend, "GET /api/v1/couples/me", "JSON/HTTPS + Bearer Token")
    Rel(backend, db, "Reads/writes accounts, check-ins, and couples", "SQL")
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
    Component(pairingHandler, "Pairing Handler", "Go / net/http", "Handles GET /api/v1/pairing, POST /api/v1/pairing/regenerate, POST /api/v1/pairing/connect, GET /api/v1/couples/me")
    Component(pairingService, "Pairing Service", "Go", "Orchestrates invite code generation and Couple formation")
    Component(pairingRepo, "Pairing Repository", "Go / sql", "Reads and writes invite_code on the accounts table")
    Component(coupleRepo, "Couple Repository", "Go / sql", "Reads and writes the couples table")
    ContainerDb(db, "Database", "PostgreSQL", "accounts, checkins, and couples tables")
    Rel(frontend, authHandler, "POST /api/v1/auth/login")
    Rel(frontend, jwtMiddleware, "GET /api/v1/users/me + Bearer Token")
    Rel(frontend, jwtMiddleware, "GET|PUT /api/v1/checkins/today + Bearer Token")
    Rel(frontend, jwtMiddleware, "GET /api/v1/pairing + Bearer Token")
    Rel(frontend, jwtMiddleware, "POST /api/v1/pairing/* + Bearer Token")
    Rel(frontend, jwtMiddleware, "GET /api/v1/couples/me + Bearer Token")
    Rel(jwtMiddleware, meHandler, "Passes verified identity")
    Rel(jwtMiddleware, checkinHandler, "Passes verified accountId")
    Rel(jwtMiddleware, pairingHandler, "Passes verified accountId")
    Rel(authHandler, authService, "login(credentials)")
    Rel(meHandler, authService, "getProfile(accountId)")
    Rel(authService, accountRepo, "findByEmail / findById")
    Rel(accountRepo, db, "SQL query")
    Rel(checkinHandler, checkinService, "getToday(accountId) / save(accountId, checkIn)")
    Rel(checkinService, checkinRepo, "findByAccountAndDate / upsert")
    Rel(checkinRepo, db, "SQL select / upsert")
    Rel(pairingHandler, pairingService, "getOrCreateCode / regenerateCode / connect / getCoupleStatus")
    Rel(pairingService, pairingRepo, "getCode / setCode / findByCode / isPaired")
    Rel(pairingService, coupleRepo, "createCouple / getCoupleSummary")
    Rel(pairingRepo, db, "SQL — accounts.invite_code")
    Rel(coupleRepo, db, "SQL — couples table")
```
