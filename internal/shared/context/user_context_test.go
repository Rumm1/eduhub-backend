package context

import (
	stdcontext "context"
	"testing"

	"github.com/google/uuid"
)

func TestUserContext(t *testing.T) {
	userID := uuid.New()
	orgID := uuid.New()
	branchID := uuid.New()

	user := UserContext{
		UserID:         userID,
		OrganizationID: &orgID,
		Roles:          []string{"ORG_ADMIN"},
		Permissions:    []string{"students.read"},
		BranchIDs:      []uuid.UUID{branchID},
	}

	ctx := WithUser(stdcontext.Background(), user)

	got, ok := GetUser(ctx)
	if !ok {
		t.Fatal("expected user context")
	}

	if got.UserID != userID {
		t.Fatal("wrong user id")
	}

	if got.OrganizationID == nil || *got.OrganizationID != orgID {
		t.Fatal("wrong organization id")
	}

	if !HasRole(ctx, "ORG_ADMIN") {
		t.Fatal("expected ORG_ADMIN role")
	}

	if !HasPermission(ctx, "students.read") {
		t.Fatal("expected students.read permission")
	}

	if HasPermission(ctx, "students.delete") {
		t.Fatal("did not expect students.delete permission")
	}
}

func TestGetUserNotFound(t *testing.T) {
	_, ok := GetUser(stdcontext.Background())
	if ok {
		t.Fatal("expected user context not found")
	}
}
