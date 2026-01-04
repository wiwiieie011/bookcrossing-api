package test

import (
	"io"
	"log/slog"
	"testing"

	"github.com/dasler-fw/bookcrossing/internal/dto"
	"github.com/dasler-fw/bookcrossing/internal/models"
	"github.com/dasler-fw/bookcrossing/internal/services"
	"github.com/dasler-fw/bookcrossing/mocks"
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



// *********************************************************************************
// *						  Тесты для exchange								   *
// *								  |											   *
// *								  V									   		   *
// *********************************************************************************

// *********************************************************************************
// *						  Тесты для exchange								   *
// *								  |											   *
// *								  V									   		   *
// *********************************************************************************
