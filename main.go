package main

import (
	"github.com/dis70rt/streamoid/internal/database"
	"github.com/dis70rt/streamoid/limiter"
	"github.com/dis70rt/streamoid/logger"
	"github.com/dis70rt/streamoid/routes"
	"github.com/gin-gonic/gin"
)

func main() {
    logger.Init()

    dbConfig, err := database.NewPSQL()
    if err != nil {
    	logger.Error("Error loading DB config", err)
		return
    }

    db, err := dbConfig.Connect()
    if err != nil {
		logger.Error("Failed to connect to database", err)
		return
    }
    defer db.Close()

    router := gin.Default()

    rateLimiter := limiter.GetLimiterManager(1,5)
    router.Use(rateLimiter.Middleware())

    routes.RegisterRoutes(router, db)
	  router.Run(":8080")
}
