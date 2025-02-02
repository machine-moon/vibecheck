package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"vibecheck/models"
	"vibecheck/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type vibecheckController struct {
	vibecheckService services.VibecheckService
	listPerPage      int
}

// NewvibecheckController creates a new vibecheck controller
func NewVibecheckController(db *sql.DB, redisClient *redis.Client, lpp int) *vibecheckController {
	vibecheckService := services.NewVibecheckService(db, redisClient)
	return &vibecheckController{vibecheckService: *vibecheckService, listPerPage: lpp}
}

// GetTweets retrieves all tweets from the database
func (vc *vibecheckController) GetTweets(c *gin.Context) {
	tweets, err := vc.vibecheckService.GetAllTweets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tweets retrieved successfully", "tweets": tweets})
}

// GetTweetsByPage retrieves a page of tweets from the database
func (vc *vibecheckController) GetTweetsByPage(c *gin.Context) {
	pageNumber, err := strconv.Atoi(c.Param("pageNumber"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}
	tweets, err := vc.vibecheckService.GetTweetsByPage(pageNumber, vc.listPerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tweets retrieved successfully", "tweets": tweets})
}

// CreateTweet creates a new tweet in the database
func (vc *vibecheckController) CreateTweet(c *gin.Context) {
	var tweet models.CreateTweet
	if err := c.ShouldBindJSON(&tweet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := vc.vibecheckService.CreateTweet(&tweet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Tweet created successfully", "tweet": tweet})
}

// GetTweet retrieves a tweet by its ID
func (vc *vibecheckController) GetTweet(c *gin.Context) {
	id := c.Param("id")
	tweet, err := vc.vibecheckService.GetTweet(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tweet retrieved successfully", "tweet": tweet})
}

// UpdateTweet updates an existing tweet in the database
func (vc *vibecheckController) UpdateTweet(c *gin.Context) {
	id := c.Param("id")
	tweet := models.Tweet{ID: id}
	if err := c.ShouldBindJSON(&tweet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := vc.vibecheckService.UpdateTweet(&tweet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tweet updated successfully", "tweet": tweet})
}

// DeleteTweet deletes a tweet from the database
func (vc *vibecheckController) DeleteTweet(c *gin.Context) {
	id := c.Param("id")
	err := vc.vibecheckService.DeleteTweet(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tweet deleted successfully"})
}

// User routes

// GetProblems retrieves all problems from the database
func (vc *vibecheckController) GetProblems(c *gin.Context) {
	problems, err := vc.vibecheckService.GetAllProblems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Problems retrieved successfully", "problems": problems})
}

// GetProblemsByPage retrieves a page of problems from the database
func (vc *vibecheckController) GetProblemsByPage(c *gin.Context) {
	pageNumber, err := strconv.Atoi(c.Param("pageNumber"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}
	problems, err := vc.vibecheckService.GetProblemsByPage(pageNumber, vc.listPerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Problems retrieved successfully", "problems": problems})
}

// CreateProblem creates a new problem in the database
func (vc *vibecheckController) CreateProblem(c *gin.Context) {
	var problem models.CreateProblem
	if err := c.ShouldBindJSON(&problem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := vc.vibecheckService.CreateProblem(&problem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Problem created successfully", "problem": problem})
}

// GetProblem retrieves a tweet without hint and answer by its ID
func (vc *vibecheckController) GetProblem(c *gin.Context) {
	id := c.Param("id")
	problem, err := vc.vibecheckService.GetProblem(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Problem retrieved successfully", "problem": problem})
}

// GetProblem retrieves a tweet without hint and answer
func (vc *vibecheckController) GetRandomProblem(c *gin.Context) {
	problem, err := vc.vibecheckService.GetProblem()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Problem retrieved successfully", "problem": problem})
}

// AnswerProblem checks if the user's solution is correct
func (vc *vibecheckController) AnswerProblem(c *gin.Context) {
	var attempt models.AttemptSolution
	if err := c.ShouldBindJSON(&attempt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	correct, err := vc.vibecheckService.CheckSolution(&attempt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"correct": correct})
}

// GetHint retrieves a hint for a tweet
func (vc *vibecheckController) GetHint(c *gin.Context) {
	id := c.Param("tweetId")
	hint, err := vc.vibecheckService.GetHint(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"hint": hint})
}
