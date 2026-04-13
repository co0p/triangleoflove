package domain_test

import (
	"errors"
	"testing"

	"triangleoflove/backend/internal/domain"
)

func TestDomain_GivenDomainPackage_WhenInspected_ThenAllModelTypesDefined(t *testing.T) {
	_ = domain.Account{ID: "1", Email: "a@b.com", HashedPassword: "h", FirstName: "A"}

	felt := 5
	_ = domain.Checkin{FeltClose: &felt}

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
