package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "ok",
		})
	})
	setUpAuthRoutes(router)
	setUpIncidentRoute(router)
	setUpPostRoutes(router)
	SetUpFeedRoute(router)
	SetUpInteractionRoutes(router)
}
