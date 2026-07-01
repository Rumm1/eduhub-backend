package jwt

import (
	"fmt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const AccessTokenType = "access"

type Manager struct {
	accessSecret string
	accessTTL    time.Duration
}

type AccessTokenPayload struct {
	UserID         uuid.UUID
	OrganizationID *uuid.UUID
	Roles          []string
	Permissions    []string
	BranchIDs      []uuid.UUID
}

func NewManager(accessSecret string, accessTTLMinutes int) *Manager {
	if accessTTLMinutes <= 0 {
		accessTTLMinutes = 15
	}

	return &Manager{
		accessSecret: accessSecret,
		accessTTL:    time.Duration(accessTTLMinutes) * time.Minute,
	}
}

func (m *Manager) GenerateAccessToken(payload AccessTokenPayload) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID:         payload.UserID,
		OrganizationID: payload.OrganizationID,
		Roles:          payload.Roles,
		Permissions:    payload.Permissions,
		BranchIDs:      payload.BranchIDs,
		TokenType:      AccessTokenType,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(now.Add(m.accessTTL)),
			IssuedAt:  gojwt.NewNumericDate(now),
			Subject:   payload.UserID.String(),
		},
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.accessSecret))
}

func (m *Manager) ParseAccessToken(tokenString string) (*Claims, error) {
	token, err := gojwt.ParseWithClaims(tokenString, &Claims{}, func(token *gojwt.Token) (interface{}, error) {
		if token.Method != gojwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(m.accessSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.TokenType != AccessTokenType {
		return nil, fmt.Errorf("invalid token type")
	}

	return claims, nil
}
