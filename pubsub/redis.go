package pubsub

import (
	"auth-service/config"
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var Rdb *redis.Client

func InitRedis() {
	redisHost := config.Get("redis-host")     // ex: xpensify-redis.redis.cache.windows.net:6379
	redisPass := config.Get("redis-password") // fourni dans Azure -> Clés d'accès

	Rdb = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPass,
		DB:       0,
	})

	// Test connexion
	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("❌ Erreur de connexion à Redis : %v", err))
	}

	fmt.Println("🔗 Redis connecté à Azure avec succès")
}

func PublishPermissionUpdate(userID string, permissions []string) error {
	event := map[string]interface{}{
		"user_id":     userID,
		"permissions": permissions,
	}

	jsonData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = Rdb.Publish(ctx, "permission_events", jsonData).Err()
	if err != nil {
		return err
	}

	fmt.Printf("📣 Permission update envoyée pour user %s: %v\n", userID, permissions)
	return nil
}
