package store

import (
	"sync"

	"backend/models"
)

var (
	Sessions = make(map[string]*models.QuizSession)
	Mu       sync.Mutex
)