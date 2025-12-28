package repository

import (
	"errors"
	"log/slog"

	"github.com/dasler-fw/bookcrossing/internal/dto"
	"github.com/dasler-fw/bookcrossing/internal/models"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("пользователь не найден")

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	Update(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	ListUsers(limit int, lastID uint) ([]models.User, error)
	ListUsers1() ([]models.User, error)
	Delete(id uint) error
}

type userRepository struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewUserRepository(db *gorm.DB, log *slog.Logger) UserRepository {
	return &userRepository{
		db:  db,
		log: log,
	}
}

func (r *userRepository) Create(user *models.User) error {
	if user == nil {
		r.log.Error("ошибка создания профиля")
		return dto.ErrUserCreateFailed
	}
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		r.log.Error("ошибка получения пользователя", "id", id, "err", err)

		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, dto.ErrUserGetFailed
	}
	return &user, nil

}

func (r *userRepository) Update(user *models.User) error {
	if user == nil || user.ID == 0 {
		r.log.Error("ошибка обновления: пустой профиль или отсутствует ID")
		return dto.ErrUserUpdateFailed
	}

	return r.db.Save(user).Error

}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		r.log.Error("ошибка получения профиля по Email")

		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, dto.ErrUserGetFailed
	}
	return &user, nil
}

func (r *userRepository) ListUsers(limit int, lastID uint) ([]models.User, error) {
	var users []models.User

	q := r.db.
		Table("users"). // явно указываем таблицу
		Order("id ASC").
		Limit(limit)

	if lastID > 0 {
		q = q.Where("id > ?", lastID)
	}

	if err := q.Find(&users).Error; err != nil {
		r.log.Error("ошибка получения пользователей", "err", err)
		return nil, err
	}

	return users, nil
}

func (r *userRepository) ListUsers1() ([]models.User, error) {
	var users []models.User

	q := r.db.
		Table("users").
		Order("id ASC")

	if err := q.Find(&users).Error; err != nil {
		r.log.Error("ошибка получения пользователей", "err", err)
		return nil, dto.ErrUserGetFailed
	}

	return users, nil
}

func (r *userRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.User{}, id).Error; err != nil {
		r.log.Error("ошибка удаления профиля")
		return dto.ErrUserDeleteFailed
	}
	return nil
}
