package importer

import (
	"context"
	"io"
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

func (s *Service) PreviewStudents(ctx context.Context, filename string, reader io.Reader) (StudentImportPreviewResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return StudentImportPreviewResponse{}, ErrTenantRequired
	}

	rows, err := ParseStudentImportFile(filename, reader)
	if err != nil {
		return StudentImportPreviewResponse{}, err
	}

	_, previewRows, summary, err := s.validateRows(ctx, *currentUser.OrganizationID, rows)
	if err != nil {
		return StudentImportPreviewResponse{}, err
	}

	return StudentImportPreviewResponse{
		Summary: summary,
		Rows:    previewRows,
	}, nil
}

func (s *Service) ConfirmStudents(ctx context.Context, filename string, reader io.Reader) (StudentImportConfirmResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return StudentImportConfirmResponse{}, ErrTenantRequired
	}

	rows, err := ParseStudentImportFile(filename, reader)
	if err != nil {
		return StudentImportConfirmResponse{}, err
	}

	validRows, previewRows, summary, err := s.validateRows(ctx, *currentUser.OrganizationID, rows)
	if err != nil {
		return StudentImportConfirmResponse{}, err
	}

	responseSummary := StudentImportConfirmSummary{
		TotalRows:   summary.TotalRows,
		ValidRows:   summary.ValidRows,
		InvalidRows: summary.InvalidRows,
		WarningRows: summary.WarningRows,
	}

	if summary.InvalidRows > 0 {
		return StudentImportConfirmResponse{
			Summary: responseSummary,
			Rows:    previewRows,
		}, nil
	}

	confirmResult, err := s.repository.ImportStudents(ctx, *currentUser.OrganizationID, validRows)
	if err != nil {
		return StudentImportConfirmResponse{}, err
	}

	responseSummary.CreatedStudents = confirmResult.CreatedStudents
	responseSummary.ReusedStudents = confirmResult.ReusedStudents
	responseSummary.CreatedParents = confirmResult.CreatedParents
	responseSummary.ReusedParents = confirmResult.ReusedParents
	responseSummary.LinkedParentsToStudents = confirmResult.LinkedParentsToStudents
	responseSummary.LinkedStudentsToGroups = confirmResult.LinkedStudentsToGroups

	return StudentImportConfirmResponse{
		Summary: responseSummary,
		Rows:    previewRows,
	}, nil
}

func (s *Service) validateRows(
	ctx context.Context,
	organizationID uuid.UUID,
	rows []StudentImportRow,
) ([]ValidStudentImportRow, []StudentImportPreviewRow, StudentImportSummary, error) {
	validRows := make([]ValidStudentImportRow, 0)
	previewRows := make([]StudentImportPreviewRow, 0, len(rows))
	groupCache := make(map[string]GroupLookup)

	summary := StudentImportSummary{
		TotalRows: len(rows),
	}

	for _, row := range rows {
		errorsList := make([]string, 0)
		warnings := make([]string, 0)

		if strings.TrimSpace(row.StudentFullName) == "" {
			errorsList = append(errorsList, "Student full name is required")
		}

		if strings.TrimSpace(row.ParentFullName) == "" {
			errorsList = append(errorsList, "Parent full name is required")
		}

		if strings.TrimSpace(row.GroupName) == "" {
			errorsList = append(errorsList, "Group name is required")
		}

		if strings.TrimSpace(row.ParentPhone) == "" && strings.TrimSpace(row.ParentEmail) == "" {
			warnings = append(warnings, "Parent phone or email is recommended")
		}

		var group GroupLookup
		if strings.TrimSpace(row.GroupName) != "" {
			cacheKey := strings.ToLower(strings.TrimSpace(row.GroupName))

			cachedGroup, ok := groupCache[cacheKey]
			if ok {
				group = cachedGroup
			} else {
				foundGroup, err := s.repository.FindGroupByName(ctx, organizationID, row.GroupName)
				if err != nil {
					if err == ErrGroupNotFound {
						errorsList = append(errorsList, "Group not found")
					} else {
						return nil, nil, StudentImportSummary{}, err
					}
				} else {
					group = foundGroup
					groupCache[cacheKey] = foundGroup
				}
			}
		}

		status := "valid"
		if len(errorsList) > 0 {
			status = "invalid"
			summary.InvalidRows++
		} else {
			summary.ValidRows++

			validRows = append(validRows, ValidStudentImportRow{
				Row:   row,
				Group: group,
			})
		}

		if len(warnings) > 0 {
			summary.WarningRows++
		}

		previewRows = append(previewRows, StudentImportPreviewRow{
			RowNumber:       row.RowNumber,
			Status:          status,
			Errors:          errorsList,
			Warnings:        warnings,
			StudentFullName: row.StudentFullName,
			StudentPhone:    row.StudentPhone,
			ParentFullName:  row.ParentFullName,
			ParentPhone:     row.ParentPhone,
			ParentEmail:     row.ParentEmail,
			GroupName:       row.GroupName,
			BranchName:      group.Branch,
			Relation:        row.Relation,
		})
	}

	return validRows, previewRows, summary, nil
}
