package permission

import (
	"context"
	"strings"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context) (ListPermissionsResponse, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return ListPermissionsResponse{}, err
	}

	return ListPermissionsResponse{
		Items: buildPermissionResponses(items),
		Total: len(items),
	}, nil
}

func (s *Service) ListGroups(ctx context.Context) (ListPermissionGroupsResponse, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return ListPermissionGroupsResponse{}, err
	}

	groupMap := make(map[string][]Permission)
	groupOrder := make([]string, 0)

	for _, item := range items {
		groupName := item.Group
		if groupName == "" {
			groupName = "other"
		}

		if _, exists := groupMap[groupName]; !exists {
			groupOrder = append(groupOrder, groupName)
		}

		groupMap[groupName] = append(groupMap[groupName], item)
	}

	groups := make([]PermissionGroupResponse, 0, len(groupOrder))

	for _, groupName := range groupOrder {
		groups = append(groups, PermissionGroupResponse{
			Name:        groupName,
			Title:       permissionGroupTitle(groupName),
			Permissions: buildPermissionResponses(groupMap[groupName]),
		})
	}

	return ListPermissionGroupsResponse{
		Items: groups,
		Total: len(groups),
	}, nil
}

func buildPermissionResponses(items []Permission) []PermissionResponse {
	responses := make([]PermissionResponse, 0, len(items))

	for _, item := range items {
		responses = append(responses, PermissionResponse{
			ID:          item.ID.String(),
			Code:        item.Code,
			Name:        item.Name,
			Description: item.Description,
			Group:       item.Group,
		})
	}

	return responses
}

func permissionGroupFromCode(code string) string {
	code = strings.TrimSpace(code)
	if code == "" {
		return "other"
	}

	parts := strings.Split(code, ".")
	if len(parts) == 0 || parts[0] == "" {
		return "other"
	}

	return parts[0]
}

func permissionGroupTitle(group string) string {
	titles := map[string]string{
		"attendance":    "Attendance",
		"audit_logs":    "Audit Logs",
		"branches":      "Branches",
		"dashboard":     "Dashboard",
		"files":         "Files",
		"groups":        "Groups",
		"homeworks":     "Homeworks",
		"lessons":       "Lessons",
		"notifications": "Notifications",
		"payments":      "Payments",
		"payroll":       "Payroll",
		"profiles":      "Profiles",
		"reports":       "Reports",
		"schedules":     "Schedules",
		"students":      "Students",
		"subjects":      "Subjects",
		"teachers":      "Teachers",
		"users":         "Users",
	}

	if title, ok := titles[group]; ok {
		return title
	}

	if group == "" {
		return "Other"
	}

	return strings.Title(strings.ReplaceAll(group, "_", " "))
}
