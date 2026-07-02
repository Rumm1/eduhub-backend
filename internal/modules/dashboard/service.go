package dashboard

import (
	"context"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/google/uuid"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetOverview(ctx context.Context) (OverviewResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return OverviewResponse{}, ErrTenantRequired
	}

	overview, err := s.repository.GetOverview(ctx, *currentUser.OrganizationID, currentUser.UserID)
	if err != nil {
		return OverviewResponse{}, err
	}

	return mapOverview(overview), nil
}

func mapOverview(item Overview) OverviewResponse {
	return OverviewResponse{
		StudentsCount:           item.StudentsCount,
		TeachersCount:           item.TeachersCount,
		GroupsCount:             item.GroupsCount,
		LessonsToday:            item.LessonsToday,
		PaymentsThisMonth:       item.PaymentsThisMonth,
		PaymentsAmountThisMonth: item.PaymentsAmountThisMonth,
		StudentDebtTotal:        item.StudentDebtTotal,
		PendingPayrollEntries:   item.PendingPayrollEntries,
		UnreadNotifications:     item.UnreadNotifications,
		RecentAuditLogs:         mapRecentAuditLogs(item.RecentAuditLogs),
	}
}

func mapRecentAuditLogs(items []RecentAuditLog) []RecentAuditLogResponse {
	result := make([]RecentAuditLogResponse, 0, len(items))

	for _, item := range items {
		userID := ""
		if item.UserID != uuid.Nil {
			userID = item.UserID.String()
		}

		entityID := ""
		if item.EntityID != uuid.Nil {
			entityID = item.EntityID.String()
		}

		result = append(result, RecentAuditLogResponse{
			ID:          item.ID.String(),
			UserID:      userID,
			Action:      item.Action,
			EntityType:  item.EntityType,
			EntityID:    entityID,
			Description: item.Description,
			CreatedAt:   item.CreatedAt,
		})
	}

	return result
}
