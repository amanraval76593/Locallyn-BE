package api

import (
	"locallyn-be/config"
	commonauth "locallyn-be/internal/common/auth"
	"locallyn-be/internal/incident"

	"github.com/gin-gonic/gin"
)

func setUpIncidentRoute(router *gin.Engine) {

	cfg := config.LoadConfig()
	incidentRepo := incident.NewRepository()
	incidentService := incident.NewService(incidentRepo)
	incidentHandler := incident.NewHandler(incidentService)
	incident := router.Group("/incident")
	{
		incident.GET("/get-nearby-incident", commonauth.RequireAuth(cfg.JwtSecret), incidentHandler.GetNearbyIncidentHandler)
		incident.GET("/get-incident", commonauth.RequireAuth(cfg.JwtSecret), incidentHandler.GetIncidentHandler)
	}

}
