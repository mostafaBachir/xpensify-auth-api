package services

import (
	"auth-service/database"
	"auth-service/models"
)

func GetAllServices() ([]models.Service, error) {
	var services []models.Service
	err := database.DB.Find(&services).Error
	return services, err
}

func CreateService(service *models.Service) error {
	return database.DB.Create(service).Error
}

func DeleteService(id uint) error {
	return database.DB.Delete(&models.Service{}, id).Error
}
func UpdateService(id uint, updated *models.Service) error {
	var service models.Service
	if err := database.DB.First(&service, id).Error; err != nil {
		return err
	}

	service.Name = updated.Name
	service.Url = updated.Url

	return database.DB.Save(&service).Error
}
