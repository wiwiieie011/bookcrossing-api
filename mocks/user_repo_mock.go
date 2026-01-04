package mocks

import (
	"github.com/dasler-fw/bookcrossing/internal/models"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepositoryMock) GetByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepositoryMock) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepositoryMock) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepositoryMock) ListUsers(limit int, lastID uint) ([]models.User, error) {
	args := m.Called(limit, lastID)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *UserRepositoryMock) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}


func (m *UserRepositoryMock) GetUserExchanges(userID uint, status string) ([]models.Exchange, error) {
	args := m.Called(userID, status)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]models.Exchange), args.Error(1)
}





