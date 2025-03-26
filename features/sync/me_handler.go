// features/sync/me_handler.go
package sync

import (
	"auth-service/database"
	"auth-service/models"

	"github.com/gofiber/fiber/v2"
)

func GetMe(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok || userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var user models.User
	err := database.DB.
		Preload("Permissions").
		Preload("Permissions.Service").
		First(&user, userID).Error

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	var perms []models.PermissionEntry
	for _, p := range user.Permissions {
		perms = append(perms, models.PermissionEntry{
			ServiceID:  p.ServiceID,
			Service:    p.Service.Name,
			Permission: p.Action,
		})
	}

	response := models.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Role:        user.Role,
		Permissions: perms,
	}

	return c.JSON(response)
}
