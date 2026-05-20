package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gokost710/subscription-service/internal/config"
	"github.com/gokost710/subscription-service/internal/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.Log.Level)

	router := gin.New()
	router.Use(gin.Recovery())

}
