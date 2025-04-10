package security

import (
	"auth-service/config"
	"auth-service/database"
	"auth-service/models"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// 🔑 Lecture centralisée des secrets JWT
var (
	accessSecretKey  = []byte(config.Get("jwt-access-secret"))
	refreshSecretKey = []byte(config.Get("jwt-refresh-secret"))
)

// 🔐 Middleware JWT : vérifie le token et injecte les claims dans ctx
func JWTMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || len(authHeader) < 7 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing or invalid token"})
	}

	tokenStr := authHeader[7:] // Supprimer "Bearer "
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return accessSecretKey, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
	}

	// Injecter les claims dans le contexte
	c.Locals("user_id", uint(claims["user_id"].(float64)))
	c.Locals("email", claims["email"])
	c.Locals("role", claims["role"])
	c.Locals("permissions", claims["permissions"])

	return c.Next()
}

// 🔐 Génère access + refresh token avec les permissions intégrées
func GenerateTokens(user models.User) (string, string, error) {
	// Charger les permissions
	permissions, err := getPermissions(user.ID)
	if err != nil {
		return "", "", err
	}

	// Mapper les permissions
	var perms []map[string]string
	for _, perm := range permissions {
		perms = append(perms, map[string]string{
			"service": perm.ServiceID,
			"action":  perm.Action,
		})
	}

	// 🎫 Access token (15 minutes)
	accessClaims := jwt.MapClaims{
		"user_id":     user.ID,
		"email":       user.Email,
		"role":        user.Role,
		"permissions": perms,
		"exp":         time.Now().Add(15 * time.Minute).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(accessSecretKey)
	if err != nil {
		return "", "", err
	}

	// 🔁 Refresh token (7 jours)
	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(refreshSecretKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// 🔍 Parse un refresh token
func ParseRefreshToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return refreshSecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("refresh token invalide ou expiré")
	}
	return token, nil
}

// 🔎 Récupérer les permissions en base
func getPermissions(userID uint) ([]models.Permission, error) {
	var permissions []models.Permission
	err := database.DB.Where("user_id = ?", userID).Find(&permissions).Error
	return permissions, err
}
