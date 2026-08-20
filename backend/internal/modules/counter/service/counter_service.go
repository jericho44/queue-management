package service

import (
	"context"
	"database/sql"
	"fmt"

	"queue-management-tenant/backend/internal/modules/counter/dto"
	counterEntity "queue-management-tenant/backend/internal/modules/counter/entity"
	counterRepo "queue-management-tenant/backend/internal/modules/counter/repository"
	orgRepo "queue-management-tenant/backend/internal/modules/organization/repository"
)

type CounterService struct {
	counterRepo *counterRepo.CounterRepository
	orgRepo     *orgRepo.OrganizationRepository
}

func NewCounterService(counterRepo *counterRepo.CounterRepository, orgRepo *orgRepo.OrganizationRepository) *CounterService {
	return &CounterService{
		counterRepo: counterRepo,
		orgRepo:     orgRepo,
	}
}

func (s *CounterService) CreateCounter(ctx context.Context, orgID int64, req dto.CreateCounterRequest) (*counterEntity.Counter, error) {
	sub, err := s.orgRepo.GetActiveSubscription(ctx, orgID)
	if err != nil {
		return nil, err
	}

	count, err := s.counterRepo.CountByOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if count >= sub.MaxCounters {
		return nil, fmt.Errorf("counter limit reached for your plan (%d max)", sub.MaxCounters)
	}

	counter := &counterEntity.Counter{
		OrganizationID: orgID,
		BranchID:       req.BranchID,
		CounterNumber:  req.CounterNumber,
		Name:           req.Name,
		Status:         "CLOSED",
	}

	if err := s.counterRepo.Create(ctx, counter, req.ServiceIDs); err != nil {
		return nil, err
	}
	return counter, nil
}

func (s *CounterService) ListCountersByBranch(ctx context.Context, orgID, branchID int64) ([]counterEntity.Counter, error) {
	return s.counterRepo.ListByBranch(ctx, orgID, branchID)
}

func (s *CounterService) OpenCounter(ctx context.Context, orgID, counterID, staffID int64) error {
	_, err := s.counterRepo.GetByID(ctx, orgID, counterID)
	if err != nil {
		return err
	}

	return s.counterRepo.UpdateStatusAndStaff(ctx, counterID, "OPEN", sql.NullInt64{Int64: staffID, Valid: true})
}

func (s *CounterService) CloseCounter(ctx context.Context, orgID, counterID int64) error {
	_, err := s.counterRepo.GetByID(ctx, orgID, counterID)
	if err != nil {
		return err
	}

	return s.counterRepo.UpdateStatusAndStaff(ctx, counterID, "CLOSED", sql.NullInt64{Valid: false})
}
