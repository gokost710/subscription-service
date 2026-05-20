package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gokost710/subscription-service/internal/http/handler"
	"github.com/gokost710/subscription-service/internal/service"
)

func New(subscriptionService service.SubscriptionProvider) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())

	healthHandler := handler.NewHealthHandler()
	r.GET("/health", healthHandler.Check)

	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)

	subscriptions := r.Group("/subscriptions")
	{
		subscriptions.POST("", subscriptionHandler.Create)
		subscriptions.GET("/:id", subscriptionHandler.GetByID)
	}

	return r
}
