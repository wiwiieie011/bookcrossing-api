package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dasler-fw/bookcrossing/internal/dto"
	"github.com/dasler-fw/bookcrossing/internal/middleware"
	"github.com/dasler-fw/bookcrossing/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type UserHandler struct {
	userServ services.UserService
	Redis    *redis.Client
}

func NewUserHandler(userServ services.UserService) *UserHandler {
	return &UserHandler{userServ: userServ}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.POST("/register", h.Register)
		users.POST("/login", h.Login)
		users.GET("/:id", middleware.JWTAuth(), h.GetProfile)
		users.PATCH("/:id", middleware.JWTAuth(), h.UpdateProfile)
		users.GET("/:id/exchanges", middleware.JWTAuth(), h.GetUserExchanges)
		// Collection endpoints
		users.GET("", h.List)       // GET /users
		users.GET("/", h.List)      // GET /users/ (trailing slash)
		users.GET("/list", h.List1) // legacy/simple list
	}

}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := h.userServ.Register(req)
	if err != nil {
		if err.Error() == "email уже используется" {
			c.JSON(http.StatusConflict, gin.H{
				"error": "email уже используется",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "не удалось зарегистрировать пользователя",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
	})

}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := h.userServ.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) GetProfile(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор пользователя"})
		return
	}

	profile, err := h.userServ.GetProfile(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	authUserID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор пользователя"})
		return
	}

	if authUserID != uint(id) {
		c.JSON(http.StatusForbidden, gin.H{"error": "доступ запрещён: нельзя редактировать чужой профиль"})
		return
	}

	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректное тело запроса"})
		return
	}
	if _, err := h.userServ.GetUserByID(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "пользователь не найден",
		})
		return
	}

	if err := h.userServ.UpdateProfile(uint(id), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "не удалось обновить профиль пользователя",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "профиль пользователя успешно обновлён",
	})

}

func (h *UserHandler) GetUserExchanges(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор пользователя"})
		return
	}

	status := c.Query("status")

	exchanges, err := h.userServ.GetUserExchanges(uint(id), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить историю обменов"})
		return
	}

	c.JSON(http.StatusOK, exchanges)
}

func (h *UserHandler) List1(c *gin.Context) {
	list, err := h.userServ.List()
	if err != nil {
		c.IndentedJSON(http.StatusHTTPVersionNotSupported, err)
		return
	}

	c.IndentedJSON(http.StatusOK, list)
}

func (h *UserHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	var lastID uint
	if lastIDStr := c.Query("last_id"); lastIDStr != "" {
		id, _ := strconv.ParseUint(lastIDStr, 10, 64)
		lastID = uint(id)
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("users:%d:%d", lastID, limit)
	nocache := c.Query("nocache") == "1"

	// 1️⃣ Проверяем кэш
	if !nocache && h.Redis != nil {
		if cached, err := h.Redis.Get(ctx, cacheKey).Result(); err == nil {
			c.Data(200, "application/json", []byte(cached))
			return
		}
	}

	// 2️⃣ Если нет в кэше — запрос из Postgres
	users, nextID, err := h.userServ.ListUsers(limit, lastID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	resp := gin.H{
		"data": users,
		"meta": gin.H{
			"limit":    limit,
			"next_id":  nextID,
			"has_next": nextID != 0,
		},
	}

	// optional debug info
	if c.Query("debug") == "1" {
		resp["debug"] = gin.H{
			"last_id":  lastID,
			"returned": len(users),
		}
	}

	jsonData, _ := json.Marshal(resp)

	// 3️⃣ Сохраняем в Redis на 5 минут (если кэш не отключён)
	if !nocache && h.Redis != nil {
		h.Redis.Set(ctx, cacheKey, jsonData, 5*time.Minute)
	}

	c.Data(200, "application/json", jsonData)
}
