package models

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model

	UserID uint `json:"user_id" gorm:"not null;primaryKey"`
	User   User `gorm:"constraint:OnDelete:CASCADE;"`

	ServiceID string  `json:"service_id" gorm:"not null;primaryKey"`
	Service   Service `gorm:"constraint:OnDelete:CASCADE;"`

	Action    string    `json:"action" gorm:"not null;primaryKey"` // ex: read, write, delete
	GrantedAt time.Time `json:"granted_at" gorm:"autoCreateTime"`
}
