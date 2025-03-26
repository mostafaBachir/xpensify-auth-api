package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string       `json:"name" gorm:"not null"`
	Email        string       `json:"email" gorm:"unique;not null"`
	Password     string       `json:"-"`
	RefreshToken string       `json:"-"`
	Role         string       `json:"role" gorm:"not null;default:'user'"`
	Permissions  []Permission `json:"permissions" gorm:"foreignKey:UserID"` // 🟢 Relation simple et safe
}
