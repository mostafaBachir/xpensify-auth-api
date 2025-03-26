package models

import "gorm.io/gorm"

type Service struct {
	gorm.Model
	Name string `json:"name" gorm:"not null"`
	Url  string `json:"url" gorm:"not null"`
}
