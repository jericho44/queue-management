package routes

import (
	"github.com/gofiber/fiber/v2"
	"queue-management-tenant/backend/internal/middleware"
	"queue-management-tenant/backend/internal/modules/counter/controller"
	"queue-management-tenant/backend/pkg/jwt"
)

func RegisterCounterRoutes(router fiber.Router, counterController *controller.CounterController, jwtSvc *jwt.JWTService) {
	countersGroup := router.Group("/counters", middleware.AuthMiddleware(jwtSvc))

	countersGroup.Post("/", middleware.RequireRoles("SUPER_ADMIN", "OWNER", "MANAGER"), counterController.CreateCounter)
	countersGroup.Get("/", counterController.ListCountersByBranch)
	countersGroup.Post("/:id/open", middleware.RequireRoles("SUPER_ADMIN", "OWNER", "MANAGER", "STAFF"), counterController.OpenCounter)
	countersGroup.Post("/:id/close", middleware.RequireRoles("SUPER_ADMIN", "OWNER", "MANAGER", "STAFF"), counterController.CloseCounter)
}
