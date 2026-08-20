package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"queue-management-tenant/backend/internal/modules/counter/dto"
	"queue-management-tenant/backend/internal/modules/counter/service"
	"queue-management-tenant/backend/pkg/response"
)

type CounterController struct {
	counterService *service.CounterService
}

func NewCounterController(counterService *service.CounterService) *CounterController {
	return &CounterController{counterService: counterService}
}

func (h *CounterController) CreateCounter(c *fiber.Ctx) error {
	orgID := c.Locals("organization_id").(int64)

	var req dto.CreateCounterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload", err.Error())
	}

	counter, err := h.counterService.CreateCounter(c.Context(), orgID, req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusCreated, "Counter created successfully", counter)
}

func (h *CounterController) ListCountersByBranch(c *fiber.Ctx) error {
	orgID := c.Locals("organization_id").(int64)
	branchID, err := strconv.ParseInt(c.Query("branch_id"), 10, 64)
	if err != nil || branchID <= 0 {
		return response.Error(c, fiber.StatusBadRequest, "Branch ID required", nil)
	}

	counters, err := h.counterService.ListCountersByBranch(c.Context(), orgID, branchID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Counter list retrieved", counters)
}

func (h *CounterController) OpenCounter(c *fiber.Ctx) error {
	orgID := c.Locals("organization_id").(int64)
	staffID := c.Locals("user_id").(int64)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid counter ID", nil)
	}

	if err := h.counterService.OpenCounter(c.Context(), orgID, id, staffID); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Counter opened successfully", nil)
}

func (h *CounterController) CloseCounter(c *fiber.Ctx) error {
	orgID := c.Locals("organization_id").(int64)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid counter ID", nil)
	}

	if err := h.counterService.CloseCounter(c.Context(), orgID, id); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Counter closed successfully", nil)
}
