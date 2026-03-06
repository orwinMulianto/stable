package auth

import (
	"github.com/gin-gonic/gin"

)

func AuthRouter (router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", RegisterHandler)
		auth.POST("/verify-email", VerifyEmailHandler)
	}
}