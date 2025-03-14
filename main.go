package main

import (
	"auth-service/database"
	"auth-service/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// 📌 Initialiser la base de données
	database.InitDB()

	// 📌 Définir les routes
	routes.SetupRoutes(app)

	// 📌 Démarrer le serveur
	log.Fatal(app.Listen(":8001"))
}
