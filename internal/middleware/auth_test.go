package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Generate a test RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}
	publicKey := &privateKey.PublicKey

	// Mock Keyfunc that returns the public key
	mockKeyFunc := func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	}

	audience := "test-audience"
	issuer := "https://test-domain/"

	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Missing Authorization Header",
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "missing or invalid Authorization header",
		},
		{
			name: "Invalid Authorization Header Format",
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				req.Header.Set("Authorization", "InvalidFormat")
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "missing or invalid Authorization header",
		},
		{
			name: "Invalid Token",
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				req.Header.Set("Authorization", "Bearer invalid.token.string")
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "invalid token",
		},
		{
			name: "Valid Token",
			setupRequest: func() *http.Request {
				// Create a valid JWT
				token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
					"sub": "test-user-id",
					"aud": audience,
					"iss": issuer,
					"exp": time.Now().Add(time.Hour).Unix(),
				})
				tokenString, err := token.SignedString(privateKey)
				if err != nil {
					t.Fatalf("Failed to sign token: %v", err)
				}

				req, _ := http.NewRequest("GET", "/", nil)
				req.Header.Set("Authorization", "Bearer "+tokenString)
				return req
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "test-user-id",
		},
		{
			name: "Expired Token",
			setupRequest: func() *http.Request {
				// Create an expired JWT
				token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
					"sub": "test-user-id",
					"aud": audience,
					"iss": issuer,
					"exp": time.Now().Add(-time.Hour).Unix(), // Expired
				})
				tokenString, err := token.SignedString(privateKey)
				if err != nil {
					t.Fatalf("Failed to sign token: %v", err)
				}

				req, _ := http.NewRequest("GET", "/", nil)
				req.Header.Set("Authorization", "Bearer "+tokenString)
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "token has invalid claims: token is expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = tt.setupRequest()

			// Initialize middleware with mock dependencies
			middleware := AuthMiddleware(mockKeyFunc, audience, issuer)

			// Dummy handler to verify success
			handler := func(c *gin.Context) {
				middleware(c)
				if !c.IsAborted() {
					userID, _ := c.Get("userID")
					c.String(http.StatusOK, userID.(string))
				}
			}

			handler(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}
		})
	}
}
