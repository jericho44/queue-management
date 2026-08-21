package routes

import (
	"github.com/gofiber/fiber/v2"
	"queue-management-tenant/backend/internal/middleware"
	"queue-management-tenant/backend/internal/modules/billing/controller"
	"queue-management-tenant/backend/pkg/jwt"
)

func RegisterBillingRoutes(router fiber.Router, billingCtrl *controller.BillingController, jwtSvc *jwt.JWTService) {
	// Public Midtrans Webhook Callback
	router.Post("/billing/webhooks", billingCtrl.HandleMidtransWebhook)

	// Protected Tenant Billing Routes
	billing := router.Group("/billing", middleware.AuthMiddleware(jwtSvc))
	billing.Get("/usage", billingCtrl.GetCurrentUsage)
	billing.Get("/invoices", billingCtrl.ListTenantInvoices)
	billing.Post("/invoices/:id/pay", middleware.RequireRoles("SUPER_ADMIN", "OWNER"), billingCtrl.CreateSnapToken)

	// Protected Superadmin Billing Routes
	superadminBilling := router.Group("/superadmin/billing", middleware.AuthMiddleware(jwtSvc), middleware.RequireRoles("SUPER_ADMIN"))
	superadminBilling.Get("/stats", billingCtrl.GetSuperadminBillingStats)
	superadminBilling.Get("/invoices", billingCtrl.ListAllInvoicesSuperadmin)
}
