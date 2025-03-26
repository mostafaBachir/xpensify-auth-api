package security

import (
	"github.com/gofiber/fiber/v2"
)

func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := c.Locals("role")
		if claims == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		role, ok := claims.(string)
		if !ok || role != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}

		return c.Next()
	}
}
