package config

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var once sync.Once

// loadEnv charge .env une seule fois (utile pour le dev local)
func loadEnv() {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println("⚠️ Aucun fichier .env trouvé, on utilise les variables d'environnement système.")
		} else {
			log.Println("✅ Variables d'environnement chargées depuis .env")
		}
	})
}

// Get récupère une variable d'environnement, compatible Key Vault ET .env local
func Get(key string) string {
	loadEnv()

	// 1. Priorité à la version Key Vault (kebab-case)
	val := os.Getenv(key)
	if val != "" {
		return val
	}

	// 2. Fallback version UPPER_CASE (ex: pg-db-host → PG_DB_HOST)
	upperKey := strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
	val = os.Getenv(upperKey)
	if val != "" {
		return val
	}

	log.Printf("⚠️ Variable d'environnement %s introuvable (ni %s).\n", key, upperKey)
	return ""
}
