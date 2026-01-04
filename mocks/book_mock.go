package mocks

import (
	"github.com/dasler-fw/bookcrossing/internal/dto"
	"github.com/dasler-fw/bookcrossing/internal/models"
	"github.com/stretchr/testify/mock"
)

type BookRepositoryMock struct {
	mock.Mock
}

func (m *BookRepositoryMock) Create(req *models.Book) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *BookRepositoryMock) GetByID(id uint) (*models.Book, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *BookRepositoryMock) Update(book *models.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *BookRepositoryMock) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}


func (m *BookRepositoryMock) Search(query dto.BookListQuery) ([]models.Book, int64, error) {
	args := m.Called(query)

	var books []models.Book
	if args.Get(0) != nil {
		books = args.Get(0).([]models.Book)
	}

	return books, args.Get(1).(int64), args.Error(2)
}


func (m *BookRepositoryMock) AttachGenres(bookID uint, genreIDs []uint) error {
	args := m.Called(bookID, genreIDs)
	return args.Error(0)
}


func (m *BookRepositoryMock) GetByUserID(userID uint, status string) ([]models.Book, error) {
	args := m.Called(userID, status)

	var books []models.Book
	if args.Get(0) != nil {
		books = args.Get(0).([]models.Book)
	}

	return books, args.Error(1)
}


 func (m *BookRepositoryMock) GetAvailable(city string) ([]models.Book, error) {
	args := m.Called(city)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	return args.Get(0).([]models.Book), args.Error(1)
 }