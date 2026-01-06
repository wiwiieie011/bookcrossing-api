package test

import (
	"io"
	"log/slog"
	"testing"

	"github.com/dasler-fw/bookcrossing/internal/dto"
	"github.com/dasler-fw/bookcrossing/internal/models"
	"github.com/dasler-fw/bookcrossing/internal/services"
	"github.com/dasler-fw/bookcrossing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// *********************************************************************************
// *						  Тесты для users								       *
// *								  |											   *
// *								  V									   		   *
// *********************************************************************************

func TestUserService_GetUserByID_OK(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	userRepo := new(mocks.UserRepositoryMock)
	bookRepo := new(mocks.BookRepositoryMock)

	svc := services.NewServiceUser(nil, userRepo, bookRepo, log)

	user := &models.User{
		ID:           1,
		Name:         "Alice",
		Email:        "alice@example.com",
		PasswordHash: "hashedpassword",
		City:         "Moscow",
		Address:      "Lenina 1",
	}

	userRepo.
		On("GetByID", uint(1)).
		Return(user, nil)

	got, err := svc.GetUserByID(1)

	require.NoError(t, err)
	require.Equal(t, user.ID, got.ID)
	require.Equal(t, "Alice", got.Name)

	userRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser_OK(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	userRepo := new(mocks.UserRepositoryMock)
	svc := services.NewServiceUser(nil, userRepo, nil, log)

	// Исходный пользователь
	user := &models.User{
		ID:           1,
		Name:         "Alice",
		Email:        "alice@example.com",
		PasswordHash: "hashedpassword",
		City:         "Moscow",
		Address:      "Lenina 1",
	}

	// Запрос на обновление
	newName := "Alice Updated"
	req := dto.UserUpdateRequest{
		Name: &newName,
	}

	// Моки: GetByID возвращает пользователя
	userRepo.On("GetByID", uint(1)).Return(user, nil)

	// Моки: Update возвращает nil (успех)
	userRepo.On("Update", user).Return(nil)

	got, err := svc.UpdateUser(1, req)

	require.NoError(t, err)
	require.Equal(t, "Alice Updated", got.Name)

	userRepo.AssertExpectations(t)
}

func TestUserService_Delete_OK(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	userRepo := new(mocks.UserRepositoryMock)
	bookRepo := new(mocks.BookRepositoryMock)

	svc := services.NewServiceUser(nil, userRepo, bookRepo, log)

	user := &models.User{
		ID:           1,
		Name:         "Alice",
		Email:        "alice@example.com",
		PasswordHash: "hashedpassword",
		City:         "Moscow",
		Address:      "Lenina 1",
	}

	// мок на GetByID
	userRepo.On("GetByID", uint(1)).Return(user, nil)
	// мок на Delete
	userRepo.On("Delete", uint(1)).Return(nil)

	err := svc.DeleteUser(1)
	require.NoError(t, err)

	userRepo.AssertExpectations(t)
}

func TestUserService_GetUserExchanges_OK(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	userRepo := new(mocks.UserRepositoryMock)
	bookRepo := new(mocks.BookRepositoryMock)

	svc := services.NewServiceUser(nil, userRepo, bookRepo, log)

	exchanges := []models.Exchange{
		{
			Model:  gorm.Model{ID: 1},
			Status: "pending",
		},
	}

	userRepo.
		On("GetUserExchanges", uint(1), "pending").
		Return(exchanges, nil)

	result, err := svc.GetUserExchanges(1, "pending")

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, uint(1), result[0].ID)

	userRepo.AssertExpectations(t)
}

// *********************************************************************************
// *						  Тесты для book								       *
// *								  |											   *
// *								  V									   		   *
// *********************************************************************************
func TestBookService_Create_OK(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	bookRepo := new(mocks.BookRepositoryMock)
	service := services.NewServiceBook(bookRepo, log)

	userID := uint(10)
	req := dto.CreateBookRequest{
		Title:       "Clean Code",
		Author:      "Robert Martin",
		Description: "About clean code",
		AISummary:   "How to write clean code",
	}

	bookRepo.
		On("Create", mock.Anything).
		Return(nil).
		Once()

	book, err := service.CreateBook(userID, req)

	require.NoError(t, err)
	require.NotNil(t, book)

	assert.Equal(t, req.Title, book.Title)
	assert.Equal(t, userID, book.UserID)

	bookRepo.AssertExpectations(t)
}

func TestBookService_GetByID_OK(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	bookRepo := new(mocks.BookRepositoryMock)
	svc := services.NewServiceBook(bookRepo, log)

	book := &models.Book{
		Model:       gorm.Model{ID: 1},
		Title:       "summer",
		Author:      "lev",
		Description: "sdladalsdlasdlsa",
		AISummary:   "dasdsadasdsa",
		Status:      "sasdasdsdsa",
		UserID:      1,
	}

	bookRepo.On("GetByID", uint(1)).Return(book, nil)
	got, err := svc.GetByID(1)
	require.NoError(t, err)
	require.Equal(t, book.ID, got.ID)
	require.Equal(t, "summer", got.Title)

	bookRepo.AssertExpectations(t)
}

func TestBookService_UpdateBook_OK(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	bookRepo := new(mocks.BookRepositoryMock)
	svc := services.NewServiceBook(bookRepo, log)

	book := &models.Book{
		Model:       gorm.Model{ID: 1},
		Title:       "summer",
		Author:      "lev",
		Description: "sdladalsdlasdlsa",
		AISummary:   "dasdsadasdsa",
		Status:      "sasdasdsdsa",
		UserID:      1,
	}

	descr := "go only up dont back"
	req := &dto.UpdateBookRequest{
		Description: &descr,
	}

	bookRepo.On("GetByID", uint(1)).Return(book, nil)
	bookRepo.On("Update", mock.Anything).Return(nil)
	got, err := svc.Update(1, 1, *req)

	require.NoError(t, err)
	require.Equal(t, descr, got.Description)

	bookRepo.AssertExpectations(t)

}

// Delete(bookID uint, userID uint) error
func TestBookService_DeleteBook_OK(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	bookRepo := new(mocks.BookRepositoryMock)
	svc := services.NewServiceBook(bookRepo, log)

	book := &models.Book{
		Model:       gorm.Model{ID: 1},
		Title:       "summer",
		Author:      "lev",
		Description: "sdladalsdlasdlsa",
		AISummary:   "dasdsadasdsa",
		Status:      "sasdasdsdsa",
		UserID:      1,
	}

	bookRepo.On("GetByID", uint(1)).Return(book, nil)
	bookRepo.On("Delete", uint(1)).Return(nil)
	err := svc.Delete(1, 1)
	require.NoError(t, err)
	bookRepo.AssertExpectations(t)
}

func TestBookService_SearchBooks_OK(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	bookRepo := new(mocks.BookRepositoryMock)
	svc := services.NewServiceBook(bookRepo, log)

	books := []models.Book{
		{
			Model:       gorm.Model{ID: 1},
			Title:       "summer",
			Author:      "lev",
			Description: "test description",
			AISummary:   "summary",
			Status:      "available",
			UserID:      1,
		},
		{
			Model:       gorm.Model{ID: 2},
			Title:       "winter",
			Author:      "anna",
			Description: "test description 2",
			AISummary:   "summary 2",
			Status:      "reserved",
			UserID:      2,
		},
	}

	var total int64 = int64(len(books))

	query := dto.BookListQuery{
		Title:     "summer",
		Page:      1,
		Limit:     2,
		SortBy:    "created_at",
		SortOrder: "desc",
	}

	bookRepo.On("Search", query).Return(books, total, nil)

	// 5️⃣ Вызываем метод сервиса
	gotBooks, gotTotal, err := svc.SearchBooks(query)

	// 6️⃣ Проверяем результат
	require.NoError(t, err)
	require.Equal(t, total, gotTotal)
	require.Len(t, gotBooks, 2)
	require.Equal(t, "summer", gotBooks[0].Title)

	bookRepo.AssertExpectations(t)
}

func TestBookService_GetBooksByUserID_OK(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	bookRepo := new(mocks.BookRepositoryMock)
	svc := services.NewServiceBook(bookRepo, log)

	books := []models.Book{
		{
			Model:       gorm.Model{ID: 1},
			Title:       "summer",
			Author:      "lev",
			Description: "test description",
			AISummary:   "summary",
			Status:      "available",
			UserID:      1,
		},
		{
			Model:       gorm.Model{ID: 2},
			Title:       "winter",
			Author:      "anna",
			Description: "test description 2",
			AISummary:   "summary 2",
			Status:      "reserved",
			UserID:      1,
		},
	}

	bookRepo.On("GetByUserID", uint(1), "available").Return(books, nil)

	got, err := svc.GetBooksByUserID(1, "available")

	require.NoError(t, err)
	require.Len(t, got, 2)
	require.Equal(t, "summer", got[0].Title)
	require.Equal(t, "winter", got[1].Title)

	bookRepo.AssertExpectations(t)
}

func TestBookService_GetAvailableBooks_OK(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	bookRepo := new(mocks.BookRepositoryMock)
	svc := services.NewServiceBook(bookRepo, log)

	books := []models.Book{
		{
			Model:       gorm.Model{ID: 1},
			Title:       "summer",
			Author:      "lev",
			Description: "test description",
			AISummary:   "summary",
			Status:      "available",
			UserID:      1,
		},
		{
			Model:       gorm.Model{ID: 2},
			Title:       "winter",
			Author:      "anna",
			Description: "test description 2",
			AISummary:   "summary 2",
			Status:      "reserved",
			UserID:      1,
		},
	}

	bookRepo.On("GetAvailable", "Moscow").Return(books, nil)

	got, err := svc.GetAvailableBooks("Moscow")
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.Equal(t, "summer", got[0].Title)
	require.Equal(t, "winter", got[1].Title)

	bookRepo.AssertExpectations(t)
}

// *********************************************************************************
// *						  Тесты для genre								  	   *
// *								  |											   *
// *								  V									   		   *
// *********************************************************************************

func TestGenreService_CreateGenre_OK(t *testing.T) {
	genreRepo := new(mocks.GenreRepositoryMock)
	svc := services.NewGenreService(genreRepo)

	req := &dto.GenreCreateRequest{
		Name: "Classic",
	}

	genreRepo.On("Create", mock.Anything).Return(nil)

	got, err := svc.Create(*req)

	require.NoError(t, err)
	require.Equal(t, req.Name, got.Name)

	genreRepo.AssertExpectations(t)
}

func TestGenreService_GetByIDGenre_OK(t *testing.T) {
	genreRepo := new(mocks.GenreRepositoryMock)
	svc := services.NewGenreService(genreRepo)

	genr := &models.Genre{
		Model: gorm.Model{ID: 1},
		Name:  "Classic",
	}

	genreRepo.On("GetByID", uint(1)).Return(genr, nil)

	got, err := svc.GetByID(1)
	require.NoError(t, err)
	require.Equal(t, genr.Name, got.Name)

	genreRepo.AssertExpectations(t)
}




// *********************************************************************************
// *						  Тесты для Review								       *
// *								  |											   *
// *								  V									   		   *
// *********************************************************************************




func TestReviewService_Create_OK(t *testing.T) {
	reviewRepo := new(mocks.ReviewRepositoryMock)
	svc := services.NewReviewService(reviewRepo)

	authorID := uint(1)
	req := dto.CreateReviewRequest{
		TargetUserID: 2,
		TargetBookID: 2,
		Text:         "bad bookjdjfjsdfj",
		Rating:       5,
	}

	reviewRepo.On("Create", mock.Anything).Return(nil)

	got, err := svc.Create(authorID, req)
	require.NoError(t, err)
	require.Equal(t, authorID, got.AuthorID)
	require.Equal(t, req.TargetUserID, got.TargetUserID)
	require.Equal(t, req.TargetBookID, got.TargetBookID)
	require.Equal(t, req.Text, got.Text)
	require.Equal(t, req.Rating, got.Rating)

	reviewRepo.AssertExpectations(t)
}

func TestReviewService_GetByTargetUserID_OK(t *testing.T) {
	reviewRepo := new(mocks.ReviewRepositoryMock)
	svc := services.NewReviewService(reviewRepo)

	review := []models.Review{
		{
			Model:        gorm.Model{ID: 1},
			AuthorID:     1,
			TargetUserID: 2,
			TargetBookID: 2,
			Text:         "t wiweir weirwe ieqwe",
			Rating:       2,
		},
		{
			Model:        gorm.Model{ID: 2},
			AuthorID:     1,
			TargetUserID: 2,
			TargetBookID: 3,
			Text:         "t wiweir weirwe ieqwe",
			Rating:       2,
		},
	}

	reviewRepo.On("GetByTargetUserID", uint(1)).Return(review, nil)

	got, err := svc.GetByUserID(uint(1))
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.Equal(t, review[0].Text, got[0].Text)
	require.Equal(t, review[1].Text, got[1].Text)
	require.Equal(t, review[0].Rating, got[0].Rating)
	require.Equal(t, review[1].Rating, got[1].Rating)
}
func TestReviewService_GetByTargetBookID_OK(t *testing.T) {
	reviewRepo := new(mocks.ReviewRepositoryMock)
	svc := services.NewReviewService(reviewRepo)

	review := []models.Review{
		{
			Model:        gorm.Model{ID: 1},
			AuthorID:     1,
			TargetUserID: 2,
			TargetBookID: 2,
			Text:         "t wiweir weirwe ieqwe",
			Rating:       2,
		},
		{
			Model:        gorm.Model{ID: 2},
			AuthorID:     2,
			TargetUserID: 2,
			TargetBookID: 2,
			Text:         "t wiweir weirwe ieqwe",
			Rating:       2,
		},
	}

	reviewRepo.On("GetByTargetBookID", uint(1)).Return(review, nil)

	got, err := svc.GetByBookID(uint(1))
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.Equal(t, review[0].Text, got[0].Text)
	require.Equal(t, review[1].Text, got[1].Text)
	require.Equal(t, review[0].Rating, got[0].Rating)
	require.Equal(t, review[1].Rating, got[1].Rating)
}

func TestReviewService_DeleteReview_OK(t *testing.T) {
	reviewRepo := new(mocks.ReviewRepositoryMock)
	svc := services.NewReviewService(reviewRepo)

	authorID := uint(1)
	review := &models.Review{
		Model:        gorm.Model{ID: 1},
		AuthorID:     1,
		TargetUserID: 2,
		TargetBookID: 2,
		Text:         "t wiweir weirwe ieqwe",
		Rating:       2,
	}

	reviewRepo.On("GetByID", uint(1)).Return(review, nil)
	reviewRepo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(uint(1), authorID)
	require.NoError(t, err)
}







// *********************************************************************************
// *						  Тесты для exchange								   *
// *								  |											   *
// *								  V									   		   *
// *********************************************************************************







func TestExchangeService_Create_OK(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	exchangeRepo := new(mocks.ExchangeRepositoryMock)
	bookRepo := new(mocks.BookRepositoryMock)

	svc := services.NewExchangeService(exchangeRepo, bookRepo, log)

	req := dto.CreateExchangeRequest{
		RecipientID:     2,
		InitiatorBookID: 10,
		RecipientBookID: 20,
	}

	// 📘 Книга инициатора
	initiatorBook := &models.Book{
		Model:  gorm.Model{ID: 10},
		UserID: 1,
		Status: "available",
	}

	// 📕 Книга получателя
	recipientBook := &models.Book{
		Model:  gorm.Model{ID: 20},
		UserID: 2,
		Status: "available",
	}

	bookRepo.On("GetByID", uint(10)).Return(initiatorBook, nil)
	bookRepo.On("GetByID", uint(20)).Return(recipientBook, nil)

	exchangeRepo.On("CreateExchange", mock.Anything).Return(nil)

	// ACT
	exchange, err := svc.CreateExchange(&req, 1)

	// ASSERT
	require.NoError(t, err)
	require.NotNil(t, exchange)
	require.Equal(t, uint(1), exchange.InitiatorID)
	require.Equal(t, uint(2), exchange.RecipientID)
	require.Equal(t, "pending", exchange.Status)

	bookRepo.AssertExpectations(t)
	exchangeRepo.AssertExpectations(t)
}

func TestExchangeService_CompleteExchange_OK(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	exchangeRepo := new(mocks.ExchangeRepositoryMock)
	bookRepo := new(mocks.BookRepositoryMock)

	svc := services.NewExchangeService(exchangeRepo, bookRepo, log)

	// Готовим обмен со статусом accepted
	exch := &models.Exchange{
		Model:           gorm.Model{ID: 1},
		InitiatorID:     1,
		RecipientID:     2,
		InitiatorBookID: 10,
		RecipientBookID: 20,
		Status:          "accepted",
	}

	// Ожидания моков
	exchangeRepo.On("GetByID", uint(1)).Return(exch, nil)
	exchangeRepo.On("CompleteExchange", exch).Return(nil)

	// Действие: инициатор завершает обмен
	err := svc.CompleteExchange(1, 1)

	// Проверка
	require.NoError(t, err)
	exchangeRepo.AssertExpectations(t)
}

func TestExchangeService_CancelExchange_OK(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	exchangeRepo := new(mocks.ExchangeRepositoryMock)
	bookRepo := new(mocks.BookRepositoryMock)

	svc := services.NewExchangeService(exchangeRepo, bookRepo, log)
	// Обмен в статусе pending может отменить только инициатор
	exch := &models.Exchange{
		Model:           gorm.Model{ID: 2},
		InitiatorID:     5,
		RecipientID:     7,
		InitiatorBookID: 10,
		RecipientBookID: 20,
		Status:          "pending",
	}

	exchangeRepo.On("GetByID", uint(2)).Return(exch, nil)
	exchangeRepo.On("CancelExchange", exch).Return(nil)

	err := svc.CancelExchange(2, 5)
	require.NoError(t, err)

	exchangeRepo.AssertExpectations(t)
}

func TestExchangeService_AcceptExchange_OK(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	exchangeRepo := new(mocks.ExchangeRepositoryMock)
	bookRepo := new(mocks.BookRepositoryMock)

	svc := services.NewExchangeService(exchangeRepo, bookRepo, log)

	// Принять pending может только получатель
	exch := &models.Exchange{
		Model:           gorm.Model{ID: 3},
		InitiatorID:     11,
		RecipientID:     22,
		InitiatorBookID: 101,
		RecipientBookID: 202,
		Status:          "pending",
	}

	exchangeRepo.On("GetByID", uint(3)).Return(exch, nil)
	// Проверяем, что статус поменяется на accepted и будет вызван Update
	exchangeRepo.On("Update", mock.MatchedBy(func(e *models.Exchange) bool {
		return e.ID == 3 && e.Status == "accepted"
	})).Return(nil)

	err := svc.AcceptExchange(3, 22)
	require.NoError(t, err)

	exchangeRepo.AssertExpectations(t)
}

func TestExchangeService_GetByID_OK(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	exchangeRepo := new(mocks.ExchangeRepositoryMock)
	bookRepo := new(mocks.BookRepositoryMock)

	svc := services.NewExchangeService(exchangeRepo, bookRepo, log)

	exch := &models.Exchange{Model: gorm.Model{ID: 99}, InitiatorID: 1, RecipientID: 2, Status: "pending"}
	exchangeRepo.On("GetByID", uint(99)).Return(exch, nil)

	got, err := svc.GetByID(99)
	require.NoError(t, err)
	require.Equal(t, uint(99), got.ID)

	exchangeRepo.AssertExpectations(t)
}

func TestExchangeService_GetAll_OK(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	exchangeRepo := new(mocks.ExchangeRepositoryMock)
	bookRepo := new(mocks.BookRepositoryMock)

	svc := services.NewExchangeService(exchangeRepo, bookRepo, log)

	list := []models.Exchange{
		{Model: gorm.Model{ID: 1}, InitiatorID: 1, RecipientID: 2, Status: "pending"},
		{Model: gorm.Model{ID: 2}, InitiatorID: 3, RecipientID: 4, Status: "accepted"},
	}

	exchangeRepo.On("GetAll").Return(list, nil)

	got, err := svc.GetAll()
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.Equal(t, uint(1), got[0].ID)
	require.Equal(t, uint(2), got[1].ID)

	exchangeRepo.AssertExpectations(t)
}
