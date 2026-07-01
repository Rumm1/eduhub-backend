package usercontext

import "context"

type contextKey string

const userKey contextKey = "user"

type User struct {
	ID             string   `json:"id"`
	Role           string   `json:"role,omitempty"`
	OrganizationID string   `json:"organization_id,omitempty"`
	BranchID       string   `json:"branch_id,omitempty"`
	Permissions    []string `json:"permissions,omitempty"`
}

func WithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func UserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userKey).(User)
	return user, ok
}
