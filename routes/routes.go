package routes

import (
	"database/sql"
	"vibecheck/controllers"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupRoutes(router *gin.Engine, db *sql.DB, redisClient *redis.Client, listPerPage int) {
	vibecheckController := controllers.NewVibecheckController(db, redisClient, listPerPage)

	// Dev routes
	router.GET("/tweets", vibecheckController.GetTweets) // For testing purposes
	router.GET("/tweets/page/:pageNumber", vibecheckController.GetTweetsByPage)

	router.POST("/tweets/create", vibecheckController.NewTweet)
	router.PUT("/tweets/:id", vibecheckController.UpdateTweet)
	router.GET("/tweets/:id", vibecheckController.GetTweet)
	router.DELETE("/tweets/:id", vibecheckController.DeleteTweet)

	// User routes

	router.GET("/problems", vibecheckController.GetProblems) // For testing purposes
	router.GET("/problems/page/:pageNumber", vibecheckController.GetProblemsByPage)
	router.POST("/problems/create", vibecheckController.NewProblem)

	// Gameplay routes
	router.GET("/problem/:id", vibecheckController.GetProblem)
	router.GET("/problem/quiz", vibecheckController.GetRandomProblem)
	router.POST("/problem/answer", vibecheckController.AnswerProblem)
	router.GET("/problem/hint/:tweetId", vibecheckController.GetHint)
}
