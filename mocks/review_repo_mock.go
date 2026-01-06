package mocks

import (
	"github.com/dasler-fw/bookcrossing/internal/models"
	"github.com/stretchr/testify/mock"
)

// Create(req *models.Review) error
// 	GetByID(id uint) (*models.Review, error)
// 	Delete(id uint) error
// 	GetByTargetUserID(id uint) ([]models.Review, error)
// 	GetByTargetBookID(id uint) ([]models.Review, error)

type ReviewRepositoryMock struct {
	mock.Mock
}

func (m *ReviewRepositoryMock) Create(req *models.Review) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *ReviewRepositoryMock) GetByID(id uint) (*models.Review, error) {
	args := m.Called(id)

	var r *models.Review
	if args.Get(0) != nil {
		r = args.Get(0).(*models.Review)
	}

	return r, args.Error(1)
}

func (m *ReviewRepositoryMock) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *ReviewRepositoryMock) GetByTargetUserID(id uint) ([]models.Review, error) {
	args := m.Called(id)
	var r []models.Review
	if args.Get(0) != nil {
		r = args.Get(0).([]models.Review)
	}

	return r, args.Error(1)
}

func (m *ReviewRepositoryMock) GetByTargetBookID(id uint) ([]models.Review, error) {
	args := m.Called(id)
	var r []models.Review
	if args.Get(0) != nil {
		r = args.Get(0).([]models.Review)
	}

	return r, args.Error(1)
}
