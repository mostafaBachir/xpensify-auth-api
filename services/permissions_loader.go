package services

import (
	"auth-service/database"
	"auth-service/models"
)

func GetPermissionsAsMap(userID uint) ([]map[string]string, error) {
	var permissions []models.Permission
	if err := database.DB.Where("user_id = ?", userID).Find(&permissions).Error; err != nil {
		return nil, err
	}

	var perms []map[string]string
	for _, perm := range permissions {
		perms = append(perms, map[string]string{
			"service": perm.ServiceID,
			"action":  perm.Action,
		})
	}
	return perms, nil
}
