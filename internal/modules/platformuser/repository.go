package platformuser

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

type UserCredentialsTarget struct {
	ID    uuid.UUID
	Email string
}

func (r *Repository) GetUserByID(ctx context.Context, userID uuid.UUID) (UserCredentialsTarget, error) {
	var user UserCredentialsTarget

	err := r.db.QueryRow(ctx, `
SELECT id, email
FROM users
WHERE id = $1
`, userID).Scan(
		&user.ID,
		&user.Email,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserCredentialsTarget{}, ErrUserNotFound
		}

		return UserCredentialsTarget{}, err
	}

	return user, nil
}

func (r *Repository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	result, err := r.db.Exec(ctx, `
UPDATE users
SET password_hash = $2,
    must_change_password = TRUE
WHERE id = $1
`, userID, passwordHash)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}
