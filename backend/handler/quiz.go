package handler

import (
	"backend/models"
	"backend/store"
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"backend/fetcher"
	"backend/models"
	"backend/store"
)

// StartQuiz initializes a new quiz
func StartQuiz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Update quiz data from Google Sheets and reload questions
	fetcher.UpdateQuizData()
	qs, err := store.LoadQuestions("quiz_data.json")
	if err != nil {
		http.Error(w, "Failed to load questions", http.StatusInternalServerError)
		return
	}
	store.Mu.Lock()
	store.Questions = qs
	store.Mu.Unlock()

	// Copy questions from store (DO NOT mutate)
	shuffled := make([]models.Question, len(store.Questions))
	copy(shuffled, store.Questions)

	// Shuffle
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

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
		"questions": public,
	})
}

// to check whether the answer is correct
func CheckAnswer(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.AnswerRequest
	json.NewDecoder(r.Body).Decode(&req)

	for _, q := range store.Questions {
		if q.ID == req.QuestionID {
			isCorrect := strings.EqualFold(req.Answer, q.Answer)
			json.NewEncoder(w).Encode(map[string]bool{
				"correct": isCorrect,
			})
			return
		}
	}

	http.Error(w, "Question not found", http.StatusNotFound)
}
