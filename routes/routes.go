package routes

import (
	"auth-service/features/sync"
	"auth-service/security"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/auth")

	// 🔹 Authentification
	api.Post("/register", sync.Register)
	api.Post("/login", sync.Login)
	api.Post("/refresh", sync.RefreshToken)
	api.Post("/logout", sync.Logout)
	api.Get("/me", security.JWTMiddleware, sync.GetMe)

	// 🔹 Gestion des services
	api.Get("/services", security.JWTMiddleware, security.RequireRole("superuser"), sync.ListServices)
	api.Post("/services", security.JWTMiddleware, security.RequireRole("superuser"), sync.CreateService)
	api.Put("/services/:id", security.JWTMiddleware, security.RequireRole("superuser"), sync.UpdateService)
	api.Delete("/services/:id", security.JWTMiddleware, security.RequireRole("superuser"), sync.DeleteService)

	// 🔹 Gestion des permissions
	api.Post("/permissions", security.JWTMiddleware, security.RequireRole("superuser"), sync.AssignPermission)
	api.Delete("/permissions/:id", security.JWTMiddleware, security.RequireRole("superuser"), sync.RemovePermission)

	// 🔹 Gestion des utilisateurs et de leurs permissions
	api.Get("/users", security.JWTMiddleware, security.RequireRole("superuser"), sync.ListUsers)
	api.Put("/users/:id/permissions", security.JWTMiddleware, security.RequireRole("superuser"), sync.UpdateUserPermissions)
}
