package audit

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, input CreateAuditLogInput) error {
	metadataBytes, err := json.Marshal(input.Metadata)
	if err != nil {
		return err
	}

	var organizationID interface{}
	if strings.TrimSpace(input.OrganizationID) != "" {
		organizationID = input.OrganizationID
	}

	var userID interface{}
	if strings.TrimSpace(input.UserID) != "" {
		userID = input.UserID
	}

	var entityID interface{}
	if strings.TrimSpace(input.EntityID) != "" {
		entityID = input.EntityID
	}

	_, err = r.db.Exec(ctx, `
INSERT INTO audit_logs (
organization_id,
user_id,
action,
entity_type,
entity_id,
description,
metadata,
ip_address,
user_agent
)
VALUES (
$1::uuid,
$2::uuid,
$3,
$4,
$5::uuid,
$6,
$7::jsonb,
$8,
$9
)
`, organizationID, userID, input.Action, input.EntityType, entityID, input.Description, string(metadataBytes), input.IPAddress, input.UserAgent)

	return err
}

func (r *Repository) List(ctx context.Context, filter AuditLogFilter) ([]AuditLog, int, error) {
	whereParts := []string{
		"al.organization_id = $1::uuid",
	}

	args := []interface{}{filter.OrganizationID}
	argIndex := 2

	if filter.UserID != "" {
		whereParts = append(whereParts, "al.user_id = $"+itoa(argIndex)+"::uuid")
		args = append(args, filter.UserID)
		argIndex++
	}

	if filter.Action != "" {
		whereParts = append(whereParts, "al.action = $"+itoa(argIndex))
		args = append(args, filter.Action)
		argIndex++
	}

	if filter.EntityType != "" {
		whereParts = append(whereParts, "al.entity_type = $"+itoa(argIndex))
		args = append(args, filter.EntityType)
		argIndex++
	}

	if filter.EntityID != "" {
		whereParts = append(whereParts, "al.entity_id = $"+itoa(argIndex)+"::uuid")
		args = append(args, filter.EntityID)
		argIndex++
	}

	if filter.FromDate != "" {
		whereParts = append(whereParts, "al.created_at >= $"+itoa(argIndex)+"::date")
		args = append(args, filter.FromDate)
		argIndex++
	}

	if filter.ToDate != "" {
		whereParts = append(whereParts, "al.created_at < ($"+itoa(argIndex)+"::date + interval '1 day')")
		args = append(args, filter.ToDate)
		argIndex++
	}

	whereSQL := strings.Join(whereParts, " AND ")

	countQuery := `
SELECT COUNT(*)
FROM audit_logs al
WHERE ` + whereSQL

	total := 0
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
SELECT
al.id::text,
COALESCE(al.organization_id::text, ''),
COALESCE(al.user_id::text, ''),
COALESCE(u.full_name, ''),
al.action,
al.entity_type,
COALESCE(al.entity_id::text, ''),
COALESCE(al.description, ''),
al.metadata::text,
COALESCE(al.ip_address, ''),
COALESCE(al.user_agent, ''),
al.created_at
FROM audit_logs al
LEFT JOIN users u ON u.id = al.user_id
WHERE ` + whereSQL + `
ORDER BY al.created_at DESC
LIMIT $` + itoa(argIndex) + `
OFFSET $` + itoa(argIndex+1)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]AuditLog, 0)

	for rows.Next() {
		var item AuditLog

		if err := rows.Scan(
			&item.ID,
			&item.OrganizationID,
			&item.UserID,
			&item.UserName,
			&item.Action,
			&item.EntityType,
			&item.EntityID,
			&item.Description,
			&item.Metadata,
			&item.IPAddress,
			&item.UserAgent,
			&item.CreatedAt,
		); err != nil {
			return nil, 0, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func IsValidUUID(value string) bool {
	if strings.TrimSpace(value) == "" {
		return true
	}

	_, err := uuid.Parse(value)
	return err == nil
}

func itoa(value int) string {
	return strconv.Itoa(value)
}
