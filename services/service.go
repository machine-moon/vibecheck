package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"strconv"
	"vibecheck/models"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type VibecheckService struct {
	db    *sql.DB
	redis *redis.Client
}

func NewVibecheckService(database *sql.DB, redisClient *redis.Client) *VibecheckService {
	return &VibecheckService{db: database, redis: redisClient}
}

// GetAllTweets retrieves all tweets from the database
func (s *VibecheckService) GetAllTweets() ([]models.Tweet, error) {
	query := "SELECT id, text, hint, answer FROM tweets"
	ctx := context.Background()

	// Try to get the cached result from Redis
	cachedTweets, err := s.redis.Get(ctx, query).Result()
	if err == nil {
		var tweets []models.Tweet
		if err := json.Unmarshal([]byte(cachedTweets), &tweets); err == nil {
			// Log cache hit
			log.Println("Cache hit for GetAllTweets")
			return tweets, nil
		}
	}

	// If cache miss or unmarshal error, query the database
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []models.Tweet
	for rows.Next() {
		var tweet models.Tweet
		if err := rows.Scan(&tweet.ID, &tweet.Text, &tweet.Hint, &tweet.Answer); err != nil {
			return nil, err
		}
		tweets = append(tweets, tweet)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Cache the result in Redis
	tweetsJSON, err := json.Marshal(tweets)
	if err == nil {
		s.redis.Set(ctx, query, tweetsJSON, 0)
	}

	return tweets, nil
}

// GetTweetsByPage retrieves a page of tweets from the database
func (s *VibecheckService) GetTweetsByPage(pageNumber int, listPerPage int) ([]models.Tweet, error) {
	if pageNumber < 1 {
		pageNumber = 1
	}

	offset := (pageNumber - 1) * listPerPage
	query := "SELECT id, text, hint, answer FROM tweets ORDER BY id LIMIT $1 OFFSET $2"
	ctx := context.Background()

	// Try to get the cached result from Redis
	cacheKey := "tweets_page_" + strconv.Itoa(pageNumber) + "_" + strconv.Itoa(listPerPage)
	cachedTweets, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var tweets []models.Tweet
		if err := json.Unmarshal([]byte(cachedTweets), &tweets); err == nil {
			return tweets, nil
		}
	}

	// If cache miss or unmarshal error, query the database
	rows, err := s.db.Query(query, listPerPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []models.Tweet
	for rows.Next() {
		var tweet models.Tweet
		if err := rows.Scan(&tweet.ID, &tweet.Text, &tweet.Hint, &tweet.Answer); err != nil {
			return nil, err
		}
		tweets = append(tweets, tweet)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Cache the result in Redis
	tweetsJSON, err := json.Marshal(tweets)
	if err == nil {
		s.redis.Set(ctx, cacheKey, tweetsJSON, 0)
	}

	return tweets, nil
}

// GetTweet retrieves a tweet by its ID from the database
func (s *VibecheckService) GetTweet(id string) (*models.Tweet, error) {
	query := "SELECT id, text, hint, answer FROM tweets WHERE id = $1"
	ctx := context.Background()

	// Try to get the cached result from Redis
	cacheKey := "tweet_" + id
	cachedTweet, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var tweet models.Tweet
		if err := json.Unmarshal([]byte(cachedTweet), &tweet); err == nil {
			return &tweet, nil
		}
	}

	// If cache miss or unmarshal error, query the database
	row := s.db.QueryRow(query, id)
	var tweet models.Tweet
	if err := row.Scan(&tweet.ID, &tweet.Text, &tweet.Hint, &tweet.Answer); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Cache the result in Redis
	tweetJSON, err := json.Marshal(tweet)
	if err == nil {
		s.redis.Set(ctx, cacheKey, tweetJSON, 0)
	}

	return &tweet, nil
}

// NewTweet creates a new tweet in the database and caches it in Redis
func (s *VibecheckService) NewTweet(tweet *models.NewTweet) error {
	query := "INSERT INTO tweets (id, text, hint, answer) VALUES ($1, $2, $3, $4)"
	id := generateNewID()
	_, err := s.db.Exec(query, id, tweet.Text, tweet.Hint, tweet.Answer)
	if err != nil {
		return err
	}

	// Cache the new tweet in Redis
	tweetJSON, err := json.Marshal(tweet)
	if err == nil {
		ctx := context.Background()
		cacheKey := "tweet_" + id
		s.redis.Set(ctx, cacheKey, tweetJSON, 0)
	}

	return nil
}

// UpdateTweet updates an existing tweet in the database and updates the cache in Redis
func (s *VibecheckService) UpdateTweet(tweet *models.Tweet) error {
	query := "UPDATE tweets SET text = $1, hint = $2, answer = $3 WHERE id = $4"
	_, err := s.db.Exec(query, tweet.Text, tweet.Hint, tweet.Answer, tweet.ID)
	if err != nil {
		return err
	}

	// Update the cache in Redis
	tweetJSON, err := json.Marshal(tweet)
	if err == nil {
		ctx := context.Background()
		cacheKey := "tweet_" + tweet.ID
		s.redis.Set(ctx, cacheKey, tweetJSON, 0)
	}

	return nil
}

// DeleteTweet deletes a tweet from the database and removes it from the cache in Redis
func (s *VibecheckService) DeleteTweet(id string) error {
	query := "DELETE FROM tweets WHERE id = $1"
	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	// Remove the tweet from the cache in Redis
	ctx := context.Background()
	cacheKey := "tweet_" + id
	s.redis.Del(ctx, cacheKey)

	return nil
}

// User routes

// GetAllProblems retrieves all problems from the database
func (s *VibecheckService) GetAllProblems() ([]models.Problem, error) {
	query := "SELECT id, text FROM tweets"
	ctx := context.Background()

	// Try to get the cached result from Redis
	cachedProblems, err := s.redis.Get(ctx, query).Result()
	if err == nil {
		var problems []models.Problem
		if err := json.Unmarshal([]byte(cachedProblems), &problems); err == nil {
			return problems, nil
		}
	}

	// If cache miss or unmarshal error, query the database
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var problems []models.Problem
	for rows.Next() {
		var problem models.Problem
		if err := rows.Scan(&problem.ID, &problem.Text); err != nil {
			return nil, err
		}
		problems = append(problems, problem)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Cache the result in Redis
	problemsJSON, err := json.Marshal(problems)
	if err == nil {
		s.redis.Set(ctx, query, problemsJSON, 0)
	}

	return problems, nil
}

// GetProblemsByPage retrieves a page of problems from the database
func (s *VibecheckService) GetProblemsByPage(pageNumber int, listPerPage int) ([]models.Problem, error) {
	if pageNumber < 1 {
		pageNumber = 1
	}

	offset := (pageNumber - 1) * listPerPage
	query := "SELECT id, text FROM tweets ORDER BY id LIMIT $1 OFFSET $2"
	ctx := context.Background()

	// Try to get the cached result from Redis
	cacheKey := "problems_page_" + strconv.Itoa(pageNumber) + "_" + strconv.Itoa(listPerPage)
	cachedProblems, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var problems []models.Problem
		if err := json.Unmarshal([]byte(cachedProblems), &problems); err == nil {
			return problems, nil
		}
	}

	// If cache miss or unmarshal error, query the database
	rows, err := s.db.Query(query, listPerPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var problems []models.Problem
	for rows.Next() {
		var problem models.Problem
		if err := rows.Scan(&problem.ID, &problem.Text); err != nil {
			return nil, err
		}
		problems = append(problems, problem)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Cache the result in Redis
	problemsJSON, err := json.Marshal(problems)
	if err == nil {
		s.redis.Set(ctx, cacheKey, problemsJSON, 0)
	}

	return problems, nil
}

// NewProblem creates a new problem in the database and caches it in Redis
func (s *VibecheckService) NewProblem(problem *models.NewProblem) error {
	query := "INSERT INTO tweets (id, text, hint, answer) VALUES ($1, $2, $3, $4)"
	id := generateNewID()
	_, err := s.db.Exec(query, id, problem.Text, problem.Hint, problem.Answer)
	if err != nil {
		return err
	}

	// Cache the new problem in Redis
	problemJSON, err := json.Marshal(problem)
	if err == nil {
		ctx := context.Background()
		cacheKey := "problem_" + id
		s.redis.Set(ctx, cacheKey, problemJSON, 0)
	}

	return nil
}

// GetProblem retrieves a tweet without hint and answer by its ID
func (s *VibecheckService) GetProblem(id string) (*models.Problem, error) {
	query := "SELECT id, text FROM tweets WHERE id = $1"
	ctx := context.Background()

	// Try to get the cached result from Redis
	cacheKey := "problem_" + id
	cachedProblem, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var problem models.Problem
		if err := json.Unmarshal([]byte(cachedProblem), &problem); err == nil {
			return &problem, nil
		}
	}

	// If cache miss or unmarshal error, query the database
	row := s.db.QueryRow(query, id)
	var problem models.Problem
	if err := row.Scan(&problem.ID, &problem.Text); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tweet not found")
		}
		return nil, err
	}

	// Cache the result in Redis
	problemJSON, err := json.Marshal(problem)
	if err == nil {
		s.redis.Set(ctx, cacheKey, problemJSON, 0)
	}

	return &problem, nil
}

// GetRandomProblem retrieves a random tweet without hint and answer
func (s *VibecheckService) GetRandomProblem() (*models.Problem, error) {
	// Get the total number of tweets
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM tweets").Scan(&count)
	if err != nil {
		return nil, err
	}

	// Check if there are any tweets
	if count == 0 {
		return nil, errors.New("no tweets available")
	}

	// Generate a random offset
	offset := rand.Intn(count)

	// Fetch the tweet at the random offset
	query := "SELECT id, text FROM tweets LIMIT 1 OFFSET $1"
	row := s.db.QueryRow(query, offset)

	var problem models.Problem
	if err := row.Scan(&problem.ID, &problem.Text); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no tweets available")
		}
		return nil, err
	}

	// Cache the result in Redis
	problemJSON, err := json.Marshal(problem)
	if err == nil {
		ctx := context.Background()
		cacheKey := "problem_random"
		s.redis.Set(ctx, cacheKey, problemJSON, 0)
	}

	return &problem, nil
}

// CheckSolution checks if the user's guess is correct
func (s *VibecheckService) CheckSolution(attempt *models.AttemptSolution) (bool, error) {
	tweet, err := s.GetTweet(attempt.ID)
	if err != nil {
		return false, err
	}
	if tweet == nil {
		return false, errors.New("tweet not found")
	}
	return tweet.Answer == attempt.Guess, nil
}

// GetHint retrieves the hint for a specific tweet
func (s *VibecheckService) GetHint(tweetID string) (string, error) {
	tweet, err := s.GetTweet(tweetID)
	if err != nil {
		return "", err
	}
	if tweet == nil {
		return "", errors.New("tweet not found")
	}
	return tweet.Hint, nil
}

func generateNewID() string {
	return uuid.New().String()
}
