package file

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

func (r *Repository) List(ctx context.Context, organizationID uuid.UUID) ([]File, error) {
	rows, err := r.db.Query(ctx, `
SELECT
id,
COALESCE(organization_id, '00000000-0000-0000-0000-000000000000'::uuid),
COALESCE(uploaded_by, '00000000-0000-0000-0000-000000000000'::uuid),
COALESCE(folder, ''),
file_name,
file_path,
COALESCE(mime_type, ''),
COALESCE(size_bytes, 0),
created_at::text
FROM files
WHERE organization_id = $1
ORDER BY created_at DESC
`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]File, 0)

	for rows.Next() {
		var item File

		if err := rows.Scan(
			&item.ID,
			&item.OrganizationID,
			&item.UploadedBy,
			&item.Folder,
			&item.FileName,
			&item.FilePath,
			&item.MimeType,
			&item.SizeBytes,
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

func (r *Repository) GetByID(ctx context.Context, organizationID uuid.UUID, fileID uuid.UUID) (File, error) {
	var item File

	err := r.db.QueryRow(ctx, `
SELECT
id,
COALESCE(organization_id, '00000000-0000-0000-0000-000000000000'::uuid),
COALESCE(uploaded_by, '00000000-0000-0000-0000-000000000000'::uuid),
COALESCE(folder, ''),
file_name,
file_path,
COALESCE(mime_type, ''),
COALESCE(size_bytes, 0),
created_at::text
FROM files
WHERE id = $1
  AND organization_id = $2
`, fileID, organizationID).Scan(
		&item.ID,
		&item.OrganizationID,
		&item.UploadedBy,
		&item.Folder,
		&item.FileName,
		&item.FilePath,
		&item.MimeType,
		&item.SizeBytes,
		&item.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return File{}, ErrFileNotFound
		}

		return File{}, err
	}

	return item, nil
}

func (r *Repository) Create(ctx context.Context, item File) (File, error) {
	err := r.db.QueryRow(ctx, `
INSERT INTO files (
id,
organization_id,
uploaded_by,
folder,
file_name,
file_path,
mime_type,
size_bytes
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
id,
organization_id,
uploaded_by,
COALESCE(folder, ''),
file_name,
file_path,
COALESCE(mime_type, ''),
COALESCE(size_bytes, 0),
created_at::text
`,
		item.ID,
		item.OrganizationID,
		item.UploadedBy,
		item.Folder,
		item.FileName,
		item.FilePath,
		item.MimeType,
		item.SizeBytes,
	).Scan(
		&item.ID,
		&item.OrganizationID,
		&item.UploadedBy,
		&item.Folder,
		&item.FileName,
		&item.FilePath,
		&item.MimeType,
		&item.SizeBytes,
		&item.CreatedAt,
	)
	if err != nil {
		return File{}, err
	}

	return item, nil
}

func (r *Repository) Delete(ctx context.Context, organizationID uuid.UUID, fileID uuid.UUID) (string, error) {
	var filePath string

	err := r.db.QueryRow(ctx, `
DELETE FROM files
WHERE id = $1
  AND organization_id = $2
RETURNING file_path
`, fileID, organizationID).Scan(&filePath)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", ErrFileNotFound
		}

		return "", err
	}

	return filePath, nil
}
