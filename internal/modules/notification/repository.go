package notification

import (
	"context"

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

func (r *Repository) ListForUser(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID) ([]Notification, error) {
	rows, err := r.db.Query(ctx, `
SELECT
id,
COALESCE(organization_id, '00000000-0000-0000-0000-000000000000'::uuid),
user_id,
title,
COALESCE(message, ''),
COALESCE(type, ''),
is_read,
created_at::text
FROM notifications
WHERE user_id = $1
  AND (organization_id = $2 OR organization_id IS NULL)
ORDER BY created_at DESC
`, userID, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Notification, 0)

	for rows.Next() {
		var item Notification

		if err := rows.Scan(
			&item.ID,
			&item.OrganizationID,
			&item.UserID,
			&item.Title,
			&item.Message,
			&item.Type,
			&item.IsRead,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) CreateMany(
	ctx context.Context,
	organizationID uuid.UUID,
	userIDs []uuid.UUID,
	title string,
	message string,
	notificationType string,
) ([]Notification, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	items := make([]Notification, 0, len(userIDs))

	for _, userID := range userIDs {
		item := Notification{
			ID:             uuid.New(),
			OrganizationID: organizationID,
			UserID:         userID,
			Title:          title,
			Message:        message,
			Type:           notificationType,
		}

		err := tx.QueryRow(ctx, `
INSERT INTO notifications (
id,
organization_id,
user_id,
title,
message,
type
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
id,
organization_id,
user_id,
title,
COALESCE(message, ''),
COALESCE(type, ''),
is_read,
created_at::text
`,
			item.ID,
			item.OrganizationID,
			item.UserID,
			item.Title,
			item.Message,
			item.Type,
		).Scan(
			&item.ID,
			&item.OrganizationID,
			&item.UserID,
			&item.Title,
			&item.Message,
			&item.Type,
			&item.IsRead,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) UserExistsInOrganization(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID) (bool, error) {
	var exists bool

	err := r.db.QueryRow(ctx, `
SELECT EXISTS (
SELECT 1
FROM user_profiles
WHERE user_id = $1
  AND organization_id = $2
  AND status = 'active'
)
`, userID, organizationID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *Repository) ListOrganizationUserIDs(ctx context.Context, organizationID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx, `
SELECT DISTINCT user_id
FROM user_profiles
WHERE organization_id = $1
  AND status = 'active'
ORDER BY user_id
`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]uuid.UUID, 0)

	for rows.Next() {
		var userID uuid.UUID

		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}

		items = append(items, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) MarkRead(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID, notificationID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `
UPDATE notifications
SET is_read = true
WHERE id = $1
  AND user_id = $2
  AND (organization_id = $3 OR organization_id IS NULL)
`, notificationID, userID, organizationID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNotificationNotFound
	}

	return nil
}

func (r *Repository) MarkAllRead(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID) (int64, error) {
	tag, err := r.db.Exec(ctx, `
UPDATE notifications
SET is_read = true
WHERE user_id = $1
  AND is_read = false
  AND (organization_id = $2 OR organization_id IS NULL)
`, userID, organizationID)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}

func (r *Repository) Delete(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID, notificationID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `
DELETE FROM notifications
WHERE id = $1
  AND user_id = $2
  AND (organization_id = $3 OR organization_id IS NULL)
`, notificationID, userID, organizationID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNotificationNotFound
	}

	return nil
}

func (r *Repository) GetByIDForUser(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID, notificationID uuid.UUID) (Notification, error) {
	var item Notification

	err := r.db.QueryRow(ctx, `
SELECT
id,
COALESCE(organization_id, '00000000-0000-0000-0000-000000000000'::uuid),
user_id,
title,
COALESCE(message, ''),
COALESCE(type, ''),
is_read,
created_at::text
FROM notifications
WHERE id = $1
  AND user_id = $2
  AND (organization_id = $3 OR organization_id IS NULL)
`, notificationID, userID, organizationID).Scan(
		&item.ID,
		&item.OrganizationID,
		&item.UserID,
		&item.Title,
		&item.Message,
		&item.Type,
		&item.IsRead,
		&item.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Notification{}, ErrNotificationNotFound
		}

		return Notification{}, err
	}

	return item, nil
}
