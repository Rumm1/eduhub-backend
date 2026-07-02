package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID         uuid.UUID   `json:"user_id"`
	ProfileID      *uuid.UUID  `json:"profile_id,omitempty"`
	OrganizationID *uuid.UUID  `json:"organization_id,omitempty"`
	Roles          []string    `json:"roles"`
	Permissions    []string    `json:"permissions"`
	BranchIDs      []uuid.UUID `json:"branch_ids"`
	TokenType      string      `json:"token_type"`

	jwt.RegisteredClaims
}
