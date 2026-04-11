package api

import (
	"locallyn-be/config"
	"locallyn-be/internal/cache"
	commonauth "locallyn-be/internal/common/auth"
	"locallyn-be/internal/user"

	"github.com/gin-gonic/gin"
)

func setUpAuthRoutes(router *gin.Engine) {
	cfg := config.LoadConfig()
	userRepo := user.NewRepository()
	redis := cache.NewRedis()
	cacheRepo := user.NewUserCache(redis)
	userService := user.NewService(userRepo, *cacheRepo, cfg.JwtSecret, cfg.AccessTokenExpiry)
	userHandler := user.NewHandler(userService)

	auth := router.Group("/auth")
	{
		auth.POST("/signup", userHandler.SignUp)
		auth.POST("/verify-user", userHandler.VerifyUser)
		auth.POST("/login", userHandler.Login)
		auth.POST("/create-profile", commonauth.RequireAuth(cfg.JwtSecret), userHandler.CreateProfile)
		auth.PATCH("/update-profile", commonauth.RequireAuth(cfg.JwtSecret), userHandler.UpdateProfile)
		auth.GET("/get-profile", commonauth.RequireAuth(cfg.JwtSecret), userHandler.GetProfile)
	}
}
