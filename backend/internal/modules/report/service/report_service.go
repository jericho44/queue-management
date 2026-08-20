package service

import (
	"context"

	"queue-management-tenant/backend/internal/modules/report/dto"
	reportRepo "queue-management-tenant/backend/internal/modules/report/repository"
)

type ReportService struct {
	reportRepo *reportRepo.ReportRepository
}

func NewReportService(reportRepo *reportRepo.ReportRepository) *ReportService {
	return &ReportService{reportRepo: reportRepo}
}

func (s *ReportService) GetDashboardStats(ctx context.Context, orgID, branchID int64) (*dto.DashboardStatsResponse, error) {
	return s.reportRepo.GetDashboardStats(ctx, orgID, branchID)
}
