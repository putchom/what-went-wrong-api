package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEntitlementChecker
type MockEntitlementChecker struct {
	mock.Mock
}

func (m *MockEntitlementChecker) CanUseAiExcuse(userID uuid.UUID) (bool, error) {
	args := m.Called(userID)
	return args.Bool(0), args.Error(1)
}

// MockAIService (re-defined here for simplicity or import if public)
type TestMockAIService struct {
	mock.Mock
}

func (m *TestMockAIService) GenerateExcuse(tone, context string) ([]string, error) {
	args := m.Called(tone, context)
	return args.Get(0).([]string), args.Error(1)
}

func TestPostAiExcuse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockEntitlement := new(MockEntitlementChecker)
		mockAI := new(TestMockAIService)
		handler := NewAIHandler(mockEntitlement, mockAI)

		userID := uuid.New()
		mockEntitlement.On("CanUseAiExcuse", userID).Return(true, nil)
		mockAI.On("GenerateExcuse", "surreal", "context").Return([]string{"excuse 1", "excuse 2"}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID.String())

		reqBody := CreateAiExcuseRequest{
			GoalID:  "g1",
			Date:    "2025-01-01",
			Tone:    "surreal",
			Context: "context",
		}
		jsonBytes, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest("POST", "/ai-excuse", bytes.NewBuffer(jsonBytes))

		handler.PostAiExcuse(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp CreateAiExcuseResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp.Candidates, 2)
		assert.Equal(t, "excuse 1", resp.Candidates[0])
	})

	t.Run("Forbidden_FreePlan", func(t *testing.T) {
		mockEntitlement := new(MockEntitlementChecker)
		mockAI := new(TestMockAIService)
		handler := NewAIHandler(mockEntitlement, mockAI)

		userID := uuid.New()
		mockEntitlement.On("CanUseAiExcuse", userID).Return(false, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID.String())

		reqBody := CreateAiExcuseRequest{
			GoalID: "g1",
			Date:   "2025-01-01",
		}
		jsonBytes, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest("POST", "/ai-excuse", bytes.NewBuffer(jsonBytes))

		handler.PostAiExcuse(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
