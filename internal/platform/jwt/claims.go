package jwt

import "time"

type Claims struct {
	Subject        string   `json:"sub"`
	Role           string   `json:"role,omitempty"`
	OrganizationID string   `json:"organization_id,omitempty"`
	BranchID       string   `json:"branch_id,omitempty"`
	Permissions    []string `json:"permissions,omitempty"`
	IssuedAt       int64    `json:"iat"`
	ExpiresAt      int64    `json:"exp"`
}

func (c Claims) Valid(now time.Time) bool {
	return c.Subject != "" && c.ExpiresAt > now.Unix()
}
