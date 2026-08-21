package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"queue-management-tenant/backend/internal/modules/service/dto"
	"queue-management-tenant/backend/internal/modules/service/service"
	"queue-management-tenant/backend/pkg/response"
)

type ServiceController struct {
	svcManagement *service.ServiceManagementService
}

func NewServiceController(svcManagement *service.ServiceManagementService) *ServiceController {
	return &ServiceController{svcManagement: svcManagement}
}

func getOrgID(c *fiber.Ctx) (int64, bool) {
	val := c.Locals("organization_id")
	if val == nil {
		return 0, false
	}
	if id, ok := val.(int64); ok {
		return id, true
	}
	if f, ok := val.(float64); ok {
		return int64(f), true
	}
	if i, ok := val.(int); ok {
		return int64(i), true
	}
	return 0, false
}

func (h *ServiceController) CreateService(c *fiber.Ctx) error {
	orgID, ok := getOrgID(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Organization context missing", nil)
	}

	var req dto.CreateServiceRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload", err.Error())
	}

	svc, err := h.svcManagement.CreateService(c.Context(), orgID, req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusCreated, "Service created successfully", svc)
}

func (h *ServiceController) ListServicesByBranch(c *fiber.Ctx) error {
	orgID, ok := getOrgID(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Organization context missing", nil)
	}

	branchID, err := strconv.ParseInt(c.Query("branch_id"), 10, 64)
	if err != nil || branchID <= 0 {
		return response.Error(c, fiber.StatusBadRequest, "Branch ID required", nil)
	}

	services, err := h.svcManagement.ListServicesByBranch(c.Context(), orgID, branchID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Service list retrieved", services)
}

func (h *ServiceController) ListServicesByBranchPublic(c *fiber.Ctx) error {
	branchParam := c.Query("branch_id")
	if branchParam == "" {
		return response.Error(c, fiber.StatusBadRequest, "Branch parameter required", nil)
	}

	services, err := h.svcManagement.ListServicesByBranchPublic(c.Context(), branchParam)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Service list retrieved", services)
}
