package profile

import (
	"stable/database/migrations"

	"github.com/gin-gonic/gin"
)

func ProfileRouter(router *gin.RouterGroup) {
	repository := NewRepository(migrations.GetDB())
	service := NewService(repository)
	handler := NewHandler(service)

	profile := router.Group("/profile")
	{
		profile.GET("/:user_id", handler.GetProfileHandler)
		profile.PATCH("/:user_id", handler.UpdateProfileHandler)
		profile.POST("/:user_id/avatar", handler.UploadAvatarHandler)
	}
}
