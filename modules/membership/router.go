package membership

import (
	"stable/database/migrations"

	"github.com/gin-gonic/gin"
)

func MembershipRouter(router *gin.RouterGroup) {
	membershipRepo := NewRepository(migrations.GetDB())
	membershipService := NewService(membershipRepo)
	membershipHandler := NewHandler(membershipService)

	membership := router.Group("/membership")
	{
		membership.POST("/trial", membershipHandler.ClaimTrialHandler)
	}
}