package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"what-went-wrong-api/internal/models"
	"what-went-wrong-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostGoals(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		db, cleanup := SetupTestDB(t)
		defer cleanup()
		handler := NewGoalHandler(db)

		userID := "auth0|test"

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		// Entitlement limit > 0
		c.Set("entitlements", services.Entitlements{MaxGoals: 3})

		reqBody := CreateGoalRequest{
			Title: "New Goal",
		}
		jsonBytes, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest("POST", "/goals", bytes.NewBuffer(jsonBytes))

		handler.PostGoals(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp CreateGoalResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "New Goal", resp.Goal.Title)
		assert.NotEmpty(t, resp.Goal.ID)
	})

	t.Run("LimitReached", func(t *testing.T) {
		db, cleanup := SetupTestDB(t)
		defer cleanup()
		handler := NewGoalHandler(db)

		userID := "auth0|test"
		// Pre-populate 3 goals
		for i := 0; i < 3; i++ {
			db.Create(&models.Goal{UserID: userID, Title: "Existing Goal"})
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		// Entitlement limit = 3
		c.Set("entitlements", services.Entitlements{MaxGoals: 3})

		reqBody := CreateGoalRequest{
			Title: "Fourth Goal",
		}
		jsonBytes, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest("POST", "/goals", bytes.NewBuffer(jsonBytes))

		handler.PostGoals(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestGetGoals(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("List", func(t *testing.T) {
		db, cleanup := SetupTestDB(t)
		defer cleanup()
		handler := NewGoalHandler(db)

		userID := "auth0|test"
		db.Create(&models.Goal{UserID: userID, Title: "Goal 1", Order: 2, CreatedAt: time.Now().Add(-time.Hour)})
		db.Create(&models.Goal{UserID: userID, Title: "Goal 2", Order: 1, CreatedAt: time.Now()})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)

		c.Request, _ = http.NewRequest("GET", "/goals", nil)

		handler.GetGoals(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp GetGoalsResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp.Goals, 2)
		// Validation order: Order asc, CreatedAt desc
		// Goal 2 has order 1, Goal 1 has order 2. So Goal 2 should be first.
		assert.Equal(t, "Goal 2", resp.Goals[0].Title)
	})
}

func TestDeleteGoal(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		db, cleanup := SetupTestDB(t)
		defer cleanup()
		handler := NewGoalHandler(db)

		userID := "auth0|test"
		goal := models.Goal{UserID: userID, Title: "To Delete"}
		db.Create(&goal)

		// Create associated excuse
		excuse := models.ExcuseEntry{UserID: userID, GoalID: goal.ID, ExcuseText: "test", Date: "2025-01-01"}
		err := db.Create(&excuse).Error
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		r.DELETE("/goals/:id", handler.DeleteGoal)

		req, _ := http.NewRequest("DELETE", "/goals/"+goal.ID.String(), nil)
		// Need to set userID in context. Since we use router, we need a middleware or hack?
		// We can use a middleware to set userID for testing
		r.Use(func(c *gin.Context) {
			c.Set("userID", userID)
			c.Next()
		})

		// ServeHTTP doesn't seem to make it easy to inject middleware AFTER defining route if route is already added?
		// Actually gin adds middleware to group.
		// Let's rebuild router structure clearly.

		r2 := gin.New()
		r2.Use(func(c *gin.Context) {
			c.Set("userID", userID)
			c.Next()
		})
		r2.DELETE("/goals/:id", handler.DeleteGoal)

		r2.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify deletion
		var count int64
		db.Model(&models.Goal{}).Where("id = ?", goal.ID).Count(&count)
		assert.Equal(t, int64(0), count)

		db.Model(&models.ExcuseEntry{}).Where("id = ?", excuse.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}
