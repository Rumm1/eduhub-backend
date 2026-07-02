package platformdashboard

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetDashboard(ctx context.Context) (DashboardResponse, error) {
	var result DashboardResponse

	err := r.db.QueryRow(ctx, `
SELECT
(SELECT COUNT(*) FROM organizations),
(SELECT COUNT(*) FROM users),
(SELECT COUNT(*) FROM branches),
(SELECT COUNT(*) FROM students),
(SELECT COUNT(*) FROM teacher_profiles),
(SELECT COUNT(*) FROM groups)
`).Scan(
		&result.OrganizationsCount,
		&result.UsersCount,
		&result.BranchesCount,
		&result.StudentsCount,
		&result.TeachersCount,
		&result.GroupsCount,
	)
	if err != nil {
		return DashboardResponse{}, err
	}

	rows, err := r.db.Query(ctx, `
SELECT
o.id::text,
o.name,
COALESCE(o.bin, ''),
COALESCE(o.status, ''),
COALESCE(o.logo_path, ''),
(SELECT COUNT(*) FROM user_profiles up WHERE up.organization_id = o.id),
(SELECT COUNT(*) FROM branches b WHERE b.organization_id = o.id),
(SELECT COUNT(*) FROM students s WHERE s.organization_id = o.id),
(SELECT COUNT(*) FROM teacher_profiles tp WHERE tp.organization_id = o.id),
(SELECT COUNT(*) FROM groups g WHERE g.organization_id = o.id)
FROM organizations o
ORDER BY o.created_at DESC
`)
	if err != nil {
		return DashboardResponse{}, err
	}
	defer rows.Close()

	organizations := make([]OrganizationSummary, 0)

	for rows.Next() {
		var item OrganizationSummary

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.BIN,
			&item.Status,
			&item.LogoPath,
			&item.UsersCount,
			&item.BranchesCount,
			&item.StudentsCount,
			&item.TeachersCount,
			&item.GroupsCount,
		); err != nil {
			return DashboardResponse{}, err
		}

		item.LogoURL = toPublicURL(item.LogoPath)

		branches, err := r.getOrganizationBranches(ctx, item.ID)
		if err != nil {
			return DashboardResponse{}, err
		}

		item.Branches = branches
		organizations = append(organizations, item)
	}

	if err := rows.Err(); err != nil {
		return DashboardResponse{}, err
	}

	result.Organizations = organizations

	return result, nil
}

func (r *Repository) getOrganizationBranches(ctx context.Context, organizationID string) ([]BranchSummary, error) {
	rows, err := r.db.Query(ctx, `
SELECT
id::text,
name,
COALESCE(address, ''),
COALESCE(phone, ''),
COALESCE(status, '')
FROM branches
WHERE organization_id = $1
ORDER BY name ASC
`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	branches := make([]BranchSummary, 0)

	for rows.Next() {
		var item BranchSummary

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Address,
			&item.Phone,
			&item.Status,
		); err != nil {
			return nil, err
		}

		branches = append(branches, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return branches, nil
}

func toPublicURL(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}

	if strings.HasPrefix(path, "/") {
		return path
	}

	return "/" + path
}
