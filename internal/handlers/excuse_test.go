package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"what-went-wrong-api/internal/models"
	"what-went-wrong-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetExcuses_Retention(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	handler := NewExcuseHandler(db)

	userID := "auth0|test"
	goalID := uuid.New()

	// Seed Old Excuse (40 days ago)
	oldDate := time.Now().AddDate(0, 0, -40).Format("2006-01-02")
	db.Create(&models.ExcuseEntry{UserID: userID, GoalID: goalID, Date: oldDate, ExcuseText: "Old"})

	// Seed New Excuse (Today)
	newDate := time.Now().Format("2006-01-02")
	db.Create(&models.ExcuseEntry{UserID: userID, GoalID: goalID, Date: newDate, ExcuseText: "New"})

	t.Run("FreeUser_Restricted", func(t *testing.T) {
		days := 30
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Set("entitlements", services.Entitlements{LogRetentionDays: &days})
		c.Params = gin.Params{{Key: "id", Value: goalID.String()}}
		c.Request, _ = http.NewRequest("GET", "/goals/"+goalID.String()+"/excuses", nil)

		handler.GetExcuses(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp GetExcusesResponse
		json.Unmarshal(w.Body.Bytes(), &resp)

		assert.Len(t, resp.Excuses, 1)
		assert.Equal(t, "New", resp.Excuses[0].ExcuseText)
	})

	t.Run("PremiumUser_SeeAll", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Set("entitlements", services.Entitlements{LogRetentionDays: nil}) // Unlimited
		c.Params = gin.Params{{Key: "id", Value: goalID.String()}}
		c.Request, _ = http.NewRequest("GET", "/goals/"+goalID.String()+"/excuses", nil)

		handler.GetExcuses(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp GetExcusesResponse
		json.Unmarshal(w.Body.Bytes(), &resp)

		assert.Len(t, resp.Excuses, 2)
	})
}

func TestPostExcuse_Upsert(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	handler := NewExcuseHandler(db)
	userID := "auth0|test"
	goalID := uuid.New()
	today := time.Now().Format("2006-01-02")

	// 1. Create New
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", userID)
	c.Set("entitlements", services.Entitlements{})
	c.Params = gin.Params{{Key: "id", Value: goalID.String()}}

	reqBody := `{"date": "` + today + `", "excuseText": "First Excuse"}`
	c.Request, _ = http.NewRequest("POST", "/goals/"+goalID.String()+"/excuses", strings.NewReader(reqBody))

	handler.PostExcuse(c)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Verify in DB
	var count int64
	db.Model(&models.ExcuseEntry{}).Where("user_id = ? AND goal_id = ?", userID, goalID).Count(&count)
	assert.Equal(t, int64(1), count)

	// 2. Update Existing
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Set("userID", userID)
	c.Set("entitlements", services.Entitlements{})
	c.Params = gin.Params{{Key: "id", Value: goalID.String()}}

	reqBody = `{"date": "` + today + `", "excuseText": "Updated Excuse"}`
	c.Request, _ = http.NewRequest("POST", "/goals/"+goalID.String()+"/excuses", strings.NewReader(reqBody))

	handler.PostExcuse(c)
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify Update
	var entry models.ExcuseEntry
	db.Where("user_id = ? AND goal_id = ?", userID, goalID).First(&entry)
	assert.Equal(t, "Updated Excuse", entry.ExcuseText)
}

func TestPostExcuse_PremiumTemplate(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	handler := NewExcuseHandler(db)
	userID := "auth0|test"
	goalID := uuid.New()

	// Create Premium Template
	db.AutoMigrate(&models.ExcuseTemplate{})
	db.Create(&models.ExcuseTemplate{ID: "tmpl-premium", Text: "Premium", IsPremium: true})

	today := time.Now().Format("2006-01-02")

	// Free user try to use premium template
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", userID)
	c.Set("entitlements", services.Entitlements{CanUsePremiumTemplates: false})
	c.Params = gin.Params{{Key: "id", Value: goalID.String()}}

	reqBody := `{"date": "` + today + `", "excuseText": "Using Premium", "templateId": "tmpl-premium"}`
	c.Request, _ = http.NewRequest("POST", "/goals/"+goalID.String()+"/excuses", strings.NewReader(reqBody))

	handler.PostExcuse(c)
	assert.Equal(t, http.StatusForbidden, w.Code)
}
