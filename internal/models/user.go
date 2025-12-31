package models

import (
	"time"

	"gorm.io/gorm"
)

// User модель с составным уникальным индексом по (email, deleted_at),
// чтобы можно было повторно использовать email после soft delete.
type User struct {
	ID        uint           `json:"-" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index:idx_users_email_deleted_at"`

	Name         string `json:"name"`
	Email        string `json:"email" gorm:"uniqueIndex:idx_users_email_deleted_at;not null"`
	PasswordHash string `json:"-"`
	City         string `json:"city"`
	Address      string `json:"address"`
}
