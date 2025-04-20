package security

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CorsMiddleware retourne une config CORS basée sur le .env
func CorsMiddleware() fiber.Handler {

	return cors.New(cors.Config{
		AllowOrigins:     "https://*.xpensify.ca,http://127.0.0.1:3001",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	})
}
