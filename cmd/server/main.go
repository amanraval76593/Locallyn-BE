package main

import (
	"locallyn-be/api"
	"locallyn-be/config"
	"locallyn-be/pkg/database"
	"locallyn-be/pkg/elasticsearch"
	"locallyn-be/pkg/redisConfig"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	database.InitPostgres(cfg.PostgresURL)

	redisConfig.InitRedis(cfg.RedisAddr)

	elasticsearch.InitElasticsearch(cfg.ElasticsearchURL)
	elasticsearch.BootstrapFeedPostsIndex(cfg.FeedSearchIndex)

	router := gin.Default()

	api.SetupRoutes(router)

	log.Println("Server running on Port: ", cfg.Port)

	router.Run(":" + cfg.Port)

}
