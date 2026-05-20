package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gokost710/subscription-service/internal/http/handler"
)

func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())

	healthHandler := handler.NewHealthHandler()
	r.GET("/health", healthHandler.Check)

	return r
}
