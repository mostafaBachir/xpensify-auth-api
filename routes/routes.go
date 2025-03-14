package routes

import (
	"auth-service/features/sync"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/auth")
	api.Post("/register", sync.Register)    // 🔹 Inscription
	api.Post("/login", sync.Login)          // 🔹 Connexion
	api.Post("/refresh", sync.RefreshToken) // 🔹 Rafraîchissement du token
	api.Post("/logout", sync.Logout)        // 🔹 Déconnexion
}
