package notification

import (
	"context"
	"strings"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/google/uuid"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context) (ListNotificationsResponse, error) {
	currentUser, err := getCurrentNotificationUser(ctx)
	if err != nil {
		return ListNotificationsResponse{}, err
	}

	items, err := s.repository.ListForUser(ctx, *currentUser.OrganizationID, currentUser.UserID)
	if err != nil {
		return ListNotificationsResponse{}, err
	}

	responseItems := mapNotifications(items)

	unread := 0
	for _, item := range items {
		if !item.IsRead {
			unread++
		}
	}

	return ListNotificationsResponse{
		Items:  responseItems,
		Total:  len(responseItems),
		Unread: unread,
	}, nil
}

func (s *Service) ListTypes() ListNotificationTypesResponse {
	items := []NotificationTypeResponse{
		{
			Code:        NotificationTypeNormal,
			Title:       "Normal",
			Description: "Default notification without high priority",
		},
		{
			Code:        NotificationTypeImportant,
			Title:       "Important",
			Description: "Important notification that should be highlighted",
		},
		{
			Code:        NotificationTypeSystem,
			Title:       "System",
			Description: "System notification from the platform",
		},
		{
			Code:        NotificationTypeWarning,
			Title:       "Warning",
			Description: "Warning or potential problem",
		},
		{
			Code:        NotificationTypePayment,
			Title:       "Payment",
			Description: "Notification related to payments or student balances",
		},
		{
			Code:        NotificationTypeSchedule,
			Title:       "Schedule",
			Description: "Notification related to schedule changes",
		},
		{
			Code:        NotificationTypeHomework,
			Title:       "Homework",
			Description: "Notification related to homework",
		},
		{
			Code:        NotificationTypePayroll,
			Title:       "Payroll",
			Description: "Notification related to payroll workflow",
		},
		{
			Code:        NotificationTypeLesson,
			Title:       "Lesson",
			Description: "Notification related to lessons",
		},
		{
			Code:        NotificationTypeMessage,
			Title:       "Message",
			Description: "Regular user or organization message",
		},
	}

	return ListNotificationTypesResponse{
		Items: items,
		Total: len(items),
	}
}

func (s *Service) Create(ctx context.Context, request CreateNotificationRequest) (CreateNotificationsResponse, error) {
	currentUser, err := getCurrentNotificationUser(ctx)
	if err != nil {
		return CreateNotificationsResponse{}, err
	}

	title := strings.TrimSpace(request.Title)
	if title == "" {
		return CreateNotificationsResponse{}, ErrTitleRequired
	}

	notificationType, err := normalizeNotificationType(request.Type)
	if err != nil {
		return CreateNotificationsResponse{}, err
	}

	targetUserIDs, err := s.resolveTargetUserIDs(ctx, *currentUser.OrganizationID, request)
	if err != nil {
		return CreateNotificationsResponse{}, err
	}

	items, err := s.repository.CreateMany(
		ctx,
		*currentUser.OrganizationID,
		targetUserIDs,
		title,
		strings.TrimSpace(request.Message),
		notificationType,
	)
	if err != nil {
		return CreateNotificationsResponse{}, err
	}

	responseItems := mapNotifications(items)

	return CreateNotificationsResponse{
		Items: responseItems,
		Total: len(responseItems),
	}, nil
}

func (s *Service) MarkRead(ctx context.Context, notificationIDRaw string) error {
	currentUser, err := getCurrentNotificationUser(ctx)
	if err != nil {
		return err
	}

	notificationID, err := uuid.Parse(notificationIDRaw)
	if err != nil {
		return ErrNotificationIDInvalid
	}

	return s.repository.MarkRead(ctx, *currentUser.OrganizationID, currentUser.UserID, notificationID)
}

func (s *Service) MarkAllRead(ctx context.Context) (MarkAllReadResponse, error) {
	currentUser, err := getCurrentNotificationUser(ctx)
	if err != nil {
		return MarkAllReadResponse{}, err
	}

	updated, err := s.repository.MarkAllRead(ctx, *currentUser.OrganizationID, currentUser.UserID)
	if err != nil {
		return MarkAllReadResponse{}, err
	}

	return MarkAllReadResponse{Updated: updated}, nil
}

func (s *Service) Delete(ctx context.Context, notificationIDRaw string) error {
	currentUser, err := getCurrentNotificationUser(ctx)
	if err != nil {
		return err
	}

	notificationID, err := uuid.Parse(notificationIDRaw)
	if err != nil {
		return ErrNotificationIDInvalid
	}

	return s.repository.Delete(ctx, *currentUser.OrganizationID, currentUser.UserID, notificationID)
}

func (s *Service) resolveTargetUserIDs(
	ctx context.Context,
	organizationID uuid.UUID,
	request CreateNotificationRequest,
) ([]uuid.UUID, error) {
	rawUserIDs := make([]string, 0)

	if strings.TrimSpace(request.UserID) != "" {
		rawUserIDs = append(rawUserIDs, strings.TrimSpace(request.UserID))
	}

	for _, rawUserID := range request.UserIDs {
		rawUserID = strings.TrimSpace(rawUserID)
		if rawUserID == "" {
			continue
		}

		rawUserIDs = append(rawUserIDs, rawUserID)
	}

	if len(rawUserIDs) == 0 {
		return s.repository.ListOrganizationUserIDs(ctx, organizationID)
	}

	seen := make(map[uuid.UUID]bool)
	targetUserIDs := make([]uuid.UUID, 0, len(rawUserIDs))

	for _, rawUserID := range rawUserIDs {
		userID, err := uuid.Parse(rawUserID)
		if err != nil {
			return nil, ErrTargetUserInvalid
		}

		if seen[userID] {
			continue
		}

		exists, err := s.repository.UserExistsInOrganization(ctx, organizationID, userID)
		if err != nil {
			return nil, err
		}

		if !exists {
			return nil, ErrTargetUserNotFound
		}

		seen[userID] = true
		targetUserIDs = append(targetUserIDs, userID)
	}

	return targetUserIDs, nil
}

func getCurrentNotificationUser(ctx context.Context) (usercontext.UserContext, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return usercontext.UserContext{}, ErrTenantRequired
	}

	return currentUser, nil
}

func normalizeNotificationType(rawType string) (string, error) {
	notificationType := strings.ToLower(strings.TrimSpace(rawType))
	if notificationType == "" {
		return NotificationTypeNormal, nil
	}

	allowedTypes := map[string]bool{
		NotificationTypeNormal:    true,
		NotificationTypeImportant: true,
		NotificationTypeSystem:    true,
		NotificationTypeWarning:   true,
		NotificationTypePayment:   true,
		NotificationTypeSchedule:  true,
		NotificationTypeHomework:  true,
		NotificationTypePayroll:   true,
		NotificationTypeLesson:    true,
		NotificationTypeMessage:   true,
	}

	if !allowedTypes[notificationType] {
		return "", ErrTypeInvalid
	}

	return notificationType, nil
}

func mapNotifications(items []Notification) []NotificationResponse {
	result := make([]NotificationResponse, 0, len(items))

	for _, item := range items {
		organizationID := ""
		if item.OrganizationID != uuid.Nil {
			organizationID = item.OrganizationID.String()
		}

		result = append(result, NotificationResponse{
			ID:             item.ID.String(),
			OrganizationID: organizationID,
			UserID:         item.UserID.String(),
			Title:          item.Title,
			Message:        item.Message,
			Type:           item.Type,
			IsRead:         item.IsRead,
			CreatedAt:      item.CreatedAt,
		})
	}

	return result
}
