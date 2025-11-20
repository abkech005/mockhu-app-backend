package onboarding

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	handler := NewHandler()

	onboard := r.Group("/v1/onboard")
	{
		onboard.POST("/basic", handler.Basic)
		onboard.POST("/profile", handler.Profile)
		onboard.POST("/interests", handler.Interests)
	}
}
