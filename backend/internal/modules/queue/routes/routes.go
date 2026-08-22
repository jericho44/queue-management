package routes

import (
	"github.com/gofiber/fiber/v2"
	"queue-management-tenant/backend/internal/middleware"
	"queue-management-tenant/backend/internal/modules/queue/controller"
	"queue-management-tenant/backend/pkg/jwt"
)

func RegisterQueueRoutes(router fiber.Router, queueController *controller.QueueController, jwtSvc *jwt.JWTService) {
	// Public route for public tracking, public display & public kiosk issuance
	router.Get("/public/tickets/:publicToken", queueController.GetPublicTicket)
	router.Post("/public/tickets", queueController.IssuePublicTicket)
	router.Get("/public/display", queueController.ListPublicTicketsByBranch)
	router.Get("/public/branches/:identifier/display", queueController.ListPublicTicketsByBranch)


	// Protected routes
	tickets := router.Group("/tickets", middleware.AuthMiddleware(jwtSvc))

	tickets.Post("/", middleware.RequireRoles("SUPER_ADMIN", "OWNER", "MANAGER", "RECEPTIONIST"), queueController.IssueTicket)
	tickets.Get("/", queueController.ListBranchTickets)
	tickets.Get("/waiting", queueController.ListWaitingTickets)
	tickets.Post("/counters/:counterId/next", middleware.RequireRoles("SUPER_ADMIN", "OWNER", "MANAGER", "STAFF"), queueController.CallNext)

	tickets.Post("/:id/recall", middleware.RequireRoles("SUPER_ADMIN", "OWNER", "MANAGER", "STAFF"), queueController.RecallTicket)
	tickets.Post("/:id/start", middleware.RequireRoles("SUPER_ADMIN", "OWNER", "MANAGER", "STAFF"), queueController.StartServing)
	tickets.Post("/:id/complete", middleware.RequireRoles("SUPER_ADMIN", "OWNER", "MANAGER", "STAFF"), queueController.CompleteTicket)
	tickets.Post("/:id/skip", middleware.RequireRoles("SUPER_ADMIN", "OWNER", "MANAGER", "STAFF"), queueController.SkipTicket)
	tickets.Post("/:id/no-show", middleware.RequireRoles("SUPER_ADMIN", "OWNER", "MANAGER", "STAFF"), queueController.NoShowTicket)
}
