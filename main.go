package main

import (
	"context"
	"database/sql"
	"log"
	"vibecheck/config"
	"vibecheck/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.LoadConfig()

	// Connect to PostgreSQL
	connStr := "host=" + cfg.DB.Host + " port=" + cfg.DB.Port + " user=" + cfg.DB.User + " password=" + cfg.DB.Password + " dbname=" + cfg.DB.Database + " sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	// Set up the Gin router
	r := gin.Default()

	// Set up routes
	routes.SetupRoutes(r, db, redisClient, cfg.ListPerPage)

	// Start the server
	r.Run(":" + cfg.ServicePort)
}
