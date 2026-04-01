package api

import (
	"locallyn-be/config"
	"locallyn-be/internal/user"

	"github.com/gin-gonic/gin"
)

func setUpAuthRoutes(router *gin.Engine) {
	cfg := config.LoadConfig()
	userRepo := user.NewRepository()
	userService := user.NewService(userRepo, cfg.JwtSecret, cfg.AccessTokenExpiry)
	userHandler := user.NewHandler(userService)

	auth := router.Group("/auth")
	{
		auth.POST("/signup", userHandler.SignUp)
		auth.POST("/login", userHandler.Login)
	}
}
