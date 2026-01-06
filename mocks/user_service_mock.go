package mocks

import (
	"github.com/dasler-fw/bookcrossing/internal/dto"
	"github.com/dasler-fw/bookcrossing/internal/models"
	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func (m *UserServiceMock) GetProfile(id uint) (*dto.UserProfileResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserProfileResponse), args.Error(1)
}

func (m *UserServiceMock) GetUserByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserServiceMock) UpdateUser(id uint, req dto.UserUpdateRequest) (*models.User, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserServiceMock) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *UserServiceMock) Register(req dto.UserCreateRequest) (string, error) {
	args := m.Called(req)
	return args.String(0), args.Error(1)
}

func (m *UserServiceMock) Login(req dto.LoginRequest) (string, error) {
	args := m.Called(req)
	return args.String(0), args.Error(1)
}

func (m *UserServiceMock) GetUserExchanges(userID uint, status string) ([]models.Exchange, error) {
	args := m.Called(userID, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Exchange), args.Error(1)
}

func (m *UserServiceMock) ListUsers(limit int, lastID uint) ([]models.User, uint, error) {
	args := m.Called(limit, lastID)
	return args.Get(0).([]models.User), args.Get(1).(uint), args.Error(2)
}
