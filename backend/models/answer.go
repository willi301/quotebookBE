package models

type AnswerRequest struct {
	SessionID  string `json:"session_id"`
	QuestionID int    `json:"question_id"`
	Answer     string `json:"answer"`
}

