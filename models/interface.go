package models

type Tweet struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Hint   string `json:"hint"`
	Answer string `json:"answer"`
}

type NewTweet struct {
	Text   string `json:"text"`
	Hint   string `json:"hint"`
	Answer string `json:"answer"`
}

type NewProblem = NewTweet

type Problem struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type HintContent struct {
	ID   string `json:"id"`
	Hint string `json:"hint"`
}

type AnswerContent struct {
	ID     string `json:"id"`
	Answer string `json:"answer"`
}

type AttemptSolution struct {
	ID    string `json:"id"`
	Guess string `json:"guess"`
}
