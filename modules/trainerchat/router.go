package trainerchat

import (
	"stable/database/migrations"
	"stable/middleware"

	"github.com/gin-gonic/gin"
)

func TrainerChatRouter(router *gin.RouterGroup) {
	repository := NewRepository(migrations.GetDB())
	service := NewService(repository)
	handler := NewHandler(service)
	hub := GetGlobalHub()

	trainerChat := router.Group("/trainer-chat")
	{
		trainerChat.GET("/trainers", handler.ListTrainersHandler)
		trainerChat.POST("/checkout", handler.CheckoutHandler)
		trainerChat.POST("/notification", handler.NotificationHandler)
		trainerChat.GET("/sessions/:session_id", handler.SessionHandler)
		trainerChat.POST("/sessions/:session_id/messages", handler.SendMessageHandler)
		trainerChat.GET("/sessions/:session_id/ws", WSHandler(hub, service))
		trainerChat.POST("/sessions/:session_id/confirm", handler.ConfirmPaymentHandler)
		trainerChat.POST("/sessions/:session_id/dev-paid", handler.DevMarkPaidHandler)

		protected := trainerChat.Group("")
		protected.Use(middleware.RequireAuth())
		{
			protected.GET("/dashboard/me", handler.DashboardMeHandler)
			protected.GET("/history", handler.HistoryHandler)
			protected.GET("/trainer-sessions", handler.GetTrainerSessionsHandler)
		}
	}
}