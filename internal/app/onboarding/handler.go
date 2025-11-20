package onboarding

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	// Add dependencies here later
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Basic(c *gin.Context) {
	var req BasicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, BasicResponse{
		TempID:  "temp-uuid-67890",
		Message: "basic_saved",
	})
}

func (h *Handler) Profile(c *gin.Context) {
	var req ProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ProfileResponse{
		TempID:            req.TempID,
		UsernameAvailable: true,
		AvatarURL:         "https://cdn.example.com/avatars/dummy.jpg",
	})
}

func (h *Handler) Interests(c *gin.Context) {
	var req InterestsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, InterestsResponse{
		Message: "interests_saved",
	})
}
