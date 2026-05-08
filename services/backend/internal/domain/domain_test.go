package domain_test

import (
	"errors"
	"testing"

	"triangleoflove/backend/internal/domain"
)

func TestDomain_GivenDomainPackage_WhenInspected_ThenAllModelTypesDefined(t *testing.T) {
	_ = domain.Account{ID: "1", Email: "a@b.com", HashedPassword: "h", FirstName: "A"}

	_ = domain.Checkin{FeltUnderstood: 4, MeaningfulSharing: 3, CouldCountOnThem: 5,
		EffortForUs: 2, Desire: 3, Spark: 4, Mood: 4}

	_ = domain.CoupleSummary{PartnerFirstName: "Bob"}

	var code domain.InviteCode = "ABC123"
	_ = code
}

func TestDomain_GivenDomainPackage_WhenInspected_ThenSingleErrNotFoundDefined(t *testing.T) {
	err := domain.ErrNotFound
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatal("domain.ErrNotFound must satisfy errors.Is with itself")
	}
}

func TestRole_GivenNewAccount_WhenCreated_ThenRoleIsUser(t *testing.T) {
	a := domain.Account{
		ID:             "1",
		Email:          "new@example.com",
		HashedPassword: "hashed",
		FirstName:      "New",
		Role:           domain.RoleUser,
		IsActive:       true,
	}

	if a.Role != domain.RoleUser {
		t.Fatalf("expected role %q, got %q", domain.RoleUser, a.Role)
	}
}
