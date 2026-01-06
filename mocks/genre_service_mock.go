package mocks

import (
	"github.com/dasler-fw/bookcrossing/internal/dto"
	"github.com/dasler-fw/bookcrossing/internal/models"
	"github.com/stretchr/testify/mock"
)


type GenreServiceMock struct {
	mock.Mock
}

func (m *GenreServiceMock) Create(req dto.GenreCreateRequest) (*models.Genre, error) {
	args := m.Called(req)
	var g *models.Genre

	if args.Get(0) != nil {
		g = args.Get(0).(*models.Genre)
	}

	return g, args.Error(1)
}


func (m *GenreServiceMock) GetByID(id uint) (*models.Genre, error) {
	args := m.Called(id)

	var g *models.Genre
	if args.Get(0) != nil {
		g = args.Get(0).(*models.Genre)
	}

	return g, args.Error(1)
}


func (m *GenreServiceMock) List() ([]models.Genre, error) {
	args := m.Called()

	var genres []models.Genre
	if args.Get(0) != nil {
		genres = args.Get(0).([]models.Genre)
	}

	return genres, args.Error(1)
}

func (m *GenreServiceMock) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}


