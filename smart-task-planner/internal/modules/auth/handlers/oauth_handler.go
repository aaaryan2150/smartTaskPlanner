package handlers

import (
	"net/http"
	"smart-task-planner/internal/modules/auth/service"

	"github.com/gin-gonic/gin"
)

// GoogleLogin redirects user to Google OAuth consent page
func GoogleLogin(c *gin.Context) {
	state := "randomstate" // In production, generate securely
	url := service.GetGoogleLoginURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback handles Google OAuth callback
func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	token, email, err := service.HandleGoogleCallback(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"email": email,
	})
}
