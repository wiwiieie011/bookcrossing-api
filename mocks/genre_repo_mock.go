package mocks

import (
	"github.com/dasler-fw/bookcrossing/internal/models"
	"github.com/stretchr/testify/mock"
)

type GenreRepositoryMock struct {
	mock.Mock
}

func (m *GenreRepositoryMock) Create(req *models.Genre) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *GenreRepositoryMock) GetByID(id uint) (*models.Genre, error) {
	args := m.Called(id)

	var g *models.Genre
	if args.Get(0) != nil {
		g = args.Get(0).(*models.Genre)
	}

	return g, args.Error(1)
}

func (m *GenreRepositoryMock) GetByName(name string) (*models.Genre, error) {
	args := m.Called(name)
	var g *models.Genre
	if args.Get(0) != nil {
		g = args.Get(0).(*models.Genre)
	}
	return g, args.Error(1)
}

func (m *GenreRepositoryMock) List() ([]models.Genre, error) {
	args := m.Called()

	var genres []models.Genre
	if args.Get(0) != nil {
		genres = args.Get(0).([]models.Genre)
	}

	return genres, args.Error(1)
}

func (m *GenreRepositoryMock) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
