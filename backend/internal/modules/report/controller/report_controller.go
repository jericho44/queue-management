package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"queue-management-tenant/backend/internal/modules/report/service"
	"queue-management-tenant/backend/pkg/response"
)

type ReportController struct {
	reportService *service.ReportService
}

func NewReportController(reportService *service.ReportService) *ReportController {
	return &ReportController{reportService: reportService}
}

func (h *ReportController) GetDashboardStats(c *fiber.Ctx) error {
	orgID := c.Locals("organization_id").(int64)
	branchID, _ := strconv.ParseInt(c.Query("branch_id"), 10, 64)

	stats, err := h.reportService.GetDashboardStats(c.Context(), orgID, branchID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Dashboard statistics retrieved", stats)
}
