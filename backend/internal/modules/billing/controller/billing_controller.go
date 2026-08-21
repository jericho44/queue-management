package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"queue-management-tenant/backend/internal/modules/billing/dto"
	"queue-management-tenant/backend/internal/modules/billing/service"
	"queue-management-tenant/backend/pkg/response"
)

type BillingController struct {
	billingService *service.BillingService
}

func NewBillingController(billingService *service.BillingService) *BillingController {
	return &BillingController{billingService: billingService}
}

func (h *BillingController) GetCurrentUsage(c *fiber.Ctx) error {
	orgID := c.Locals("organization_id").(int64)

	usage, err := h.billingService.GetCurrentUsage(c.Context(), orgID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Current usage meter retrieved", usage)
}

func (h *BillingController) ListTenantInvoices(c *fiber.Ctx) error {
	orgID := c.Locals("organization_id").(int64)

	invoices, err := h.billingService.ListTenantInvoices(c.Context(), orgID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Tenant invoices retrieved", invoices)
}

func (h *BillingController) CreateSnapToken(c *fiber.Ctx) error {
	orgID := c.Locals("organization_id").(int64)
	invoiceID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid invoice ID", nil)
	}

	resp, err := h.billingService.CreateSnapToken(c.Context(), orgID, invoiceID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Payment Snap token generated", resp)
}

func (h *BillingController) HandleMidtransWebhook(c *fiber.Ctx) error {
	var payload dto.MidtransWebhookPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid webhook payload", err.Error())
	}

	if err := h.billingService.ProcessMidtransWebhook(c.Context(), payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Webhook processed successfully", nil)
}

func (h *BillingController) ListAllInvoicesSuperadmin(c *fiber.Ctx) error {
	invoices, err := h.billingService.ListAllInvoices(c.Context())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "All invoices retrieved", invoices)
}

func (h *BillingController) GetSuperadminBillingStats(c *fiber.Ctx) error {
	stats, err := h.billingService.GetSuperadminStats(c.Context())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Superadmin billing stats retrieved", stats)
}
