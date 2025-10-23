package routes

import (
	"net/http"
	"os"

	"github.com/dis70rt/streamoid/internal/database"
	"github.com/dis70rt/streamoid/logger"
	"github.com/gin-gonic/gin"
)

func UploadCSV(db *database.DatabaseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
			return
		}
		
		if gin.Mode() == gin.TestMode {
            fileHandle, err := file.Open()
            if err != nil {
                logger.Error("failed to open uploaded file", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
                return
            }
            defer fileHandle.Close()

            totalCount, failedProducts, err := db.ProcessCSVFromReader(fileHandle)
            if err != nil {
                logger.Error("Failed to process CSV", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }

            c.JSON(http.StatusOK, gin.H{"stored": totalCount, "failed": failedProducts})
            return
        }
	

		tempPath := "/temp/" + file.Filename
		if err := c.SaveUploadedFile(file, tempPath); err != nil {
			logger.Error("failed to save file", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		defer os.Remove(tempPath)

		totalCount, failedProducts, err := db.UploadCSV("/temp/" + file.Filename)
		if err != nil {
			logger.Error("Failed to get newProducts", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"stored": totalCount, "failed": failedProducts}) 
	}
}