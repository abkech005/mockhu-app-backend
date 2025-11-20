package upload

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	// Add storage service later
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Avatar(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	_ = file // Use file later

	c.JSON(http.StatusOK, AvatarResponse{
		AvatarURL: "https://cdn.example.com/avatars/uploaded.jpg",
	})
}
