package service

import (
	"context"

	"queue-management-tenant/backend/internal/modules/service/dto"
	"queue-management-tenant/backend/internal/modules/service/entity"
	svcRepo "queue-management-tenant/backend/internal/modules/service/repository"
)

type ServiceManagementService struct {
	serviceRepo *svcRepo.ServiceRepository
}

func NewServiceManagementService(serviceRepo *svcRepo.ServiceRepository) *ServiceManagementService {
	return &ServiceManagementService{serviceRepo: serviceRepo}
}

func (s *ServiceManagementService) CreateService(ctx context.Context, orgID int64, req dto.CreateServiceRequest) (*entity.Service, error) {
	avgTime := req.AvgServiceTimeSec
	if avgTime <= 0 {
		avgTime = 480
	}
	priorityWeight := req.PriorityWeight
	if priorityWeight <= 0 {
		priorityWeight = 1
	}

	svc := &entity.Service{
		OrganizationID:   orgID,
		BranchID:         req.BranchID,
		Name:             req.Name,
		Code:             req.Code,
		Prefix:           req.Prefix,
		AvgServiceTimeSec: avgTime,
		PriorityWeight:   priorityWeight,
		Status:           "ACTIVE",
	}

	if err := s.serviceRepo.Create(ctx, svc); err != nil {
		return nil, err
	}
	return svc, nil
}

func (s *ServiceManagementService) ListServicesByBranch(ctx context.Context, orgID, branchID int64) ([]entity.Service, error) {
	return s.serviceRepo.ListByBranch(ctx, orgID, branchID)
}

func (s *ServiceManagementService) ListServicesByBranchPublic(ctx context.Context, branchIdentifier string) ([]entity.Service, error) {
	return s.serviceRepo.ListByBranchPublic(ctx, branchIdentifier)
}


func (s *ServiceManagementService) GetService(ctx context.Context, orgID, id int64) (*entity.Service, error) {
	return s.serviceRepo.GetByID(ctx, orgID, id)
}

