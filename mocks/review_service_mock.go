package mocks

import (
	"github.com/dasler-fw/bookcrossing/internal/dto"
	"github.com/dasler-fw/bookcrossing/internal/models"
	"github.com/stretchr/testify/mock"
)

// Create(authorID uint, req dto.CreateReviewRequest) error
// 	GetByUserID(userID uint) ([]models.Review, error)
// 	GetByBookID(bookID uint) ([]models.Review, error)
// 	Delete(reviewID uint, authorID uint) error

type ReviewServiceMock struct {
	mock.Mock
}

func (m *ReviewServiceMock) Create(authorID uint, req dto.CreateReviewRequest) error {
	args := m.Called(authorID, req)
	return args.Error(0)
}

func (m *ReviewServiceMock) GetByUserID(userID uint) ([]models.Review, error) {
	args := m.Called(userID)

	var r []models.Review
	if args.Get(0) != nil {
		r = args.Get(0).([]models.Review)
	}

	return r, args.Error(1)
}

func (m *ReviewServiceMock) GetByBookID(bookID uint) ([]models.Review, error) {
	args := m.Called(bookID)

	var r []models.Review
	if args.Get(0) != nil {
		r = args.Get(0).([]models.Review)
	}

	return r, args.Error(1)
}

func (m *ReviewServiceMock) Delete(reviewID uint, authorID uint) error {
	args := m.Called(reviewID,authorID)
	return args.Error(0)
}
