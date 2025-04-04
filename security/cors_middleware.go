package security

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CorsMiddleware retourne une config CORS basée sur le .env
func CorsMiddleware() fiber.Handler {
	allowOrigins := "*"
	if allowOrigins == "" {
		allowOrigins = "*" // fallback si non défini
	}

	return cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3001,https://www.xpensify.ca",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	})
}
