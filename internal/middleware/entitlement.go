package middleware

import (
	"net/http"
	"what-went-wrong-api/internal/models"
	"what-went-wrong-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EntitlementManager interface {
	GetPlan(userID uuid.UUID) (*models.UserPlan, error)
	GetEntitlements(planName string) services.Entitlements
}

func NewEntitlementMiddleware(service EntitlementManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		plan, err := service.GetPlan(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user plan"})
			c.Abort()
			return
		}

		entitlements := service.GetEntitlements(plan.Plan)

		c.Set("entitlements", entitlements)
		c.Next()
	}
}
