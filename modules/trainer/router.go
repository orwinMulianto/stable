package trainer

import (
	"stable/database/migrations"
	"stable/middleware"

	"github.com/gin-gonic/gin"
)

func TrainerRouter(router *gin.RouterGroup) {
	trainerRepo := NewRepository(migrations.GetDB())
	trainerService := NewService(trainerRepo)
	trainerHandler := NewHandler(trainerService)

	trainer := router.Group("/trainers")
	{
		// Public routes
		trainer.GET("", trainerHandler.GetAllTrainersHandler)
		trainer.GET("/detail/:id", trainerHandler.GetTrainerByIDHandler)

		// Protected routes
		protected := trainer.Group("")
		protected.Use(middleware.RequireAuth())
		{
			protected.GET("/me", trainerHandler.GetMyProfileHandler)
			protected.POST("/profile", trainerHandler.CreateProfileHandler)
			protected.PUT("/profile", trainerHandler.UpdateProfileHandler)
			protected.DELETE("/profile", trainerHandler.DeleteProfileHandler)
		}
	}
}