package user

import (
	"github.com/gin-gonic/gin"

)

func UserRouter(router *gin.RouterGroup) {
	user := router.Group("/users")
	{
		user.GET("/:id", GetUserbyIDHandler)
	}
}