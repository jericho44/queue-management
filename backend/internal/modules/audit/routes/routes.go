package routes

import (
	"github.com/gofiber/fiber/v2"
	"queue-management-tenant/backend/internal/middleware"
	"queue-management-tenant/backend/internal/modules/audit/controller"
	"queue-management-tenant/backend/pkg/jwt"
)

func RegisterAuditRoutes(router fiber.Router, auditController *controller.AuditController, jwtSvc *jwt.JWTService) {
	audit := router.Group("/audit-logs", middleware.AuthMiddleware(jwtSvc))

	audit.Get("/", middleware.RequireRoles("SUPER_ADMIN", "OWNER"), auditController.List)
}
