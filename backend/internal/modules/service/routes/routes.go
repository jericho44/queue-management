package routes

import (
	"github.com/gofiber/fiber/v2"
	"queue-management-tenant/backend/internal/middleware"
	"queue-management-tenant/backend/internal/modules/service/controller"
	"queue-management-tenant/backend/pkg/jwt"
)

func RegisterServiceRoutes(router fiber.Router, serviceController *controller.ServiceController, jwtSvc *jwt.JWTService) {
	servicesGroup := router.Group("/services", middleware.AuthMiddleware(jwtSvc))

	servicesGroup.Post("/", middleware.RequireRoles("SUPER_ADMIN", "OWNER", "MANAGER"), serviceController.CreateService)
	servicesGroup.Get("/", serviceController.ListServicesByBranch)
}
