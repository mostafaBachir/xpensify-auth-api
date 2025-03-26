package security

import (
	"auth-service/database"
	"auth-service/models"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))
var RefreshSecretKey = []byte(os.Getenv("JWT_REFRESH_SECRET"))

// 🔐 Génère access + refresh token avec les permissions intégrées
func GenerateTokens(user models.User) (string, string, error) {
	// Charger les permissions
	permissions, err := getPermissions(user.ID)
	if err != nil {
		return "", "", err
	}

	// 🔁 Transformer les permissions en structure simple
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
		"exp":         time.Now().Add(time.Minute * 15).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	// 🔁 Refresh token (7 jours)
	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(RefreshSecretKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// 🔍 Parse & vérifie un refresh token
func ParseRefreshToken(refreshToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Signature HMAC requise
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Méthode de signature invalide")
		}
		return RefreshSecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("Refresh token invalide ou expiré")
	}

	return token, nil
}
func getPermissions(userID uint) ([]models.Permission, error) {
	var permissions []models.Permission
	err := database.DB.Where("user_id = ?", userID).Find(&permissions).Error
	return permissions, err
}
