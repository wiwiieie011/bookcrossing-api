package transport

import (
	"log/slog"

	"github.com/dasler-fw/bookcrossing/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(
	router *gin.Engine,
	log *slog.Logger,
	bookService services.BookService,
	exchangeService services.ExchangeService,
	genreService services.GenreService,
	reviewService services.ReviewService,
	userService services.UserService,
	rdb *redis.Client,
) {
	bookHandler := NewBookHandler(bookService)
	exchangeHandler := NewExchangeHandler(exchangeService)
	genreHandler := NewGenreHandler(genreService)
	reviewHandler := NewReviewHandler(reviewService)
	userHandler := NewUserHandler(userService)
	// wire Redis client for handlers that use caching
	userHandler.Redis = rdb

	bookHandler.RegisterRoutes(router)
	exchangeHandler.RegisterExchangeRoutes(router)
	genreHandler.RegisterGenreRoutes(router)
	reviewHandler.RegisterReviewRoutes(router)
	userHandler.RegisterRoutes(router)
}
