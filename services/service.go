package services

import (
	"database/sql"
	"errors"
	"math/rand"
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
	return tweets, nil
}

// GetTweetsByPage retrieves a page of tweets from the database
func (s *VibecheckService) GetTweetsByPage(pageNumber int, listPerPage int) ([]models.Tweet, error) {
	if pageNumber < 1 {
		pageNumber = 1
	}

	offset := (pageNumber - 1) * listPerPage
	query := "SELECT id, text, hint, answer FROM tweets ORDER BY id LIMIT $1 OFFSET $2"
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
	return tweets, nil
}

// GetTweet retrieves a tweet by its ID from the database
func (s *VibecheckService) GetTweet(id string) (*models.Tweet, error) {
	query := "SELECT id, text, hint, answer FROM tweets WHERE id = $1"
	row := s.db.QueryRow(query, id)

	var tweet models.Tweet
	if err := row.Scan(&tweet.ID, &tweet.Text, &tweet.Hint, &tweet.Answer); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &tweet, nil
}

// CreateTweet creates a new tweet in the database
func (s *VibecheckService) CreateTweet(tweet *models.CreateTweet) error {
	query := "INSERT INTO tweets (id, text, hint, answer) VALUES ($1, $2, $3, $4)"
	_, err := s.db.Exec(query, generateNewID(), tweet.Text, tweet.Hint, tweet.Answer)
	if err != nil {
		return err
	}
	return nil
}

// UpdateTweet updates an existing tweet in the database
func (s *VibecheckService) UpdateTweet(tweet *models.Tweet) error {
	query := "UPDATE tweets SET text = $1, hint = $2, answer = $3 WHERE id = $4"
	_, err := s.db.Exec(query, tweet.Text, tweet.Hint, tweet.Answer, tweet.ID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTweet deletes a tweet from the database
func (s *VibecheckService) DeleteTweet(id string) error {
	query := "DELETE FROM tweets WHERE id = $1"
	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

// User routes

// GetAllProblems retrieves all problems from the database
func (s *VibecheckService) GetAllProblems() ([]models.Problem, error) {
	query := "SELECT id, text FROM tweets"
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
	return problems, nil
}

// GetProblemsByPage retrieves a page of problems from the database
func (s *VibecheckService) GetProblemsByPage(pageNumber int, listPerPage int) ([]models.Problem, error) {
	if pageNumber < 1 {
		pageNumber = 1
	}

	offset := (pageNumber - 1) * listPerPage
	query := "SELECT id, text FROM tweets ORDER BY id LIMIT $1 OFFSET $2"
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
	return problems, nil
}

// CreateProblem creates a new problem in the database
func (s *VibecheckService) CreateProblem(problem *models.CreateTweet) error {
	query := "INSERT INTO tweets (id, text, hint, answer) VALUES ($1, $2, $3, $4)"
	_, err := s.db.Exec(query, generateNewID(), problem.Text, problem.Hint, problem.Answer)
	if err != nil {
		return err
	}
	return nil
}

// GetProblem retrieves a tweet without hint and answer by its ID
func (s *VibecheckService) GetProblem(id string) (*models.Problem, error) {
	query := "SELECT id, text FROM tweets WHERE id = $1"
	row := s.db.QueryRow(query, id)

	var tweet models.Problem
	if err := row.Scan(&tweet.ID, &tweet.Text); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tweet not found")
		}
		return nil, err
	}
	return &tweet, nil
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

	var tweet models.Problem
	if err := row.Scan(&tweet.ID, &tweet.Text); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no tweets available")
		}
		return nil, err
	}
	return &tweet, nil
}

// CheckSolution checks if the user's guess is correct
func (s *VibecheckService) CheckSolution(attempt *models.AttemptSolution) (bool, error) {
	tweet, err := s.GetTweetByID(attempt.ID)
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
	tweet, err := s.GetTweetByID(tweetID)
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
