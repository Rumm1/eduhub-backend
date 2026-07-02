package context

import (
	stdcontext "context"

	"github.com/google/uuid"
)

type contextKey string

const userContextKey contextKey = "user_context"

type UserContext struct {
	UserID         uuid.UUID
	ProfileID      *uuid.UUID
	OrganizationID *uuid.UUID
	Roles          []string
	Permissions    []string
	BranchIDs      []uuid.UUID
}

func WithUser(ctx stdcontext.Context, user UserContext) stdcontext.Context {
	return stdcontext.WithValue(ctx, userContextKey, user)
}

func GetUser(ctx stdcontext.Context) (UserContext, bool) {
	user, ok := ctx.Value(userContextKey).(UserContext)
	return user, ok
}

func MustGetUser(ctx stdcontext.Context) UserContext {
	user, ok := GetUser(ctx)
	if !ok {
		panic("user context not found")
	}

	return user
}

func HasPermission(ctx stdcontext.Context, permission string) bool {
	user, ok := GetUser(ctx)
	if !ok {
		return false
	}

	for _, item := range user.Permissions {
		if item == permission {
			return true
		}
	}

	return false
}

func HasRole(ctx stdcontext.Context, role string) bool {
	user, ok := GetUser(ctx)
	if !ok {
		return false
	}

	for _, item := range user.Roles {
		if item == role {
			return true
		}
	}

	return false
}
