package test

import "github.com/gin-gonic/gin"

func setupGin() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}
