package services

import (
	"auth-service/database"
	"auth-service/models"
	"auth-service/security"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// 📌 Fonction pour l'inscription d'un utilisateur
func RegisterUser(user *models.User) error {
	// Vérifier si l'email est déjà utilisé
	var existingUser models.User
	if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return errors.New("email already in use")
	}

	// Hasher le mot de passe avec bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error hashing password")
	}
	user.Password = string(hashedPassword)

	// Sauvegarder l'utilisateur en BDD
	return database.DB.Create(user).Error
}

// 📌 Fonction pour la connexion utilisateur
func LoginUser(email, password string) (string, string, error) {
	var user models.User

	// Vérifier si l'utilisateur existe
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", "", errors.New("invalid email or password")
	}

	// Vérifier le mot de passe
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("invalid email or password")
	}

	// Générer Access Token & Refresh Token
	accessToken, refreshToken, err := security.GenerateTokens(user.ID)
	if err != nil {
		return "", "", errors.New("could not generate tokens")
	}

	// Sauvegarde du Refresh Token en BDD
	user.RefreshToken = refreshToken
	if err := database.DB.Save(&user).Error; err != nil {
		return "", "", errors.New("could not save refresh token")
	}

	return accessToken, refreshToken, nil
}

// 📌 Fonction pour rafraîchir un token expiré
func RefreshUserToken(refreshToken string) (string, string, error) {
	// 📌 Vérifier et parser le Refresh Token
	token, err := security.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	// 📌 Extraire les claims (payload du token)
	claims, ok := token.Claims.(jwt.MapClaims) // ✅ Correction ici
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	// 📌 Récupérer l'ID utilisateur
	userID := uint(claims["user_id"].(float64))
	var user models.User

	// 📌 Vérifier si l'utilisateur existe et si le token correspond à celui en BDD
	if err := database.DB.First(&user, userID).Error; err != nil || user.RefreshToken != refreshToken {
		return "", "", errors.New("invalid refresh token")
	}

	// 📌 Générer un nouveau couple de tokens
	accessToken, newRefreshToken, err := security.GenerateTokens(userID)
	if err != nil {
		return "", "", errors.New("could not generate new tokens")
	}

	// 📌 Mettre à jour le Refresh Token en base de données
	user.RefreshToken = newRefreshToken
	database.DB.Save(&user)

	return accessToken, newRefreshToken, nil
}

// 📌 Fonction pour la déconnexion (révocation du Refresh Token)
func LogoutUser(userID uint) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// Supprimer le Refresh Token en base de données
	user.RefreshToken = ""
	return database.DB.Save(&user).Error
}

// 📌 Fonction pour récupérer un utilisateur par son ID
func GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
