package api

import (
	"locallyn-be/config"
	commonauth "locallyn-be/internal/common/auth"
	"locallyn-be/internal/feed"

	"github.com/gin-gonic/gin"
)

func SetUpFeedRoute(router *gin.Engine) {
	cfg := config.LoadConfig()
	feedRepo := feed.NewRepository()
	feedService := feed.NewService(feedRepo)
	feedHandler := feed.NewHandler(feedService)

	feedRouter := router.Group("/feed")
	{
		feedRouter.GET("/get-feed-by-location", commonauth.RequireAuth(cfg.JwtSecret), feedHandler.GetFeedByLocationHandler)
		feedRouter.GET("/get-incident-posts/:incidentId", commonauth.RequireAuth(cfg.JwtSecret), feedHandler.GetIncidentPostsHandler)
	}
}
