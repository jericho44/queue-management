package routes

import (
	"github.com/gofiber/fiber/v2"
	"queue-management-tenant/backend/internal/middleware"
	"queue-management-tenant/backend/internal/modules/auth/controller"
	"queue-management-tenant/backend/pkg/jwt"
)

func RegisterAuthRoutes(router fiber.Router, authController *controller.AuthController, jwtSvc *jwt.JWTService) {
	auth := router.Group("/auth")

	// Public auth routes
	auth.Post("/register-org", authController.RegisterOrganization)
	auth.Post("/login", authController.Login)

	// Protected auth routes
	authMiddleware := middleware.AuthMiddleware(jwtSvc)
	auth.Get("/me", authMiddleware, authController.Me)
	auth.Post("/users", authMiddleware, middleware.RequireRoles("SUPER_ADMIN", "OWNER", "MANAGER"), authController.CreateUser)
}
