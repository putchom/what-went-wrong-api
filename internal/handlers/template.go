package handlers

import (
	"errors"
	"net/http"
	"what-went-wrong-api/internal/models"
	"what-went-wrong-api/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TemplateHandler struct {
	db *gorm.DB
}

func NewTemplateHandler(db *gorm.DB) *TemplateHandler {
	return &TemplateHandler{db: db}
}

// GetTemplates godoc
// @Summary List excuse templates
// @Description Get excuse templates. Can filter by pack_id. Premium users can access all. Free users restricted from premium packs.
// @Tags templates
// @Accept json
// @Produce json
// @Param pack_id query string false "Pack ID to filter"
// @Success 200 {object} GetTemplatesResponse
// @Router /excuse-templates [get]
func (h *TemplateHandler) GetTemplates(c *gin.Context) {
	packID := c.Query("pack_id")

	// Entitlement check
	entitlementsInterface, exists := c.Get("entitlements")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Entitlements not found"})
		return
	}
	entitlements := entitlementsInterface.(services.Entitlements)

	query := h.db.Model(&models.ExcuseTemplate{})

	if packID != "" {
		// If requesting a premium pack, check entitlement
		// For now, let's assume packs starting with "premium-" are premium, or we check specific IDs.
		// OR, we can rely on `IsPremium` flag on the template if we had one.
		// The model `ExcuseTemplate` has `IsPremium`. Let's use that logic?
		// Wait, `ExcuseTemplate` in `models/excuse.go` has `IsPremium`.
		// If filtering by pack_id, we should check if *that pack* contains premium templates?
		// A simpler logic: If user is not premium, exclude `is_premium = true`.

		query = query.Where("pack_id = ?", packID)
	}

	if !entitlements.CanUsePremiumTemplates {
		// Enforce free user restriction: Cannot see premium templates
		query = query.Where("is_premium = ?", false)
	}

	var templates []models.ExcuseTemplate
	if err := query.Find(&templates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch templates"})
		return
	}

	res := GetTemplatesResponse{Templates: make([]ExcuseTemplateResponse, len(templates))}
	for i, t := range templates {
		// Assuming Tags is stored as JSON or string array in DB, GORM handles it if configured
		// In `models/excuse.go`, tags is `type:text[]`. Gorm's postgres driver handles `lib/pq` array?
		// Or we need a scanner. For now, let's assume it works or is handled.
		// Looking at `models/excuse.go`, it's `pq.StringArray` likely.
		// Let's verify `models/excuse.go` again. It says `Tags []string gorm:"type:text[]"`.

		res.Templates[i] = ExcuseTemplateResponse{
			ID:         t.ID,
			PackID:     t.PackID,
			ExcuseText: t.Text,
			Tags:       t.Tags,
			IsPremium:  t.IsPremium,
			CreatedAt:  t.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, res)
}

// GetTemplate godoc
// @Summary Get template details
// @Tags templates
// @Accept json
// @Produce json
// @Param id path string true "Template ID"
// @Success 200 {object} ExcuseTemplateResponse
// @Router /excuse-templates/{id} [get]
func (h *TemplateHandler) GetTemplate(c *gin.Context) {
	id := c.Param("id")

	entitlementsInterface, exists := c.Get("entitlements")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Entitlements not found"})
		return
	}
	entitlements := entitlementsInterface.(services.Entitlements)

	var t models.ExcuseTemplate
	if err := h.db.First(&t, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch template"})
		return
	}

	if t.IsPremium && !entitlements.CanUsePremiumTemplates {
		c.JSON(http.StatusForbidden, gin.H{"error": "This template requires a premium plan"})
		return
	}

	res := ExcuseTemplateResponse{
		ID:         t.ID,
		PackID:     t.PackID,
		ExcuseText: t.Text,
		Tags:       t.Tags,
		IsPremium:  t.IsPremium,
		CreatedAt:  t.CreatedAt,
	}
	c.JSON(http.StatusOK, res)
}
