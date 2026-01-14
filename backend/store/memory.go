package store

import (
	"sync"

	"backend/models"
)

var (
	Questions []models.Question
	Mu        sync.Mutex
)
