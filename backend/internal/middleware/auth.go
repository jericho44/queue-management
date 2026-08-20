package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"queue-management-tenant/backend/pkg/jwt"
	"queue-management-tenant/backend/pkg/response"
)

func AuthMiddleware(jwtService *jwt.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Missing authorization header", nil)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid authorization header format", nil)
		}

		claims, err := jwtService.ValidateToken(parts[1])
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid or expired token", err.Error())
		}

		c.Locals("user", claims)
		c.Locals("user_id", claims.UserID)
		c.Locals("user_uuid", claims.UserUUID)
		c.Locals("organization_id", claims.OrganizationID)
		c.Locals("org_uuid", claims.OrgUUID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func RequireRoles(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleVal := c.Locals("role")
		if roleVal == nil {
			return response.Error(c, fiber.StatusForbidden, "Access denied: missing role", nil)
		}

		userRole, ok := roleVal.(string)
		if !ok {
			return response.Error(c, fiber.StatusForbidden, "Access denied: invalid role format", nil)
		}

		// Super admin has access to everything
		if userRole == "SUPER_ADMIN" {
			return c.Next()
		}

		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}

		return response.Error(c, fiber.StatusForbidden, "Access denied: insufficient permissions", nil)
	}
}
