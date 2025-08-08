package handlers

import (
	"api/internal/middleware"
	"api/internal/services"
	"api/internal/utils"
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.GET("/login", h.Login)
		auth.GET("/callback", h.Callback)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/logout", h.Logout)
		auth.GET("/me", middleware.AuthMiddleware(nil), h.GetProfile)
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	state, err := generateState()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate state")
		return
	}

	authURL := h.authService.GetAuthURL(state)

	utils.SuccessResponse(c, "Authentication URL generated", gin.H{
		"auth_url": authURL,
		"state":    state,
	})
}

func (h *AuthHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	
	if code == "" {
		utils.BadRequestResponse(c, "Authorization code is required")
		return
	}

	if state == "" {
		utils.BadRequestResponse(c, "State parameter is required")
		return
	}

	loginResponse, err := h.authService.HandleCallback(code)
	if err != nil {
		utils.BadRequestResponse(c, "Authentication failed: "+err.Error())
		return
	}

	utils.SuccessResponse(c, "Login successful", loginResponse)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var request struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	newToken, err := h.authService.RefreshToken(request.Token)
	if err != nil {
		utils.UnauthorizedResponse(c, "Invalid or expired token")
		return
	}

	utils.SuccessResponse(c, "Token refreshed successfully", gin.H{
		"token": newToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	utils.SuccessResponse(c, "Logged out successfully", nil)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	userEmail := middleware.GetUserEmailFromContext(c)
	userRole := middleware.GetUserRoleFromContext(c)

	utils.SuccessResponse(c, "Profile retrieved successfully", gin.H{
		"id":    userID,
		"email": userEmail,
		"role":  userRole,
	})
}

func generateState() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}