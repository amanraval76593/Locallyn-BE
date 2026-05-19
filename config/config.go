package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresURL       string
	RedisAddr         string
	ElasticsearchURL  string
	FeedSearchIndex   string
	Port              string
	JwtSecret         string
	AccessTokenExpiry int
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	expiry, _ := strconv.Atoi(getEnv("ACCESS_TOKEN_EXPIRY", "60"))
	return &Config{
		PostgresURL:       getEnv("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/locallyn"),
		RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
		ElasticsearchURL:  getEnv("ELASTICSEARCH_URL", "http://localhost:9200"),
		FeedSearchIndex:   getEnv("FEED_SEARCH_INDEX", "feed_posts_v1"),
		Port:              getEnv("PORT", "8080"),
		JwtSecret:         getEnv("JWT_SECRET", ""),
		AccessTokenExpiry: expiry,
	}

}

func getEnv(key string, fallback any) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	switch v := fallback.(type) {
	case string:
		return v
	case int:
		return string(rune(v))
	default:
		return ""
	}
}
