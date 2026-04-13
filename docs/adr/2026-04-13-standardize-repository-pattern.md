# Standardise Repository Pattern with Domain Package and Spring Data Naming

**Date**: 2026-04-13  
**Status**: Accepted

## Context

The initial repository layer had several inconsistencies that made the codebase harder to read and test:

1. **Mixed types across layers** — Each repository owned its own model types (`repository.Account`, `repository.Checkin`, `repository.CoupleSummary`). Service and handler code imported these types from the `repository` package, coupling them to the storage layer.

2. **Multiple sentinel errors** — Each repository defined its own not-found sentinel (`ErrAccountNotFound`, `ErrCheckinNotFound`, `ErrCodeNotFound`). Test mocks had to import the `repository` package just to reference these values.

3. **Inconsistent method naming** — Methods mixed styles (`GetCode`, `FindAccountByCode`, `IsAccountPaired`, `CreateCouple`, `GetCoupleSummary`) with no shared convention. The `CoupleRepository.GetCoupleSummary` return signature included a boolean found-flag alongside the error (`(CoupleSummary, bool, error)`), diverging from the rest of the codebase.

4. **Boolean found-flag anti-pattern** — Returning a `bool` alongside an error for "not found" forces callers to handle three independent return values. The idiomatic Go approach is to signal absence via a sentinel error.

## Decision

### 1. Domain package

All model types are moved to `internal/domain`. No other package defines domain types.

```
internal/domain/domain.go
  var ErrNotFound = errors.New("not found")
  type InviteCode string
  type Account struct { ... }
  type Checkin struct { ... }
  type CoupleSummary struct { PartnerFirstName string; FormedOn time.Time }
```

Repository implementations import `domain`; services and handlers import only `domain` (not `repository`) for type references.

### 2. Single sentinel error

`domain.ErrNotFound` replaces all per-entity sentinels. Repository methods return this error when a row is absent. Callers use `errors.Is(err, domain.ErrNotFound)`.

### 3. Spring Data naming conventions

Repository method names follow Spring Data Repository conventions:

| Pattern | Example |
|---|---|
| `FindBy<Field>` | `FindByEmail`, `FindByInviteCode`, `FindByAccountID` |
| `FindBy<Entity>And<Field>` | `FindByAccountAndDate` |
| `Find<Entity>By<Field>` | `FindInviteCodeByAccountID` |
| `Save` | `Save` (insert/upsert) |
| `Save<Entity>` | `SaveInviteCode` |
| `ExistsBy<Field>` | `ExistsCoupleByAccountID` |

### 4. Boolean found-flag removed

`CoupleRepository.GetCoupleSummary` (formerly returning `(CoupleSummary, bool, error)`) is replaced by `FindByAccountID` returning `(domain.CoupleSummary, error)`. Absence is signalled by `domain.ErrNotFound`.

## Consequences

- **Positive**: Service and handler tests no longer import `repository`; mocks depend only on `domain` and `service`.
- **Positive**: A single `domain.ErrNotFound` is callable with `errors.Is` anywhere in the stack.
- **Positive**: Spring Data naming makes method intent unambiguous without reading documentation.
- **Negative**: Future repository methods must follow the naming map; ad-hoc names require a deliberate decision to extend the map.
- **Neutral**: `InviteCode` is a typed string (`type InviteCode string`), not a struct — it carries no metadata beyond the string value.
