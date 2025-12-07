package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func NewAuthMiddleware() (gin.HandlerFunc, error) {
	domain := os.Getenv("AUTH0_DOMAIN")
	audience := os.Getenv("AUTH0_AUDIENCE")
	if domain == "" || audience == "" {
		return nil, errors.New("AUTH0_DOMAIN or AUTH0_AUDIENCE is missing")
	}

	jwksURL := fmt.Sprintf("https://%s/.well-known/jwks.json", domain)
	issuer := fmt.Sprintf("https://%s/", domain)

	jwks, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("failed to create JWKS from resource at the given URL: %w", err)
	}

	return AuthMiddleware(jwks.Keyfunc, audience, issuer), nil
}

func AuthMiddleware(kf jwt.Keyfunc, audience, issuer string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, kf,
			jwt.WithAudience(audience),
			jwt.WithIssuer(issuer),
		)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token", "detail": err.Error()})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			sub, _ := claims["sub"].(string)
			c.Set("userID", sub)
			c.Set("claims", claims)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}
	}
}
