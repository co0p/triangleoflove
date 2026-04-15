package domain

import (
	"errors"
	"time"
)

// ErrNotFound is returned by any repository method when the requested record does not exist.
var ErrNotFound = errors.New("not found")

// InviteCode is a typed string representing a 6-character uppercase alphanumeric pairing code.
type InviteCode string

// Account holds the core identity fields for an authenticated user.
type Account struct {
	ID             string
	Email          string
	HashedPassword string
	FirstName      string
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
