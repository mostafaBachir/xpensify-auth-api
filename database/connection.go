package database

import (
	"auth-service/models"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // Driver pur pour sql.Open
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
	createDatabaseIfNotExists()

	host := os.Getenv("PG_DB_HOST")
	port := os.Getenv("PG_DB_PORT")
	user := os.Getenv("PG_DB_USER")
	password := os.Getenv("PG_DB_PASSWORD")
	dbname := os.Getenv("PG_DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Erreur de connexion à la base de données :", err)
	}
	fmt.Println("✅ Connexion à la base réussie :", dbname)
	return db
}

func createDatabaseIfNotExists() {
	host := os.Getenv("PG_DB_HOST")
	port := os.Getenv("PG_DB_PORT")
	user := os.Getenv("PG_DB_USER")
	password := os.Getenv("PG_DB_PASSWORD")
	dbname := os.Getenv("PG_DB_NAME")

	postgresDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		host, port, user, password,
	)

	db, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		log.Fatal("❌ Erreur connexion Postgres (postgres db) :", err)
	}
	defer db.Close()

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	err = db.QueryRow(query, dbname).Scan(&exists)
	if err != nil {
		log.Fatal("❌ Erreur vérification existence DB :", err)
	}

	if !exists {
		_, err = db.Exec("CREATE DATABASE " + dbname)
		if err != nil {
			log.Fatal("❌ Erreur création base :", err)
		}
		fmt.Println("✅ Base de données créée :", dbname)
	} else {
		fmt.Println("ℹ️ Base de données déjà existante :", dbname)
	}
}
