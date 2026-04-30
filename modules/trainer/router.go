package trainer

import (
	"github.com/gin-gonic/gin"
	"stable/database/migrations"
	"stable/middleware"
)

func TrainerRouter(router *gin.RouterGroup) {
	trainerRepo    := NewRepository(migrations.GetDB())
	trainerService := NewService(trainerRepo)
	trainerHandler := NewHandler(trainerService)

	trainer := router.Group("/trainers")
	{
		// Public - tidak perlu auth
		trainer.GET("", trainerHandler.GetAllTrainersHandler)           // GET /api/v1/trainers
		trainer.GET("/:id", trainerHandler.GetTrainerByIDHandler)       // GET /api/v1/trainers/:id

		// Protected - perlu auth
		protected := trainer.Group("")
		protected.Use(middleware.RequireAuth())
		{
			protected.GET("/me", trainerHandler.GetMyProfileHandler)        // GET /api/v1/trainers/me
			protected.POST("/profile", trainerHandler.CreateProfileHandler) // POST /api/v1/trainers/profile
			protected.PUT("/profile", trainerHandler.UpdateProfileHandler)  // PUT /api/v1/trainers/profile
			protected.DELETE("/profile", trainerHandler.DeleteProfileHandler) // DELETE /api/v1/trainers/profile
		}
	}
}