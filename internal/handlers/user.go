package handlers

import (
	"net/http"
	"what-went-wrong-api/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ErrorResponse struct {
	Error string `json:"error" example:"missing or invalid Authorization header"`
}

// GetUsers godoc
// @Summary Get all users
// @Description Get a list of all users from database. Requires Bearer token authentication.
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Security BearerAuth
// @Router /users [get]
func GetUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		if err := db.Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}
