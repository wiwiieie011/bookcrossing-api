package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`
	Name         string `json:"name"`
	Email        string `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string `json:"-"`
	City         string `json:"-"`
	Address      string `json:"-"`
}
