package repository_test

import (
	"io"
	"log/slog"
	"testing"

	"github.com/glebarez/sqlite" // драйвер от Глебареза
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/dasler-fw/bookcrossing/internal/dto"
	"github.com/dasler-fw/bookcrossing/internal/models"
	"github.com/dasler-fw/bookcrossing/internal/repository"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(
		sqlite.Open("file::memory:?cache=shared"),
		&gorm.Config{},
	)
	require.NoError(t, err)

	err = db.AutoMigrate(
		&models.User{},
		&models.Genre{},
		&models.Book{},
		&models.Review{},
		&models.Exchange{},
	)
	require.NoError(t, err)

	return db
}

// *********************************************************************************
// *						  Тесты для user									   *
// *								  |											   *
// *								  V									   		   *
// *********************************************************************************

func TestUserRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo := repository.NewUserRepository(db, log)

	//  Create
	user := &models.User{
		Name:         "Alice",
		Email:        "alice@example.com",
		PasswordHash: "hash",
		City:         "Moscow",
		Address:      "Lenina 1",
	}
	err := repo.Create(user)
	require.NoError(t, err)
	require.NotZero(t, user.ID)
	//  GetByID
	got, err := repo.GetByID(user.ID)
	require.NoError(t, err)
	require.Equal(t, user.Email, got.Email)

	//  GetByEmail
	gotByEmail, err := repo.GetByEmail("alice@example.com")
	require.NoError(t, err)
	require.Equal(t, user.ID, gotByEmail.ID)

	//  Update
	newName := "Alice Updated"
	user.Name = newName
	err = repo.Update(user)
	require.NoError(t, err)
	got, _ = repo.GetByID(user.ID)
	require.Equal(t, newName, got.Name)

	//  ListUsers
	users, err := repo.ListUsers(10, 0)
	require.NoError(t, err)
	require.Len(t, users, 1)

	//  Delete
	err = repo.Delete(user.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(user.ID)
	require.ErrorIs(t, err, repository.ErrUserNotFound)
}

// *********************************************************************************
// *						  Тесты для book									   *
// *								  |											   *
// *								  V									   		   *
// *********************************************************************************

func TestBookRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	repo := repository.NewBookRepository(db, log)

	//  Создаём пользователя
	user := &models.User{
		Name:         "Alice",
		Email:        "alice@example.com",
		PasswordHash: "hash",
	}
	require.NoError(t, db.Create(user).Error)
	require.NotZero(t, user.ID) // важно! иначе Book.UserID будет 0

	//  Создаём жанр
	genre := &models.Genre{Name: "Fiction"}
	require.NoError(t, db.Create(genre).Error)
	require.NotZero(t, genre.ID)

	//  Create книги
	book := &models.Book{
		Title:       "Test Book",
		Author:      "Author 1",
		Description: "Description",
		Status:      "available",
		UserID:      user.ID,
		Genres:      []models.Genre{*genre},
	}
	require.NoError(t, repo.Create(book))
	require.NotZero(t, book.ID)

	//  GetByID и проверка полей
	got, err := repo.GetByID(book.ID)
	require.NoError(t, err)
	require.Equal(t, "Test Book", got.Title)
	require.Equal(t, "Description", got.Description)
	require.Equal(t, user.ID, got.UserID)
	require.Len(t, got.Genres, 1)
	require.Equal(t, "Fiction", got.Genres[0].Name)

	//  Update книги
	book.Title = "Updated Title"
	require.NoError(t, repo.Update(book))

	got, err = repo.GetByID(book.ID)
	require.NoError(t, err)
	require.Equal(t, "Updated Title", got.Title)

	//  Delete книги
	require.NoError(t, repo.Delete(book.ID))

	_, err = repo.GetByID(book.ID)
	require.ErrorIs(t, err, dto.ErrorBookNotFound)
}

// *********************************************************************************
// *						  Тесты для review									   *
// *								  |											   *
// *								  V									   		   *
// *********************************************************************************

func TestReviewRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo := repository.NewReviewRepository(db, log)

	// Создаём пользователей
	author := &models.User{Name: "Author", Email: "author@example.com", PasswordHash: "hash"}
	targetUser := &models.User{Name: "Target", Email: "target@example.com", PasswordHash: "hash"}
	require.NoError(t, db.Create(author).Error)
	require.NoError(t, db.Create(targetUser).Error)

	// Создаём книгу
	book := &models.Book{Title: "Book 1", Author: "Author1", Status: "available", UserID: targetUser.ID}
	require.NoError(t, db.Create(book).Error)

	//  Create Review
	review := &models.Review{
		AuthorID:     author.ID,
		TargetUserID: targetUser.ID,
		TargetBookID: book.ID,
		Text:         "Great book!",
		Rating:       5,
	}
	require.NoError(t, repo.Create(review))
	require.NotZero(t, review.ID)

	//  GetByID
	got, err := repo.GetByID(review.ID)
	require.NoError(t, err)
	require.Equal(t, review.Text, got.Text)
	require.Equal(t, review.Rating, got.Rating)

	//  GetByTargetUserID
	reviewsByUser, err := repo.GetByTargetUserID(targetUser.ID)
	require.NoError(t, err)
	require.Len(t, reviewsByUser, 1)
	require.Equal(t, review.ID, reviewsByUser[0].ID)

	//  GetByTargetBookID
	reviewsByBook, err := repo.GetByTargetBookID(book.ID)
	require.NoError(t, err)
	require.Len(t, reviewsByBook, 1)
	require.Equal(t, review.ID, reviewsByBook[0].ID)

	//  Delete
	require.NoError(t, repo.Delete(review.ID))
	_, err = repo.GetByID(review.ID)
	require.ErrorIs(t, err, dto.ErrReviewNotFound)
}

// *********************************************************************************
// *						  Тесты для genre									   *
// *								  |											   *
// *								  V									   		   *
// *********************************************************************************

func TestGenreRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo := repository.NewGenreRepository(db, log)

	// Create
	genre := &models.Genre{Name: "classic"}
	require.NoError(t, repo.Create(genre))
	require.NotZero(t, genre.ID)

	// GetByID
	got, err := repo.GetByID(genre.ID)
	require.NoError(t, err)
	require.Equal(t, genre.Name, got.Name)

	// GetByName
	gotByName, err := repo.GetByName("classic")
	require.NoError(t, err)
	require.Equal(t, genre.ID, gotByName.ID)

	// Delete
	require.NoError(t, repo.Delete(genre.ID))

	// After delete → not found
	_, err = repo.GetByID(genre.ID)
	require.Error(t, err)
}


// *********************************************************************************
// *						  Тесты для exchange								   *
// *								  |											   *
// *								  V									   		   *
// *********************************************************************************

func TestExchangeRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo := repository.NewExchangeRepository(db, log)

	//  Создаём инициатора и получателя
	initiator := &models.User{Name: "Initiator", Email: "initiator@example.com", PasswordHash: "hash"}
	recipient := &models.User{Name: "Recipient", Email: "recipient@example.com", PasswordHash: "hash"}
	require.NoError(t, db.Create(initiator).Error)
	require.NoError(t, db.Create(recipient).Error)

	//  Создаём книги
	book1 := &models.Book{Title: "Book1", Author: "Author1", Status: "available", UserID: initiator.ID}
	book2 := &models.Book{Title: "Book2", Author: "Author2", Status: "available", UserID: recipient.ID}
	require.NoError(t, db.Create(book1).Error)
	require.NoError(t, db.Create(book2).Error)

	//  Create Exchange
	exchange := &models.Exchange{
		InitiatorID:     initiator.ID,
		RecipientID:     recipient.ID,
		InitiatorBookID: book1.ID,
		RecipientBookID: book2.ID,
		Status:          "pending",
	}
	require.NoError(t, repo.CreateExchange(exchange))
	require.NotZero(t, exchange.ID)

	// Проверяем статус книг после создания обмена
	var b1, b2 models.Book
	require.NoError(t, db.First(&b1, book1.ID).Error)
	require.NoError(t, db.First(&b2, book2.ID).Error)
	require.Equal(t, "reserved", b1.Status)
	require.Equal(t, "reserved", b2.Status)

	//  GetByID
	got, err := repo.GetByID(exchange.ID)
	require.NoError(t, err)
	require.Equal(t, exchange.ID, got.ID)
	require.Equal(t, exchange.Status, got.Status)

	//  CompleteExchange
	require.NoError(t, repo.CompleteExchange(got))

	// Проверяем, что книги обновились
	require.NoError(t, db.First(&b1, book1.ID).Error)
	require.NoError(t, db.First(&b2, book2.ID).Error)
	require.Equal(t, recipient.ID, b1.UserID)
	require.Equal(t, initiator.ID, b2.UserID)
	require.Equal(t, "available", b1.Status)
	require.Equal(t, "available", b2.Status)

	//  CancelExchange (на другом обмене)
	exchange2 := &models.Exchange{
		InitiatorID:     initiator.ID,
		RecipientID:     recipient.ID,
		InitiatorBookID: book1.ID,
		RecipientBookID: book2.ID,
		Status:          "pending",
	}
	require.NoError(t, repo.CreateExchange(exchange2))
	require.NoError(t, repo.CancelExchange(exchange2))

	// Проверка, что статус отменён
	got2, err := repo.GetByID(exchange2.ID)
	require.NoError(t, err)
	require.Equal(t, "cancelled", got2.Status)
	require.Nil(t, got2.CompletedAt)

	//  Delete Exchange
	require.NoError(t, db.Delete(&models.Exchange{}, exchange.ID).Error)
	_, err = repo.GetByID(exchange.ID)
	require.Error(t, err)
}