package services

import (
	"auth-service/database"
	"auth-service/models"
	"errors"
)

func HasPermission(userID uint, serviceID string, action string) bool {
	var count int64
	database.DB.Model(&models.Permission{}).
		Where("user_id = ? AND service_id = ? AND action = ?", userID, serviceID, action).
		Count(&count)

	return count > 0
}

func GetPermissionsByUser(userID uint) ([]models.Permission, error) {
	var permissions []models.Permission
	result := database.DB.Where("user_id = ?", userID).Find(&permissions)
	return permissions, result.Error
}

func AssignPermission(p *models.Permission) error {
	var existing models.Permission
	if err := database.DB.First(&existing, "user_id = ? AND service_id = ? AND action = ?", p.UserID, p.ServiceID, p.Action).Error; err == nil {
		return errors.New("permission déjà existante")
	}
	return database.DB.Create(p).Error
}

func DeletePermission(userID uint, serviceID string, action string) error {
	return database.DB.
		Where("user_id = ? AND service_id = ? AND action = ?", userID, serviceID, action).
		Delete(&models.Permission{}).Error
}
