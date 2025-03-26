package sync

import (
	"auth-service/database"
	"auth-service/models"
	"auth-service/services"

	"github.com/gofiber/fiber/v2"
)

// 🔹 Liste tous les utilisateurs avec leurs permissions et services
func ListUsers(c *fiber.Ctx) error {
	var users []models.User
	if err := database.DB.Preload("Permissions.Service").Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erreur lors de la récupération des utilisateurs"})
	}

	var response []fiber.Map

	for _, user := range users {
		permissions := make([]fiber.Map, 0)

		for _, p := range user.Permissions {
			permissions = append(permissions, fiber.Map{
				"service_id": p.ServiceID,
				"service":    p.Service.Name,
				"permission": p.Action,
			})
		}

		response = append(response, fiber.Map{
			"id":          user.ID,
			"name":        user.Name,
			"email":       user.Email,
			"role":        user.Role,
			"permissions": permissions,
		})
	}

	return c.JSON(response)
}

// 🔹 Met à jour les permissions d'un utilisateur
func UpdateUserPermissions(c *fiber.Ctx) error {
	var payload []models.Permission
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID utilisateur invalide"})
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Corps de requête invalide"})
	}

	if err := services.UpdateUserPermissions(uint(userID), payload); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erreur lors de la mise à jour des permissions"})
	}

	return c.JSON(fiber.Map{"message": "Permissions mises à jour avec succès"})
}
