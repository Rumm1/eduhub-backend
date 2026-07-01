package jwt

import (
	"testing"

	"github.com/google/uuid"
)

func TestGenerateAndParseAccessToken(t *testing.T) {
	manager := NewManager("test_secret", 15)

	userID := uuid.New()
	orgID := uuid.New()
	branchID := uuid.New()

	token, err := manager.GenerateAccessToken(AccessTokenPayload{
		UserID:         userID,
		OrganizationID: &orgID,
		Roles:          []string{"ORG_ADMIN"},
		Permissions:    []string{"students.read", "students.create"},
		BranchIDs:      []uuid.UUID{branchID},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	claims, err := manager.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("expected no parse error, got %v", err)
	}

	if claims.UserID != userID {
		t.Fatalf("expected user id %s, got %s", userID, claims.UserID)
	}

	if claims.OrganizationID == nil || *claims.OrganizationID != orgID {
		t.Fatalf("expected organization id %s", orgID)
	}

	if len(claims.Roles) != 1 || claims.Roles[0] != "ORG_ADMIN" {
		t.Fatal("expected ORG_ADMIN role")
	}

	if len(claims.Permissions) != 2 {
		t.Fatal("expected 2 permissions")
	}

	if len(claims.BranchIDs) != 1 || claims.BranchIDs[0] != branchID {
		t.Fatal("expected branch id")
	}
}

func TestParseInvalidAccessToken(t *testing.T) {
	manager := NewManager("test_secret", 15)

	_, err := manager.ParseAccessToken("invalid.token.value")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}
