package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"what-went-wrong-api/internal/models"
	"what-went-wrong-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetMePlan(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success_DefaultFree", func(t *testing.T) {
		db, cleanup := SetupTestDB(t)
		defer cleanup()

		entitlementService := services.NewEntitlementService(db)
		handler := NewPlanHandler(entitlementService)

		userID := uuid.New()
		// No plan pre-seeded, should default to free

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID.String())

		c.Request, _ = http.NewRequest("GET", "/me/plan", nil)
		handler.GetMePlan(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp GetMePlanResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "free", resp.Plan)
		assert.Equal(t, 3, resp.Entitlements.MaxGoals)

		// Verify DB
		var plan models.UserPlan
		db.Where("user_id = ?", userID).First(&plan)
		assert.Equal(t, "free", plan.Plan)
	})

	t.Run("Success_ExistingPremium", func(t *testing.T) {
		db, cleanup := SetupTestDB(t)
		defer cleanup()

		entitlementService := services.NewEntitlementService(db)
		handler := NewPlanHandler(entitlementService)

		userID := uuid.New()
		db.Create(&models.UserPlan{UserID: userID, Plan: "premium"})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID.String())

		c.Request, _ = http.NewRequest("GET", "/me/plan", nil)
		handler.GetMePlan(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp GetMePlanResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "premium", resp.Plan)
		assert.Equal(t, 100, resp.Entitlements.MaxGoals)
	})
}

func TestPostMePlan(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("UpdateToPremium", func(t *testing.T) {
		db, cleanup := SetupTestDB(t)
		defer cleanup()

		entitlementService := services.NewEntitlementService(db)
		handler := NewPlanHandler(entitlementService)

		userID := uuid.New()
		// Initial: Free
		db.Create(&models.UserPlan{UserID: userID, Plan: "free"})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID.String())

		reqBody := PostMePlanRequest{Plan: "premium"}
		jsonBytes, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest("POST", "/me/plan", bytes.NewBuffer(jsonBytes))

		handler.PostMePlan(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp PostMePlanResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "premium", resp.Plan)
		assert.Equal(t, 100, resp.Entitlements.MaxGoals)

		// Verify DB
		var plan models.UserPlan
		db.Where("user_id = ?", userID).First(&plan)
		assert.Equal(t, "premium", plan.Plan)
	})
}
