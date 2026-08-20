package controller

import (
	"github.com/gofiber/fiber/v2"
	"queue-management-tenant/backend/internal/modules/auth/dto"
	"queue-management-tenant/backend/internal/modules/auth/service"
	"queue-management-tenant/backend/pkg/response"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (h *AuthController) RegisterOrganization(c *fiber.Ctx) error {
	var req dto.RegisterOrganizationRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload", err.Error())
	}

	res, err := h.authService.RegisterOrganization(c.Context(), req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusCreated, "Organization registered successfully", res)
}

func (h *AuthController) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload", err.Error())
	}

	res, err := h.authService.Login(c.Context(), req)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusOK, "Login successful", res)
}

func (h *AuthController) Me(c *fiber.Ctx) error {
	user := c.Locals("user")
	return response.Success(c, fiber.StatusOK, "Current user profile", user)
}

func (h *AuthController) CreateUser(c *fiber.Ctx) error {
	orgID := c.Locals("organization_id").(int64)

	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload", err.Error())
	}

	user, err := h.authService.CreateUser(c.Context(), orgID, req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return response.Success(c, fiber.StatusCreated, "User created successfully", user)
}
