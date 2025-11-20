package upload

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	handler := NewHandler()

	upload := r.Group("/v1/upload")
	{
		upload.POST("/avatar", handler.Avatar)
	}
}
