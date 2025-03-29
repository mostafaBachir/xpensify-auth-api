package security

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CorsMiddleware retourne une config CORS basée sur le .env
func CorsMiddleware() fiber.Handler {
	allowOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowOrigins == "" {
		allowOrigins = "*" // fallback si non défini
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	})
}
