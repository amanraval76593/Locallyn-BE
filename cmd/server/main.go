package main

import (
	"locallyn-be/api"
	"locallyn-be/config"
	"locallyn-be/pkg/database"
	"locallyn-be/pkg/redis"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	database.InitPostgres(cfg.PostgresURL)

	redis.InitRedis(cfg.RedisAddr)

	router := gin.Default()

	api.SetupRoutes(router)

	log.Println("Server running on Port: ", cfg.Port)

	router.Run(":" + cfg.Port)

}
