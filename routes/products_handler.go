package routes

import (
	"net/http"

	"github.com/dis70rt/streamoid/internal/database"
	"github.com/gin-gonic/gin"
)

func GetAllProducts(db *database.DatabaseService) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req database.ProductRequest

        if err := c.ShouldBindQuery(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
            return
        }

        req.Default()

        products, err := db.GetAll(req.Page, req.Limit)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"products": products})
    }
}


func SearchProducts(db *database.DatabaseService) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req database.ProductRequest

        if err := c.ShouldBindQuery(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
            return
        }

        req.Default()

        products, err := db.Search(&req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"products": products})
    }
}
