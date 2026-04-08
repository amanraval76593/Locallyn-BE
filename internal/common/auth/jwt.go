package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userId, email, secret string, expiryMinutes int) (string, error) {
	claims := &Claims{
		UserId: userId,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ParseAccessToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid token signing method")
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

type contextKey string

const claimsContextKey contextKey = "auth_claims"

func RequireAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "authorization header is required",
			})
			return
		}

		tokenString, found := strings.CutPrefix(authHeader, "Bearer ")
		if !found || strings.TrimSpace(tokenString) == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "invalid authorization header",
			})
			return
		}

		claims, err := ParseAccessToken(strings.TrimSpace(tokenString), secret)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		ctx := context.WithValue(c.Request.Context(), claimsContextKey, claims)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func GetClaimsFromContext(ctx context.Context) (*Claims, error) {
	claims, ok := ctx.Value(claimsContextKey).(*Claims)
	if !ok || claims == nil {
		return nil, errors.New("authentication claims not found")
	}

	return claims, nil
}
