package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"what-went-wrong-api/internal/models"
	"what-went-wrong-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetTemplates(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("FreeUser_FilterPremium", func(t *testing.T) {
		db, cleanup := SetupTestDB(t)
		defer cleanup()

		// Ensure Migration for ExcuseTemplate
		db.AutoMigrate(&models.ExcuseTemplate{})

		handler := NewTemplateHandler(db)

		// Seed templates
		freeTmpl := models.ExcuseTemplate{
			ID:        "tmpl-1",
			PackID:    "pack-1",
			Text:      "Free Excuse",
			IsPremium: false,
			Tags:      pq.StringArray{"funny", "work"},
		}
		premiumTmpl := models.ExcuseTemplate{
			ID:        "tmpl-2",
			PackID:    "pack-premium",
			Text:      "Premium Excuse",
			IsPremium: true,
			Tags:      pq.StringArray{"vip"},
		}
		db.Create(&freeTmpl)
		db.Create(&premiumTmpl)

		userID := uuid.New()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID.String())
		// Free user
		c.Set("entitlements", services.Entitlements{CanUsePremiumTemplates: false})
		c.Request, _ = http.NewRequest("GET", "/excuse-templates", nil)

		handler.GetTemplates(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp GetTemplatesResponse
		json.Unmarshal(w.Body.Bytes(), &resp)

		// Should only see free template
		assert.Len(t, resp.Templates, 1)
		assert.Equal(t, "Free Excuse", resp.Templates[0].ExcuseText)
		assert.Contains(t, resp.Templates[0].Tags, "funny")
	})

	t.Run("PremiumUser_SeeAll", func(t *testing.T) {
		db, cleanup := SetupTestDB(t)
		defer cleanup()
		db.AutoMigrate(&models.ExcuseTemplate{})
		handler := NewTemplateHandler(db)

		// Seed templates
		freeTmpl := models.ExcuseTemplate{ID: "t1", PackID: "pack-1", IsPremium: false}
		premiumTmpl := models.ExcuseTemplate{ID: "t2", PackID: "pack-premium", IsPremium: true}
		db.Create(&freeTmpl)
		db.Create(&premiumTmpl)

		userID := uuid.New()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID.String())
		// Premium user
		c.Set("entitlements", services.Entitlements{CanUsePremiumTemplates: true})
		c.Request, _ = http.NewRequest("GET", "/excuse-templates", nil)

		handler.GetTemplates(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp GetTemplatesResponse
		json.Unmarshal(w.Body.Bytes(), &resp)

		assert.Len(t, resp.Templates, 2)
	})

	t.Run("FilterByPackID", func(t *testing.T) {
		db, cleanup := SetupTestDB(t)
		defer cleanup()
		db.AutoMigrate(&models.ExcuseTemplate{})
		handler := NewTemplateHandler(db)

		t1 := models.ExcuseTemplate{ID: "t1", PackID: "pack-1"}
		t2 := models.ExcuseTemplate{ID: "t2", PackID: "pack-2"}
		db.Create(&t1)
		db.Create(&t2)

		userID := uuid.New()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID.String())
		c.Set("entitlements", services.Entitlements{CanUsePremiumTemplates: true})
		// Filter query
		c.Request, _ = http.NewRequest("GET", "/excuse-templates?pack_id=pack-1", nil)

		handler.GetTemplates(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp GetTemplatesResponse
		json.Unmarshal(w.Body.Bytes(), &resp)

		assert.Len(t, resp.Templates, 1)
		assert.Equal(t, "pack-1", resp.Templates[0].PackID)
	})
}
