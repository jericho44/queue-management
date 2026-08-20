package routes

import (
	"github.com/gofiber/fiber/v2"
	"queue-management-tenant/backend/internal/middleware"
	"queue-management-tenant/backend/internal/modules/report/controller"
	"queue-management-tenant/backend/pkg/jwt"
)

func RegisterReportRoutes(router fiber.Router, reportController *controller.ReportController, jwtSvc *jwt.JWTService) {
	reports := router.Group("/reports", middleware.AuthMiddleware(jwtSvc))

	reports.Get("/dashboard", middleware.RequireRoles("SUPER_ADMIN", "OWNER", "MANAGER"), reportController.GetDashboardStats)
}
