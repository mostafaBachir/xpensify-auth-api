package services

import (
	"auth-service/database"
	"auth-service/models"
	"auth-service/security"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// 📌 Inscription d'un utilisateur
func RegisterUser(user *models.User) error {
	// Vérifie si l'email est déjà utilisé
	var existingUser models.User
	if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return errors.New("email already in use")
	}

	// Hash du mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error hashing password")
	}
	user.Password = string(hashedPassword)

	// Enregistrement en base de données
	return database.DB.Create(user).Error
}

// 📌 Connexion utilisateur
func LoginUser(email, password string) (string, string, error) {
	var user models.User

	// Charger l'utilisateur avec ses permissions
	if err := database.DB.Preload("Permissions").Where("email = ?", email).First(&user).Error; err != nil {
		return "", "", errors.New("invalid email or password")
	}

	// Vérification du mot de passe
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("invalid email or password")
	}

	// Génération des tokens
	accessToken, refreshToken, err := security.GenerateTokens(user)
	if err != nil {
		return "", "", errors.New("could not generate tokens")
	}

	// Sauvegarde du Refresh Token
	user.RefreshToken = refreshToken
	if err := database.DB.Save(&user).Error; err != nil {
		return "", "", errors.New("could not save refresh token")
	}

	return accessToken, refreshToken, nil
}

// 📌 Rafraîchissement de token
func RefreshUserToken(refreshToken string) (string, string, error) {
	// Valider le Refresh Token
	token, err := security.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	// Extraire les claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	// Récupérer l'utilisateur
	userID := uint(claims["user_id"].(float64))
	var user models.User
	if err := database.DB.Preload("Permissions").First(&user, userID).Error; err != nil || user.RefreshToken != refreshToken {
		return "", "", errors.New("invalid refresh token")
	}

	// Générer les nouveaux tokens
	accessToken, newRefreshToken, err := security.GenerateTokens(user)
	if err != nil {
		return "", "", errors.New("could not generate new tokens")
	}

	// Mettre à jour en BDD
	user.RefreshToken = newRefreshToken
	database.DB.Save(&user)

	return accessToken, newRefreshToken, nil
}

// 📌 Déconnexion
func LogoutUser(userID uint) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	user.RefreshToken = ""
	return database.DB.Save(&user).Error
}

// 📌 Récupérer un utilisateur par son ID
func GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
