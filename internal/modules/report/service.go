package report

import (
	"context"
	"errors"
	"strings"
	"time"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/google/uuid"
)

var (
	ErrTenantRequired    = errors.New("tenant organization is required")
	ErrTeacherIDRequired = errors.New("teacher id is required")
	ErrTeacherIDInvalid  = errors.New("teacher id is invalid")
	ErrTeacherNotFound   = errors.New("teacher not found in organization")
	ErrFromDateRequired  = errors.New("from date is required")
	ErrToDateRequired    = errors.New("to date is required")
	ErrFromDateInvalid   = errors.New("from date is invalid")
	ErrToDateInvalid     = errors.New("to date is invalid")
	ErrDateRangeInvalid  = errors.New("to date must be after or equal from date")
	ErrForbiddenReport   = errors.New("forbidden report access")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetTeacherSchedule(
	ctx context.Context,
	teacherIDRaw string,
	fromDateRaw string,
	toDateRaw string,
) (TeacherScheduleReportResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return TeacherScheduleReportResponse{}, ErrTenantRequired
	}

	fromDate := strings.TrimSpace(fromDateRaw)
	if fromDate == "" {
		return TeacherScheduleReportResponse{}, ErrFromDateRequired
	}

	parsedFromDate, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return TeacherScheduleReportResponse{}, ErrFromDateInvalid
	}

	toDate := strings.TrimSpace(toDateRaw)
	if toDate == "" {
		return TeacherScheduleReportResponse{}, ErrToDateRequired
	}

	parsedToDate, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		return TeacherScheduleReportResponse{}, ErrToDateInvalid
	}

	if parsedToDate.Before(parsedFromDate) {
		return TeacherScheduleReportResponse{}, ErrDateRangeInvalid
	}

	teacherIDText := strings.TrimSpace(teacherIDRaw)

	if teacherIDText == "" {
		if hasRole(currentUser.Roles, "TEACHER") {
			teacherIDText = currentUser.UserID.String()
		} else {
			return TeacherScheduleReportResponse{}, ErrTeacherIDRequired
		}
	}

	teacherID, err := uuid.Parse(teacherIDText)
	if err != nil {
		return TeacherScheduleReportResponse{}, ErrTeacherIDInvalid
	}

	if hasRole(currentUser.Roles, "TEACHER") && teacherID != currentUser.UserID && !hasRole(currentUser.Roles, "ORG_ADMIN") {
		return TeacherScheduleReportResponse{}, ErrForbiddenReport
	}

	report, err := s.repo.GetTeacherSchedule(ctx, *currentUser.OrganizationID, teacherID, fromDate, toDate)
	if err != nil {
		return TeacherScheduleReportResponse{}, err
	}

	return buildTeacherScheduleResponse(report), nil
}

func buildTeacherScheduleResponse(report TeacherScheduleReport) TeacherScheduleReportResponse {
	items := make([]TeacherScheduleReportItem, 0, len(report.Items))

	for _, item := range report.Items {
		items = append(items, TeacherScheduleReportItem{
			LessonID:            item.LessonID,
			LessonDate:          item.LessonDate,
			StartTime:           item.StartTime,
			EndTime:             item.EndTime,
			Hours:               item.Hours,
			Topic:               item.Topic,
			Status:              item.Status,
			GroupID:             item.GroupID,
			GroupName:           item.GroupName,
			BranchID:            item.BranchID,
			BranchName:          item.BranchName,
			SubjectID:           item.SubjectID,
			SubjectName:         item.SubjectName,
			PlannedTeacherID:    item.PlannedTeacherID,
			PlannedTeacherName:  item.PlannedTeacherName,
			ActualTeacherID:     item.ActualTeacherID,
			ActualTeacherName:   item.ActualTeacherName,
			IsSubstitution:      item.IsSubstitution,
			TeacherRoleInLesson: item.TeacherRoleInLesson,
			SubstitutionReason:  item.SubstitutionReason,
		})
	}

	return TeacherScheduleReportResponse{
		TeacherID:          report.TeacherID,
		TeacherName:        report.TeacherName,
		FromDate:           report.FromDate,
		ToDate:             report.ToDate,
		TotalLessons:       report.TotalLessons,
		ActualLessons:      report.ActualLessons,
		PlannedOnlyLessons: report.PlannedOnlyLessons,
		Substitutions:      report.Substitutions,
		TotalActualHours:   report.TotalActualHours,
		Items:              items,
	}
}

func hasRole(roles []string, role string) bool {
	for _, item := range roles {
		if item == role {
			return true
		}
	}

	return false
}
