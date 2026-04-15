package api

import (
	"locallyn-be/config"
	commonauth "locallyn-be/internal/common/auth"
	"locallyn-be/internal/incident"
	"locallyn-be/internal/post"

	"github.com/gin-gonic/gin"
)

func setUpPostRoutes(router *gin.Engine) {
	cfg := config.LoadConfig()
	postRepo := post.NewRepository()
	incidentRepo := incident.NewRepository()
	incidentService := incident.NewService(incidentRepo)
	postService := post.NewService(postRepo, incidentService)
	postHandler := post.NewHandler(postService)

	auth := router.Group("/post")
	{
		auth.POST("/create-post", commonauth.RequireAuth(cfg.JwtSecret), postHandler.CreatePostHandler)
	}

}
