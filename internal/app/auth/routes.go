package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	handler := NewHandler()

	auth := r.Group("/v1/auth")
	{
		auth.POST("/signup", handler.Signup)
		auth.POST("/verify", handler.Verify)
		auth.POST("/login", handler.Login)
		auth.POST("/refresh", handler.Refresh)
		auth.POST("/logout", handler.Logout)
		auth.POST("/resend", handler.Resend)
	}
}
