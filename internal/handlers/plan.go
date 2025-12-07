package handlers

import (
	"net/http"
	"what-went-wrong-api/internal/models"
	"what-went-wrong-api/internal/services"

	"github.com/gin-gonic/gin"
)

type EntitlementManager interface {
	GetPlan(userID string) (*models.UserPlan, error)
	GetEntitlements(planName string) services.Entitlements
	UpdatePlan(userID string, planName string) (*models.UserPlan, error)
}

type PlanHandler struct {
	entitlementService EntitlementManager
}

func NewPlanHandler(entitlementService EntitlementManager) *PlanHandler {
	return &PlanHandler{entitlementService: entitlementService}
}

// GetMePlan godoc
// @Summary Get current user plan and entitlements
// @Description Returns the user's current subscription plan and their active entitlements.
// @Tags plan
// @Accept json
// @Produce json
// @Success 200 {object} GetMePlanResponse
// @Router /me/plan [get]
func (h *PlanHandler) GetMePlan(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDStr.(string)

	plan, err := h.entitlementService.GetPlan(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get plan"})
		return
	}

	entitlements := h.entitlementService.GetEntitlements(plan.Plan)

	c.JSON(http.StatusOK, GetMePlanResponse{
		Plan:         plan.Plan,
		Entitlements: entitlements,
	})
}

// PostMePlan godoc
// @Summary Update user plan
// @Description Updates the user's subscription plan (e.g., from 'free' to 'premium').
// @Tags plan
// @Accept json
// @Produce json
// @Param request body PostMePlanRequest true "Request body"
// @Success 200 {object} PostMePlanResponse
// @Router /me/plan [post]
func (h *PlanHandler) PostMePlan(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDStr.(string)

	var req PostMePlanRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPlan, err := h.entitlementService.UpdatePlan(userID, req.Plan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update plan"})
		return
	}

	entitlements := h.entitlementService.GetEntitlements(updatedPlan.Plan)

	c.JSON(http.StatusOK, PostMePlanResponse{
		Plan:         updatedPlan.Plan,
		Entitlements: entitlements,
	})
}
