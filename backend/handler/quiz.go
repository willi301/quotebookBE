package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
	"backend/models"
	"backend/store"
)

// StartQuiz initializes a new quiz session
func StartQuiz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := generateSessionID()

	// Copy questions from store (DO NOT mutate)
	shuffled := make([]models.Question, len(store.Questions))
	copy(shuffled, store.Questions)

	// Shuffle
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	// Create session
	session := &models.QuizSession{
		ID:        sessionID,
		Questions: shuffled,
		Score:     0,
	}

	// Save session
	store.Mu.Lock()
	store.Sessions[sessionID] = session
	store.Mu.Unlock()

	// Create public questions
	public := make([]models.PublicQuestion, len(shuffled))
	for i, q := range shuffled {
		public[i] = models.PublicQuestion{
			ID:   q.ID,
			Text: q.Text,
		}
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"session_id": sessionID,
		"questions":  public,
	})
}


// to randomly generate a session ID
func generateSessionID() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}