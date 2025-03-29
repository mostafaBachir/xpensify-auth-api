package main

import (
	"auth-service/config"
	"auth-service/database"
	"auth-service/pubsub"
	"auth-service/routes"
	"auth-service/security"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// 🔐 Middleware CORS
	app.Use(security.CorsMiddleware())

	// 🔌 Init Redis PubSub
	pubsub.InitRedis()

	// 🗄️ Init DB
	database.InitDB()

	// 🚦 Définir les routes
	routes.SetupRoutes(app)

	// 🚀 Démarrer le serveur
	port := config.Get("auth-service-port")
	if port == "" {
		port = "8001"
	}
	log.Fatal(app.Listen(":" + port))
}
