package database

import (
	"auth-service/config"
	"auth-service/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	db := connect()
	DB = db
	AutoMigrateModels(db)
	return db
}

func InitDBWithoutAutoMigrate() *gorm.DB {
	db := connect()
	DB = db
	fmt.Println("✅ Connexion à la base établie (sans migration)")
	return db
}

func AutoMigrateModels(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Service{},
		&models.Permission{},
	)
	if err != nil {
		log.Fatal("❌ Erreur lors de la migration automatique :", err)
	}
	fmt.Println("✅ Migration exécutée avec succès.")
}

func connect() *gorm.DB {
	// Lecture via Key Vault-compatible noms (avec tirets)
	host := config.Get("pg-db-host")
	port := config.Get("pg-db-port")
	user := config.Get("pg-db-user")
	password := config.Get("pg-db-password")
	dbname := config.Get("pg-db-name")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Erreur de connexion à la base de données :", err)
	}
	return db
}
