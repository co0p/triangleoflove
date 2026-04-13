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
type Checkin struct {
	FeltClose            *int   `json:"felt_close"`
	PositiveEnergy       *int   `json:"positive_energy"`
	Supported            *int   `json:"supported"`
	CommunicationHealthy *int   `json:"communication_healthy"`
	StressLevel          *int   `json:"stress_level"`
	Note                 string `json:"note"`
}

// CoupleSummary is a read model describing a paired relationship.
type CoupleSummary struct {
	PartnerFirstName string
	FormedOn         time.Time
}
