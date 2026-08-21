package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"queue-management-tenant/backend/internal/modules/queue/dto"
	"queue-management-tenant/backend/internal/modules/queue/service"
	"queue-management-tenant/backend/pkg/response"
)

type QueueController struct {
	queueService *service.QueueService
}

func NewQueueController(queueService *service.QueueService) *QueueController {
	return &QueueController{queueService: queueService}
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

func getUserID(c *fiber.Ctx) (int64, bool) {
	val := c.Locals("user_id")
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

func (h *QueueController) IssueTicket(c *fiber.Ctx) error {
	orgID, ok := getOrgID(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Organization context missing", nil)
	}

	var req dto.IssueTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload", err.Error())
	}

	ticket, err := h.queueService.IssueTicket(c.Context(), orgID, req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusCreated, "Ticket issued successfully", ticket)
}

func (h *QueueController) IssuePublicTicket(c *fiber.Ctx) error {
	var req dto.IssueTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload", err.Error())
	}

	ticket, err := h.queueService.IssuePublicTicket(c.Context(), req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusCreated, "Ticket issued successfully", ticket)
}

func (h *QueueController) CallNext(c *fiber.Ctx) error {
	orgID, ok := getOrgID(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Organization context missing", nil)
	}
	staffID, _ := getUserID(c)

	counterID, err := strconv.ParseInt(c.Params("counterId"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid counter ID", nil)
	}

	ticket, err := h.queueService.CallNext(c.Context(), orgID, counterID, staffID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Called next ticket", ticket)
}

func (h *QueueController) RecallTicket(c *fiber.Ctx) error {
	orgID, ok := getOrgID(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Organization context missing", nil)
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid ticket ID", nil)
	}

	ticket, err := h.queueService.RecallTicket(c.Context(), orgID, id)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Ticket recalled", ticket)
}

func (h *QueueController) StartServing(c *fiber.Ctx) error {
	orgID, ok := getOrgID(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Organization context missing", nil)
	}
	staffID, _ := getUserID(c)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid ticket ID", nil)
	}

	ticket, err := h.queueService.StartServing(c.Context(), orgID, id, staffID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Serving started", ticket)
}

func (h *QueueController) CompleteTicket(c *fiber.Ctx) error {
	orgID, ok := getOrgID(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Organization context missing", nil)
	}
	staffID, _ := getUserID(c)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid ticket ID", nil)
	}

	ticket, err := h.queueService.CompleteTicket(c.Context(), orgID, id, staffID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Ticket completed", ticket)
}

func (h *QueueController) SkipTicket(c *fiber.Ctx) error {
	orgID, ok := getOrgID(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Organization context missing", nil)
	}
	staffID, _ := getUserID(c)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid ticket ID", nil)
	}

	ticket, err := h.queueService.SkipTicket(c.Context(), orgID, id, staffID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Ticket skipped", ticket)
}

func (h *QueueController) NoShowTicket(c *fiber.Ctx) error {
	orgID, ok := getOrgID(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Organization context missing", nil)
	}
	staffID, _ := getUserID(c)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid ticket ID", nil)
	}

	ticket, err := h.queueService.NoShowTicket(c.Context(), orgID, id, staffID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Ticket marked as no show", ticket)
}

func (h *QueueController) GetPublicTicket(c *fiber.Ctx) error {
	token := c.Params("publicToken")

	ticket, err := h.queueService.GetByPublicToken(c.Context(), token)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Ticket not found", nil)
	}

	return response.Success(c, fiber.StatusOK, "Ticket status retrieved", ticket)
}

func (h *QueueController) ListBranchTickets(c *fiber.Ctx) error {
	orgID, ok := getOrgID(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Organization context missing", nil)
	}

	branchID, err := strconv.ParseInt(c.Query("branch_id"), 10, 64)
	if err != nil || branchID <= 0 {
		return response.Error(c, fiber.StatusBadRequest, "Branch ID required", nil)
	}
	status := c.Query("status")

	tickets, err := h.queueService.ListBranchTickets(c.Context(), orgID, branchID, status)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Queue list retrieved", tickets)
}

func (h *QueueController) ListWaitingTickets(c *fiber.Ctx) error {
	orgID, ok := getOrgID(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Organization context missing", nil)
	}

	branchID, err := strconv.ParseInt(c.Query("branch_id"), 10, 64)
	if err != nil || branchID <= 0 {
		return response.Error(c, fiber.StatusBadRequest, "Branch ID required", nil)
	}

	tickets, err := h.queueService.ListBranchTickets(c.Context(), orgID, branchID, "WAITING")
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Waiting queue list retrieved", tickets)
}
