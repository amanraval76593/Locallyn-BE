package api

import (
	"locallyn-be/config"
	commonauth "locallyn-be/internal/common/auth"
	"locallyn-be/internal/interactions"

	"github.com/gin-gonic/gin"
)

func SetUpInteractionRoutes(router *gin.Engine) {
	cfg := config.LoadConfig()
	interactionsRepo := interactions.NewRepository()
	interactionsService := interactions.NewService(interactionsRepo)
	interactionsHandler := interactions.NewHandler(interactionsService)
	interactionRoute := router.Group("/interactions")
	{
		interactionRoute.POST("/confirm-incident/:incidentId", commonauth.RequireAuth(cfg.JwtSecret), interactionsHandler.ConfirmIncidentHandler)
		interactionRoute.POST("/feedback/:postId", commonauth.RequireAuth(cfg.JwtSecret), interactionsHandler.PostFeedbackHandler)
		interactionRoute.POST("/report-post", commonauth.RequireAuth(cfg.JwtSecret), interactionsHandler.ReportPost)
	}
}
