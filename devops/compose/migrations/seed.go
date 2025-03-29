package main

import (
	"auth-service/database"
	"auth-service/models"
	"auth-service/security"
	"fmt"
	"log"
)

func main() {
	fmt.Println("🧹 Initialisation de la base de données...")

	db := database.InitDBWithoutAutoMigrate()

	// 🔥 Drop & AutoMigrate
	if err := db.Migrator().DropTable(&models.Permission{}, &models.Service{}, &models.User{}); err != nil {
		log.Fatalf("❌ Échec lors du drop des tables : %v", err)
	}
	database.AutoMigrateModels(db)

	// 👑 Superutilisateur
	superuser := models.User{
		Name:         "Bachir Mostafa",
		Email:        "mostafa.bachir@gmail.com",
		Password:     mustHash("rapido31"), // 💡 à passer via env plus tard
		Role:         "superuser",
		RefreshToken: "",
	}
	db.Create(&superuser)

	// 👥 Autres utilisateurs
	users := []models.User{
		{Name: "Alice Admin", Email: "alice@example.com", Password: mustHash("alice123"), Role: "admin"},
		{Name: "Bob Viewer", Email: "bob@example.com", Password: mustHash("bob123"), Role: "user"},
	}
	for _, u := range users {
		db.Create(&u)
	}

	// 📦 Services
	services := []models.Service{
		{Name: "Xpensify Receipt API", Url: "http://localhost:8002"},
		{Name: "Xpensify Dashboard", Url: "http://localhost:3000"},
	}
	for _, s := range services {
		db.Create(&s)
	}

	// 🔁 Récupération dynamique des services
	var receiptService, dashboardService models.Service
	db.Where("name = ?", "Xpensify Receipt API").First(&receiptService)
	db.Where("name = ?", "Xpensify Dashboard").First(&dashboardService)

	// 🔐 Permissions
	permissions := []models.Permission{
		{UserID: 1, ServiceID: fmt.Sprintf("%d", receiptService.ID), Action: "read"},
		{UserID: 1, ServiceID: fmt.Sprintf("%d", receiptService.ID), Action: "write"},
		{UserID: 1, ServiceID: fmt.Sprintf("%d", dashboardService.ID), Action: "manage"},

		{UserID: 2, ServiceID: fmt.Sprintf("%d", receiptService.ID), Action: "read"},
		{UserID: 2, ServiceID: fmt.Sprintf("%d", dashboardService.ID), Action: "read"},

		{UserID: 3, ServiceID: fmt.Sprintf("%d", receiptService.ID), Action: "read"},
	}
	for _, p := range permissions {
		db.Create(&p)
	}

	fmt.Println("✅ Base de données initialisée avec succès.")
}

func mustHash(password string) string {
	hash, err := security.HashPassword(password)
	if err != nil {
		panic("❌ Erreur de hash : " + err.Error())
	}
	return hash
}
