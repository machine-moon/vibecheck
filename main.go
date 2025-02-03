package main

import (
	"context"
	"database/sql"
	"log"
	"vibecheck/config"
	"vibecheck/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.LoadConfig()
	log.Printf("Loaded config: %+v\n", cfg)

	connStr := "host=" + cfg.DB.Host + " port=" + cfg.DB.Port + " user=" + cfg.DB.User + " password=" + cfg.DB.Password + " dbname=" + cfg.DB.Database + " sslmode=disable"
	log.Printf("Connecting to PostgreSQL with connection string: %s\n", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AddAllowHeaders("Content-Type")
	corsConfig.AddAllowHeaders("X-Requested-With")
	corsConfig.AddAllowHeaders("Accept")
	corsConfig.AddAllowHeaders("Origin")
	corsConfig.AddAllowHeaders("X-CSRF-Token")

	r.Use(cors.New(corsConfig))

	routes.SetupRoutes(r, db, redisClient, cfg.ListPerPage)

	r.Run(":" + cfg.ServicePort)
}
