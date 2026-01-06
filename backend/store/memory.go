package store

import (
	"sync"

	"backend/models"
)

var (
	Questions []models.Question
	Sessions  = make(map[string]*models.QuizSession)
	Mu        sync.Mutex
)
