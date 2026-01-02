package models


type QuizSession struct {
	ID        string
	Questions []Question
	Score     int
}