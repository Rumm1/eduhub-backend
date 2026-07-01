package audit

import (
	"context"
	"errors"
	"strings"
	"time"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
)

var (
	ErrTenantRequired  = errors.New("tenant organization is required")
	ErrUserIDInvalid   = errors.New("user id is invalid")
	ErrEntityIDInvalid = errors.New("entity id is invalid")
	ErrFromDateInvalid = errors.New("from date is invalid")
	ErrToDateInvalid   = errors.New("to date is invalid")
	ErrLimitInvalid    = errors.New("limit is invalid")
	ErrOffsetInvalid   = errors.New("offset is invalid")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input CreateAuditLogInput) error {
	return s.repo.Create(ctx, input)
}

func (s *Service) List(
	ctx context.Context,
	userID string,
	action string,
	entityType string,
	entityID string,
	fromDate string,
	toDate string,
	limit int,
	offset int,
) (AuditLogsListResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return AuditLogsListResponse{}, ErrTenantRequired
	}

	userID = strings.TrimSpace(userID)
	if !IsValidUUID(userID) {
		return AuditLogsListResponse{}, ErrUserIDInvalid
	}

	entityID = strings.TrimSpace(entityID)
	if !IsValidUUID(entityID) {
		return AuditLogsListResponse{}, ErrEntityIDInvalid
	}

	fromDate = strings.TrimSpace(fromDate)
	if fromDate != "" {
		if _, err := time.Parse("2006-01-02", fromDate); err != nil {
			return AuditLogsListResponse{}, ErrFromDateInvalid
		}
	}

	toDate = strings.TrimSpace(toDate)
	if toDate != "" {
		if _, err := time.Parse("2006-01-02", toDate); err != nil {
			return AuditLogsListResponse{}, ErrToDateInvalid
		}
	}

	if limit == 0 {
		limit = 50
	}

	if limit < 1 || limit > 200 {
		return AuditLogsListResponse{}, ErrLimitInvalid
	}

	if offset < 0 {
		return AuditLogsListResponse{}, ErrOffsetInvalid
	}

	items, total, err := s.repo.List(ctx, AuditLogFilter{
		OrganizationID: currentUser.OrganizationID.String(),
		UserID:         userID,
		Action:         strings.TrimSpace(action),
		EntityType:     strings.TrimSpace(entityType),
		EntityID:       entityID,
		FromDate:       fromDate,
		ToDate:         toDate,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return AuditLogsListResponse{}, err
	}

	return buildListResponse(items, limit, offset, total), nil
}

func buildListResponse(items []AuditLog, limit int, offset int, total int) AuditLogsListResponse {
	responses := make([]AuditLogResponse, 0, len(items))

	for _, item := range items {
		responses = append(responses, AuditLogResponse{
			ID:             item.ID,
			OrganizationID: item.OrganizationID,
			UserID:         item.UserID,
			UserName:       item.UserName,
			Action:         item.Action,
			EntityType:     item.EntityType,
			EntityID:       item.EntityID,
			Description:    item.Description,
			Metadata:       item.Metadata,
			IPAddress:      item.IPAddress,
			UserAgent:      item.UserAgent,
			CreatedAt:      item.CreatedAt.Format(time.RFC3339),
		})
	}

	return AuditLogsListResponse{
		Items:  responses,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}
}
