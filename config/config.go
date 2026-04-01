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
	Port              string
	JwtSecret         string
	AccessTokenExpiry int
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	expiry, _ := strconv.Atoi(getEnv("ACCESS_TOKEN_EXPIRY", "15"))
	return &Config{
		PostgresURL:       getEnv("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/locallyn"),
		RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
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
