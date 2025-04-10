package sync

import (
	"auth-service/config"
	"auth-service/models"
	"auth-service/services"

	"github.com/gofiber/fiber/v2"
)

// 📌 Handler pour l'inscription d'un utilisateur
func Register(c *fiber.Ctx) error {
	var payload struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// 📌 Vérification du JSON envoyé
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// 📌 Vérification des champs obligatoires
	if payload.Name == "" || payload.Email == "" || payload.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "All fields are required"})
	}

	// 📌 Appel du service pour enregistrer l'utilisateur
	user := &models.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
	}

	if err := services.RegisterUser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User registered successfully"})
}

// 📌 Handler pour la connexion d'un utilisateur
func Login(c *fiber.Ctx) error {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// 📌 Vérification du JSON envoyé
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// 📌 Vérification des champs
	if payload.Email == "" || payload.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Email and password are required"})
	}

	// 📌 Appel du service d'authentification
	accessToken, refreshToken, err := services.LoginUser(payload.Email, payload.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	// 📌 Retour des tokens en JSON
	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"jwt_secret":    config.Get("jwt-access-secret"),
		// ou config.Get("jwt-access-secret")

	})
}

// 📌 Handler pour rafraîchir un token expiré
func RefreshToken(c *fiber.Ctx) error {
	var payload struct {
		RefreshToken string `json:"refresh_token"`
	}

	// 📌 Vérification du JSON envoyé
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// 📌 Vérification du Refresh Token
	if payload.RefreshToken == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Refresh token is required"})
	}

	// 📌 Appel du service pour rafraîchir le token
	accessToken, newRefreshToken, err := services.RefreshUserToken(payload.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	// 📌 Retour des nouveaux tokens en JSON
	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}

// 📌 Handler pour la déconnexion
func Logout(c *fiber.Ctx) error {
	var payload struct {
		UserID uint `json:"user_id"`
	}

	// 📌 Vérification du JSON envoyé
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// 📌 Vérification de l'ID utilisateur
	if payload.UserID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "User ID is required"})
	}

	// 📌 Appel du service pour révoquer le Refresh Token
	if err := services.LogoutUser(payload.UserID); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User logged out successfully"})
}
