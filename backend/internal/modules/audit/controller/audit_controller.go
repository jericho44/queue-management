package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"queue-management-tenant/backend/internal/modules/audit/repository"
	"queue-management-tenant/backend/pkg/response"
)

type AuditController struct {
	auditRepo *repository.AuditRepository
}

func NewAuditController(auditRepo *repository.AuditRepository) *AuditController {
	return &AuditController{auditRepo: auditRepo}
}

func (h *AuditController) List(c *fiber.Ctx) error {
	orgID := c.Locals("organization_id").(int64)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	logs, total, err := h.auditRepo.List(c.Context(), orgID, page, perPage)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Paginated(c, "Audit logs retrieved", logs, page, perPage, total)
}
