package domain

import (
	"errors"
	"time"
)

// ErrNotFound is returned by any repository method when the requested record does not exist.
var ErrNotFound = errors.New("not found")

// ErrDuplicateEmail is returned by AccountRepository.Register when the email is already taken.
var ErrDuplicateEmail = errors.New("duplicate email")

// Role represents the access level of an Account. It is fixed at creation time.
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// InviteCode is a typed string representing a 6-character uppercase alphanumeric pairing code.
type InviteCode string

// Account holds the core identity fields for an authenticated user.
type Account struct {
	ID             string
	Email          string
	HashedPassword string
	FirstName      string
	Role           Role
	IsActive       bool
	CreatedAt      time.Time
}

// Checkin holds a single daily check-in record.
// Rating fields use 0 as the unset sentinel; valid entered values are 1–5.
type Checkin struct {
	FeltUnderstood    int    `json:"felt_understood"`
	MeaningfulSharing int    `json:"meaningful_sharing"`
	CouldCountOnThem  int    `json:"could_count_on_them"`
	EffortForUs       int    `json:"effort_for_us"`
	Desire            int    `json:"desire"`
	Spark             int    `json:"spark"`
	Mood              int    `json:"mood"`
	Note              string `json:"note"`
}

// CoupleSummary is a read model describing a paired relationship.
type CoupleSummary struct {
	PartnerFirstName string
	FormedOn         time.Time
}

// AccountSummary is a read model projected for the admin user list.
type AccountSummary struct {
	ID        string
	Email     string
	FirstName string
	Role      Role
	IsActive  bool
	CreatedAt time.Time
}
