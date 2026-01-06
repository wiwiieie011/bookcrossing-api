package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dasler-fw/bookcrossing/internal/dto"
	"github.com/dasler-fw/bookcrossing/internal/models"
	"github.com/dasler-fw/bookcrossing/internal/transport"
	"github.com/dasler-fw/bookcrossing/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupGin() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func TestUserHandler_GetProfile_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 🔹 mock сервиса
	userService := new(mocks.UserServiceMock)

	// 🔹 handler
	handler := transport.NewUserHandler(userService)

	// 🔹 ожидаемый профиль
	profile := &dto.UserProfileResponse{
		ID:                       1,
		Name:                     "Alice",
		City:                     "Moscow",
		BooksCount:               0,
		SuccessfulExchangesCount: 0,
	}

	// 🔹 ожидание вызова сервиса
	userService.
		On("GetProfile", uint(1)).
		Return(profile, nil)

	// 🔹 HTTP запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)

	r := gin.New()
	r.GET("/users/:id", handler.GetProfile)

	// 🔹 выполняем запрос
	r.ServeHTTP(w, req)

	// 🔹 проверки
	require.Equal(t, http.StatusOK, w.Code)

	var resp dto.UserProfileResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	require.Equal(t, profile.ID, resp.ID)
	require.Equal(t, profile.Name, resp.Name)

	userService.AssertExpectations(t)
}

func TestUserHandler_Register_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userService := new(mocks.UserServiceMock)
	handler := transport.NewUserHandler(userService)

	reqBody := dto.UserCreateRequest{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "pass",
		City:     "Moscow",
		Address:  "Lenina 1",
	}

	userService.On("Register", reqBody).Return("TOKEN", nil)

	b, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	r := gin.New()
	r.POST("/users/register", handler.Register)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var body map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.Equal(t, "TOKEN", body["token"])

	userService.AssertExpectations(t)
}

func TestUserHandler_Login_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userService := new(mocks.UserServiceMock)
	handler := transport.NewUserHandler(userService)

	reqBody := dto.LoginRequest{Email: "bob@example.com", Password: "pass"}
	userService.On("Login", reqBody).Return("TOKEN", nil)

	b, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	r := gin.New()
	r.POST("/users/login", handler.Login)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var body map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.Equal(t, "TOKEN", body["token"])
}

func TestUserHandler_UpdateProfile_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userService := new(mocks.UserServiceMock)
	handler := transport.NewUserHandler(userService)

	name := "New Name"
	updReq := dto.UserUpdateRequest{Name: &name}

	userService.On("GetUserByID", uint(1)).Return(&models.User{ID: 1}, nil).Maybe()
	// Use Anything for request to avoid deep equal issues with pointer fields
	userService.On("UpdateUser", uint(1), mock.Anything).Return(&models.User{ID: 1, Name: "New Name"}, nil)

	b, _ := json.Marshal(updReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/users/1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	r := gin.New()
	// middleware to inject user_id into context
	r.PATCH("/users/:id", func(c *gin.Context) { c.Set("user_id", uint(1)) }, handler.UpdateProfile)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_List_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userService := new(mocks.UserServiceMock)
	handler := transport.NewUserHandler(userService)

	// no Redis configured in handler => falls back to service
	users := []models.User{{ID: 1, Name: "A"}}
	userService.On("ListUsers", 50, uint(0)).Return(users, uint(0), nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)

	r := gin.New()
	r.GET("/users", handler.List)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var body struct {
		Data []models.User `json:"data"`
		Meta struct {
			Limit   int  `json:"limit"`
			NextID  uint `json:"next_id"`
			HasNext bool `json:"has_next"`
		} `json:"meta"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.Len(t, body.Data, 1)
	require.Equal(t, 50, body.Meta.Limit)
	require.False(t, body.Meta.HasNext)
}
