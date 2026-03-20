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
		auth.POST("/login",                authHandler.LoginHandler)       // ← tambahan
		auth.POST("/verify-email",         VerifyEmailHandler)
		auth.POST("/resend-verification",  ResendVerificationCodeHandler)
	}
}