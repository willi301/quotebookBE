package handler

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
	"backend/models"
	"backend/store"
)

// StartQuiz initializes a new quiz session
func StartQuiz(w http.ResponseWriter, r *http.Request) {
	// making sure the request method will be post
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Generate session ID
	sessionID := generateSessionID()
	
	// Create a copy of questions and randomize them
	shuffledQuestions := make([]models.Question, len(questions))
	copy(shuffledQuestions, questions)
	
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(shuffledQuestions), func(i, j int) {
		shuffledQuestions[i], shuffledQuestions[j] = shuffledQuestions[j], shuffledQuestions[i]
	})
	
	// Create new quiz session with shuffled questions
	session := &models.QuizSession{
		ID:        sessionID,
		Questions: shuffledQuestions,
		Score:     0,
		StartedAt: time.Now(),
	}
	
	// Store session in centralized store
	store.Mu.Lock()
	store.Sessions[sessionID] = session
	userAnswers[sessionID] = make(map[int]string)
	store.Mu.Unlock()
	
	// Convert to public questions (without answers)
	publicQuestions := make([]models.PublicQuestion, len(questions))
	for i, q := range questions {
		publicQuestions[i] = models.PublicQuestion{
			ID:   q.ID,
			Text: q.Text,
		}
	}
	
	response := map[string]interface{}{
		"session_id": sessionID,
		"questions":  publicQuestions,
		"message":    "Quiz started successfully",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	
	log.Printf("Started quiz session: %s", sessionID)
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