package handler

import (
	"net/http"

	"github.com/NekKkMirror/medods-tz.git/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// IssueTokens handles the issuing of access and refresh tokens for a user.
func (h *AuthHandler) IssueTokens(c *gin.Context) {
	userID := c.Query("userID")
	clientIP := c.ClientIP()

	accessToken, refreshToken, err := h.authService.IssueTokens(userID, clientIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

// RefreshTokens handles the refreshing of access and refresh tokens for a user.
func (h *AuthHandler) RefreshTokens(c *gin.Context) {
	refreshToken := c.PostForm("refreshToken")
	clientIP := c.ClientIP()

	accessToken, newRefreshToken, err := h.authService.RefreshToken(refreshToken, clientIP)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": newRefreshToken,
	})
}
