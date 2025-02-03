# Vibecheck

## Overview
Vibecheck is a web application that allows users to create, retrieve, update, and delete tweets. It also features a gameplay mode where users can solve problems based on tweets and receive hints.

## Project Structure
- `controllers/`: Handles HTTP request logic.
- `models/`: Defines the application's data models.
- `services/`: Contains business logic and database/Redis interactions.
- `config/`: Manages configuration loading.
- `routes/`: Defines application routes.
- `db/`: Includes database initialization scripts and data files.
- `docker/`: Contains Docker Compose files for setting up database and Redis services.
- `Dockerfile`: Builds the application container.
- `Makefile`: Automates building and running the application.
- `main.go`: Entry point of the application.
- `go.mod` and `go.sum`: Manage dependencies.

## Setup Instructions
1. Clone the repository:
    ```sh
    git clone https://github.com/machine-moon/vibecheck.git
    cd vibecheck
    ```

2. Set up the environment variables:
    ```sh
    source secrets.yml
    ```

3. Build and run the application using Docker Compose:
    ```sh
    docker-compose up --build
    ```

4. Alternatively, build and run the application locally:
    ```sh
    make run
    ```

## Usage
- Access the application at `http://localhost:8080`.
- Use the following endpoints to interact with the application:
  - `GET /tweets`: Retrieve all tweets.
  - `GET /tweets/page/:pageNumber`: Retrieve a page of tweets.
  - `POST /tweets/create`: Create a new tweet.
  - `PUT /tweets/:id`: Update an existing tweet.
  - `GET /tweets/:id`: Retrieve a tweet by its ID.
  - `DELETE /tweets/:id`: Delete a tweet.
  - `GET /problems`: Retrieve all problems.
  - `GET /problems/page/:pageNumber`: Retrieve a page of problems.
  - `POST /problems/create`: Create a new problem.
  - `GET /problem/:id`: Retrieve a problem by its ID.
  - `GET /problem/quiz`: Retrieve a random problem.
  - `POST /problem/answer`: Check if the user's solution is correct.
  - `GET /problem/hint/:tweetId`: Retrieve a hint for a problem.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
