package routes

import (
	"database/sql"
	"net/http"

	"github.com/dis70rt/streamoid/internal/database"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, sqlDB *sql.DB) {
	db := database.NewDatabaseService(sqlDB)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status":"healthy"})
	})

	router.POST("/upload", UploadCSV(db))
	router.GET("/products", GetAllProducts(db))
	router.GET("/products/search", SearchProducts(db))
}