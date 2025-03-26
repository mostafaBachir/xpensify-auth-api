package database

import (
	"auth-service/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initialise la base avec migration automatique (main.go)
func InitDB() *gorm.DB {
	db := connect()
	DB = db
	AutoMigrateModels(db)
	return db
}

// InitDBWithoutAutoMigrate utilisée dans le seed
func InitDBWithoutAutoMigrate() *gorm.DB {
	db := connect()
	DB = db
	fmt.Println("✅ Connexion à la base établie (sans migration)")
	return db
}

// AutoMigrateModels applique les migrations dans le bon ordre
func AutoMigrateModels(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},    // 🔑 FK dans Permission
		&models.Service{}, // 🔑 FK dans Permission
		&models.Permission{},
	)
	if err != nil {
		log.Fatal("❌ Erreur lors de la migration automatique :", err)
	}
	fmt.Println("✅ Migration exécutée avec succès.")
}

// connect établit la connexion GORM
func connect() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost", "11002", "bachir", "rapido31", "auth_service",
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Erreur de connexion à la base de données :", err)
	}
	return db
}
