package database

import (
	"auth-service/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
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
	// Charger le fichier .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Aucun fichier .env trouvé, on utilise les variables d'environnement système.")
	}

	// Lire les variables d'env
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Erreur de connexion à la base de données :", err)
	}
	return db
}
