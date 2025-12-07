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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Goal{}, &models.ExcuseEntry{})
	return db
}

func TestPostGoals(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		db := setupTestDB()
		handler := NewGoalHandler(db)

		userID := uuid.New()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID.String())
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
		db := setupTestDB()
		handler := NewGoalHandler(db)

		userID := uuid.New()
		// Pre-populate 3 goals
		for i := 0; i < 3; i++ {
			db.Create(&models.Goal{UserID: userID, Title: "Existing Goal"})
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID.String())
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
		db := setupTestDB()
		handler := NewGoalHandler(db)

		userID := uuid.New()
		db.Create(&models.Goal{UserID: userID, Title: "Goal 1", Order: 2, CreatedAt: time.Now().Add(-time.Hour)})
		db.Create(&models.Goal{UserID: userID, Title: "Goal 2", Order: 1, CreatedAt: time.Now()})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID.String())

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
		db := setupTestDB()
		handler := NewGoalHandler(db)

		userID := uuid.New()
		goal := models.Goal{UserID: userID, Title: "To Delete"}
		db.Create(&goal)

		// Create associated excuse
		excuse := models.ExcuseEntry{UserID: userID, GoalID: goal.ID, ExcuseText: "test"}
		db.Create(&excuse)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID.String())
		c.Params = []gin.Param{{Key: "id", Value: goal.ID.String()}}

		c.Request, _ = http.NewRequest("DELETE", "/goals/"+goal.ID.String(), nil)

		handler.DeleteGoal(c)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify deletion
		var count int64
		db.Model(&models.Goal{}).Where("id = ?", goal.ID).Count(&count)
		assert.Equal(t, int64(0), count)

		db.Model(&models.ExcuseEntry{}).Where("id = ?", excuse.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}
