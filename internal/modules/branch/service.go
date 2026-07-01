package branch

import (
	"context"
	"errors"
	"strings"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/google/uuid"
)

var (
	ErrBranchNameRequired = errors.New("branch name is required")
	ErrTenantRequired     = errors.New("tenant organization is required")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateBranchRequest) (BranchResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return BranchResponse{}, ErrTenantRequired
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return BranchResponse{}, ErrBranchNameRequired
	}

	newBranch := Branch{
		ID:             uuid.New(),
		OrganizationID: *currentUser.OrganizationID,
		Name:           name,
		Address:        strings.TrimSpace(req.Address),
		Phone:          strings.TrimSpace(req.Phone),
		Status:         "active",
	}

	createdBranch, err := s.repo.Create(ctx, newBranch)
	if err != nil {
		return BranchResponse{}, err
	}

	return buildBranchResponse(createdBranch), nil
}

func (s *Service) List(ctx context.Context) (ListBranchesResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return ListBranchesResponse{}, ErrTenantRequired
	}

	branches, err := s.repo.ListByOrganizationID(ctx, *currentUser.OrganizationID)
	if err != nil {
		return ListBranchesResponse{}, err
	}

	items := make([]BranchResponse, 0, len(branches))
	for _, branch := range branches {
		items = append(items, buildBranchResponse(branch))
	}

	return ListBranchesResponse{
		Items: items,
		Total: len(items),
	}, nil
}

func buildBranchResponse(branch Branch) BranchResponse {
	return BranchResponse{
		ID:             branch.ID.String(),
		OrganizationID: branch.OrganizationID.String(),
		Name:           branch.Name,
		Address:        branch.Address,
		Phone:          branch.Phone,
		Status:         branch.Status,
	}
}
