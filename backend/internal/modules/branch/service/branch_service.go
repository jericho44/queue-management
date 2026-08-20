package service

import (
	"context"
	"database/sql"
	"fmt"

	"queue-management-tenant/backend/internal/modules/branch/dto"
	"queue-management-tenant/backend/internal/modules/branch/entity"
	branchRepo "queue-management-tenant/backend/internal/modules/branch/repository"
	orgRepo "queue-management-tenant/backend/internal/modules/organization/repository"
)

type BranchService struct {
	branchRepo *branchRepo.BranchRepository
	orgRepo    *orgRepo.OrganizationRepository
}

func NewBranchService(branchRepo *branchRepo.BranchRepository, orgRepo *orgRepo.OrganizationRepository) *BranchService {
	return &BranchService{
		branchRepo: branchRepo,
		orgRepo:    orgRepo,
	}
}

func (s *BranchService) CreateBranch(ctx context.Context, orgID int64, req dto.CreateBranchRequest) (*entity.Branch, error) {
	sub, err := s.orgRepo.GetActiveSubscription(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subscription: %w", err)
	}

	count, err := s.branchRepo.CountByOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if count >= sub.MaxBranches {
		return nil, fmt.Errorf("branch limit reached for your plan (%d max)", sub.MaxBranches)
	}

	branch := &entity.Branch{
		OrganizationID: orgID,
		Name:           req.Name,
		Code:           req.Code,
		Address:        sql.NullString{String: req.Address, Valid: req.Address != ""},
		Phone:          sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		Status:         "ACTIVE",
	}

	if err := s.branchRepo.Create(ctx, branch); err != nil {
		return nil, err
	}
	return branch, nil
}

func (s *BranchService) ListBranches(ctx context.Context, orgID int64) ([]entity.Branch, error) {
	return s.branchRepo.ListByOrg(ctx, orgID)
}

func (s *BranchService) GetBranch(ctx context.Context, orgID, id int64) (*entity.Branch, error) {
	return s.branchRepo.GetByID(ctx, orgID, id)
}
