package service

import (
	"context"
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
		Address:        req.Address,
		Phone:          req.Phone,

		Status:         "ACTIVE",
		KioskEnabled:   true,
		KioskMode:      "DUAL",
		PaperSize:      "58mm",
		ReceiptHeader:  "",
		ReceiptFooter:  "Terima kasih atas kunjungan Anda. Harap menunggu hingga nomor Anda dipanggil.",
		AutoPrint:      false,
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

func (s *BranchService) GetBranchPublic(ctx context.Context, identifier string) (*entity.Branch, error) {
	return s.branchRepo.GetByIDPublic(ctx, identifier)
}


func (s *BranchService) UpdateKioskSettings(ctx context.Context, orgID, branchID int64, req dto.UpdateKioskSettingsRequest) (*entity.Branch, error) {
	branch, err := s.branchRepo.GetByID(ctx, orgID, branchID)
	if err != nil {
		return nil, err
	}

	if req.KioskEnabled != nil {
		branch.KioskEnabled = *req.KioskEnabled
	}
	if req.KioskMode != nil {
		branch.KioskMode = *req.KioskMode
	}
	if req.PaperSize != nil {
		branch.PaperSize = *req.PaperSize
	}
	if req.ReceiptHeader != nil {
		branch.ReceiptHeader = *req.ReceiptHeader
	}
	if req.ReceiptFooter != nil {
		branch.ReceiptFooter = *req.ReceiptFooter
	}
	if req.AutoPrint != nil {
		branch.AutoPrint = *req.AutoPrint
	}

	if err := s.branchRepo.UpdateKioskSettings(ctx, branch); err != nil {
		return nil, err
	}

	return branch, nil
}

