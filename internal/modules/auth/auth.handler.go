package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kapilkumar9395/pypr/internal/common"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")

	auth.POST("/register", h.RegisterUser)
	auth.POST("/login", h.Login)
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	// Placeholder for register handler logic
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.service.RegisterUser(
		req.EmailId,
		req.FirstName,
		req.LastName,
		req.Password,
		nil, // organizationId
		req.InvitationToken,
		req.Timezone,
	)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "User registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	loginResponse, err := h.service.UserLogin(
		req.EmailId,
		req.Password,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "accessToken",
		Value:    loginResponse.AccessToken,
		HttpOnly: true,
		Secure:   false, // true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	c.JSON(http.StatusOK, common.GenericResponse{
		Message: "User logged in successfully",
		Data:    nil,
	})
}
