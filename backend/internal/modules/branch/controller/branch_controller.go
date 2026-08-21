package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"queue-management-tenant/backend/internal/modules/branch/dto"
	"queue-management-tenant/backend/internal/modules/branch/service"
	"queue-management-tenant/backend/pkg/response"
)

type BranchController struct {
	branchService *service.BranchService
}

func NewBranchController(branchService *service.BranchService) *BranchController {
	return &BranchController{branchService: branchService}
}

func (h *BranchController) CreateBranch(c *fiber.Ctx) error {
	orgID := c.Locals("organization_id").(int64)

	var req dto.CreateBranchRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload", err.Error())
	}

	branch, err := h.branchService.CreateBranch(c.Context(), orgID, req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusCreated, "Branch created successfully", branch)
}

func (h *BranchController) ListBranches(c *fiber.Ctx) error {
	orgID := c.Locals("organization_id").(int64)

	branches, err := h.branchService.ListBranches(c.Context(), orgID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Branch list retrieved", branches)
}

func (h *BranchController) GetBranch(c *fiber.Ctx) error {
	orgID := c.Locals("organization_id").(int64)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid branch ID", nil)
	}

	branch, err := h.branchService.GetBranch(c.Context(), orgID, id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Branch not found", nil)
	}

	return response.Success(c, fiber.StatusOK, "Branch details retrieved", branch)
}

func (h *BranchController) GetBranchPublic(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid branch ID", nil)
	}

	branch, err := h.branchService.GetBranchPublic(c.Context(), id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Branch not found", nil)
	}

	return response.Success(c, fiber.StatusOK, "Branch details retrieved", branch)
}

func (h *BranchController) UpdateKioskSettings(c *fiber.Ctx) error {
	orgID := c.Locals("organization_id").(int64)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid branch ID", nil)
	}

	var req dto.UpdateKioskSettingsRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload", err.Error())
	}

	branch, err := h.branchService.UpdateKioskSettings(c.Context(), orgID, id, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Branch kiosk settings updated", branch)
}

