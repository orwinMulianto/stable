package auth

import (
	"stable/database/migrations"
	"stable/modules/users"

	"github.com/gin-gonic/gin"
)

func AuthRouter(router *gin.RouterGroup) {
	userRepo    := users.NewRepository(migrations.GetDB())
	authService := NewService(userRepo)
	authHandler := NewHandler(authService)

	auth := router.Group("/auth")
	{
		auth.POST("/register",             authHandler.RegisterHandler)
		auth.POST("/login",                authHandler.LoginHandler)    
		auth.POST("/verify-email",         authHandler.VerifyEmailHandler)
		auth.POST("/resend-verification",  authHandler.ResendVerificationCodeHandler)
		auth.GET("/google/login", authHandler.GoogleLoginHandler)
		auth.GET("/google/callback", authHandler.GoogleCallbackHandler)
	}
}