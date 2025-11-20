package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	// Add dependencies here later
	// service *Service
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, SignupResponse{
		UserID:              "user-uuid-12345",
		VerificationNeeded:  true,
		VerificationChannel: req.Method,
	})
}

func (h *Handler) Verify(c *gin.Context) {
	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, VerifyResponse{
		AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.dummy",
		RefreshToken: "refresh-token-dummy",
		ExpiresIn:    900,
		User: &UserInfo{
			ID:       req.UserID,
			Username: "john123",
			Email:    "user@example.com",
		},
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.dummy",
		RefreshToken: "refresh-token-dummy",
		ExpiresIn:    900,
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, RefreshResponse{
		AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.new-dummy",
		RefreshToken: "new-refresh-token-dummy",
		ExpiresIn:    900,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, LogoutResponse{
		Message: "logged_out",
	})
}

func (h *Handler) Resend(c *gin.Context) {
	var req ResendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ResendResponse{
		Message: "code_sent",
	})
}
