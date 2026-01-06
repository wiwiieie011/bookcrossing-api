package mocks

import (
	"github.com/dasler-fw/bookcrossing/internal/dto"
	"github.com/dasler-fw/bookcrossing/internal/models"
	"github.com/stretchr/testify/mock"
)

type BookServiceMock struct {
	mock.Mock
}

func (m *BookServiceMock) CreateBook(userID uint, req dto.CreateBookRequest) (*models.Book, error) {
	args := m.Called(userID, req) // передаём параметры в testify.Mock

	// Проверяем, что первый аргумент возвращённый не nil
	var book *models.Book
	if args.Get(0) != nil {
		book = args.Get(0).(*models.Book)
	}

	return book, args.Error(1) // второй аргумент — это ошибка
}

func (m *BookServiceMock) Update(bookID uint, userID uint, req dto.UpdateBookRequest) (*models.Book, error) {
	args := m.Called(bookID, userID, req)

	var book *models.Book
	if args.Get(0) != nil {
		book = args.Get(0).(*models.Book)
	}
	return book, args.Error(1)
}

func (m *BookServiceMock) Delete(bookID uint, userID uint) error {
	args := m.Called(bookID, userID)
	return args.Error(0)
}

func (m *BookServiceMock) SearchBooks(query dto.BookListQuery) ([]models.Book, int64, error) {
	args := m.Called(query)
	return args.Get(0).([]models.Book), args.Get(1).(int64), args.Error(2)

}

func (m *BookServiceMock) GetBooksByUserID(userID uint, status string) ([]models.Book, error) {
	args := m.Called(userID, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Book), args.Error(1)
}

func (m *BookServiceMock) GetAvailableBooks(city string) ([]models.Book, error) {
	args := m.Called(city)
	var books []models.Book
	if args.Get(0) != nil {
		books = args.Get(0).([]models.Book)
	}

	return books, args.Error(1)
}
