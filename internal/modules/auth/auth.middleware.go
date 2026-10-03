package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ContextUserId         = "userId"
	ContextOrganizationId = "organizationId"
)

type AuthMiddleware struct {
	secretKey []byte
}

func NewAuthMiddleware(secretKey string) *AuthMiddleware {
	return &AuthMiddleware{
		secretKey: []byte(secretKey),
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("accessToken")
		if err != nil || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Authentication required",
				"data":    nil,
			})
			return
		}

		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (any, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, jwt.ErrSignatureInvalid
				}

				return m.secretKey, nil
			},
		)

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid or expired token",
				"data":    nil,
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token claims",
				"data":    nil,
			})
			return
		}

		userId, ok := claims["userId"].(string)
		if !ok || userId == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid user ID in token",
				"data":    nil,
			})
			return
		}

		organizationId, ok := claims["organizationId"].(string)
		if !ok || organizationId == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid organization ID in token",
				"data":    nil,
			})
			return
		}

		c.Set(ContextUserId, userId)
		c.Set(ContextOrganizationId, organizationId)

		c.Next()
	}
}
