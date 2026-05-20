package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gokost710/subscription-service/internal/http/handler"
	"github.com/gokost710/subscription-service/internal/http/swagger"
	"github.com/gokost710/subscription-service/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New(subscriptionService service.SubscriptionProvider) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	swagger.Register()

	r := gin.New()
	r.Use(gin.Recovery())

	healthHandler := handler.NewHealthHandler()
	r.GET("/health", healthHandler.Check)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)

	subscriptions := r.Group("/subscriptions")
	{
		subscriptions.POST("", subscriptionHandler.Create)
		subscriptions.GET("", subscriptionHandler.List)
		subscriptions.GET("/summary", subscriptionHandler.TotalPrice)
		subscriptions.GET("/:id", subscriptionHandler.GetByID)
		subscriptions.PUT("/:id", subscriptionHandler.Update)
		subscriptions.DELETE("/:id", subscriptionHandler.Delete)
	}

	return r
}
