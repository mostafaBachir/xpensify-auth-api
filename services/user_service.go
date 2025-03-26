package services

import (
	"auth-service/database"
	"auth-service/models"
	"errors"
)

// 🔹 Liste tous les utilisateurs avec leurs permissions et services
func GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := database.DB.Preload("Permissions.Service").Find(&users).Error
	return users, err
}

// 🔹 Met à jour les permissions d’un utilisateur
func UpdateUserPermissions(userID uint, permissions []models.Permission) error {
	// Vérifie que l’utilisateur existe
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return errors.New("utilisateur introuvable")
	}

	// Supprime les anciennes permissions
	if err := database.DB.Where("user_id = ?", userID).Delete(&models.Permission{}).Error; err != nil {
		return err
	}

	// Ajoute les nouvelles permissions
	for i := range permissions {
		permissions[i].UserID = userID
	}
	return database.DB.Create(&permissions).Error
}
