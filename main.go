package main

import (
	"auth-service/database"
	"auth-service/pubsub"
	"auth-service/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	// Active CORS avec configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3001, http://127.0.0.1:3001, http://192.168.0.99:3001",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))
	pubsub.InitRedis()

	// 📌 Initialiser la base de données
	database.InitDB()

	// 📌 Définir les routes
	routes.SetupRoutes(app)

	// 📌 Démarrer le serveur
	log.Fatal(app.Listen(":8001"))
}
