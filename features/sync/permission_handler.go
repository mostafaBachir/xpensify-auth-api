package sync

import (
	"auth-service/models"
	"auth-service/services"

	"github.com/gofiber/fiber/v2"
)

// POST /permissions
func AssignPermission(c *fiber.Ctx) error {
	var p models.Permission
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := services.AssignPermission(&p); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Permission assigned successfully",
	})
}

// DELETE /permissions
func RemovePermission(c *fiber.Ctx) error {
	var body struct {
		UserID    uint   `json:"user_id"`
		ServiceID string `json:"service_id"`
		Action    string `json:"action"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := services.DeletePermission(body.UserID, body.ServiceID, body.Action); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
