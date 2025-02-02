package config

import (
	"os"
	"strconv"
)

type Config struct {
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Database string
	}
	Redis struct {
		Host     string
		Port     string
		Password string
		DB       int
	}
	ServicePort string
	ListPerPage int
}

func LoadConfig() Config {
	listPerPage, err := strconv.Atoi(getEnv("LIST_PER_PAGE", "10"))
	if err != nil {
		listPerPage = 10
	}
	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		redisDB = 0
	}

	config := Config{}
	config.DB.Host = getEnv("DB_HOST", "certainlyNotLocalhost")
	config.DB.Port = getEnv("DB_PORT", "6969")
	config.DB.User = getEnv("DB_USER", "notTarek")
	config.DB.Password = getEnv("DB_PASSWORD", "notMohammed")
	config.DB.Database = getEnv("DB_NAME", "whyareyoureadingthis")
	config.Redis.Host = getEnv("REDIS_HOST", "goodquestion")
	config.Redis.Port = getEnv("REDIS_PORT", "6379")
	config.Redis.Password = getEnv("REDIS_PASSWORD", "1H@t3R3dis")
	config.Redis.DB = redisDB
	config.ListPerPage = listPerPage
	config.ServicePort = getEnv("SERVICE_PORT", "8080")
	return config
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
