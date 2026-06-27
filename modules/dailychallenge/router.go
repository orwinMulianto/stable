package dailychallenge

import (
	"stable/database/migrations"

	"github.com/gin-gonic/gin"
)

func DailyChallengeRouter(router *gin.RouterGroup) {
	repository := NewRepository(migrations.GetDB())
	service := NewService(repository)
	handler := NewHandler(service)

	dailyChallenge := router.Group("/daily-challenge")
	{
		dailyChallenge.GET("/today", handler.GetTodayHandler)
		dailyChallenge.POST("/complete", handler.CompleteTodayHandler)
	}
}
