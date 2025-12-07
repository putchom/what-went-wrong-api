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
	"github.com/stretchr/testify/mock"
)

type MockEntitlementManager struct {
	mock.Mock
}

func (m *MockEntitlementManager) GetPlan(userID uuid.UUID) (*models.UserPlan, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserPlan), args.Error(1)
}

func (m *MockEntitlementManager) GetEntitlements(planName string) services.Entitlements {
	args := m.Called(planName)
	return args.Get(0).(services.Entitlements)
}

func (m *MockEntitlementManager) UpdatePlan(userID uuid.UUID, planName string) (*models.UserPlan, error) {
	args := m.Called(userID, planName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserPlan), args.Error(1)
}

func TestGetMePlan(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockManager := new(MockEntitlementManager)
		handler := NewPlanHandler(mockManager)

		userID := uuid.New()
		mockPlan := &models.UserPlan{UserID: userID, Plan: "free"}
		mockEntitlements := services.Entitlements{MaxGoals: 3}

		mockManager.On("GetPlan", userID).Return(mockPlan, nil)
		mockManager.On("GetEntitlements", "free").Return(mockEntitlements)

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
	})
}

func TestPostMePlan(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockManager := new(MockEntitlementManager)
		handler := NewPlanHandler(mockManager)

		userID := uuid.New()
		updatedPlan := &models.UserPlan{
			UserID:    userID,
			Plan:      "premium",
			UpdatedAt: time.Now(),
		}
		premiumEntitlements := services.Entitlements{MaxGoals: 100}

		mockManager.On("UpdatePlan", userID, "premium").Return(updatedPlan, nil)
		mockManager.On("GetEntitlements", "premium").Return(premiumEntitlements)

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
	})
}
