package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var Rdb *redis.Client

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // à adapter si tu as un container ou une URL externe
		Password: "",               // si protégé, mets ton mot de passe ici
		DB:       0,
	})
	fmt.Println("🔗 Redis client initialized")
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

	fmt.Printf("📣 Published permission update for user %s: %v\n", userID, permissions)
	return nil
}
