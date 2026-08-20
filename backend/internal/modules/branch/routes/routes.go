package routes

import (
	"github.com/gofiber/fiber/v2"
	"queue-management-tenant/backend/internal/middleware"
	"queue-management-tenant/backend/internal/modules/branch/controller"
	"queue-management-tenant/backend/pkg/jwt"
)

func RegisterBranchRoutes(router fiber.Router, branchController *controller.BranchController, jwtSvc *jwt.JWTService) {
	branches := router.Group("/branches", middleware.AuthMiddleware(jwtSvc))

	branches.Post("/", middleware.RequireRoles("SUPER_ADMIN", "OWNER", "MANAGER"), branchController.CreateBranch)
	branches.Get("/", branchController.ListBranches)
	branches.Get("/:id", branchController.GetBranch)
}
