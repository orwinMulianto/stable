package users

import (
	"github.com/gin-gonic/gin"
	"stable/database/migrations"

)

func UserRouter(router *gin.RouterGroup) {
	userRepo := NewRepository(migrations.GetDB())
	userService := NewService(userRepo)
	userHandler := NewHandler(userService)

	user := router.Group("/users")
	{
		user.GET("/:id", userHandler.GetUserByIDHandler)}
}