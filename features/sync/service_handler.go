package sync

import (
	"auth-service/models"
	"auth-service/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GET /services
func ListServices(c *fiber.Ctx) error {
	servicesList, err := services.GetAllServices()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch services"})
	}
	return c.JSON(servicesList)
}

// POST /services
func CreateService(c *fiber.Ctx) error {
	var s models.Service
	if err := c.BodyParser(&s); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	if err := services.CreateService(&s); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(s)
}

// DELETE /services/:id
func DeleteService(c *fiber.Ctx) error {

	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := services.DeleteService(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(204)
}
func UpdateService(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	var input models.Service

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	if err := services.UpdateService(uint(id), &input); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Service updated successfully"})
}
