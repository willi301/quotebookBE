package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"backend/fetcher"
	"backend/handler"
	"backend/store"
)

// Middleware to allow the Frontend to talk to the Backend
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow the Frontend's specific URL (Change to vercel soon)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")

		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle "Preflight" requests (Browser checking permission before sending data)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Seed random ONCE
	rand.Seed(time.Now().UnixNano())

	// Load quiz data ONCE
	qs, err := store.LoadQuestions("quiz_data.json")
	if err != nil {
		log.Fatal("failed to load questions:", err)
	}
	store.Questions = qs

	// Setup router
	mux := http.NewServeMux()
	mux.HandleFunc("/api/quiz/start", handler.StartQuiz)
	mux.HandleFunc("/api/quiz/answer", handler.CheckAnswer)
	mux.HandleFunc("/api/quiz/refresh", fetcher.UpdateQuizData)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(mux)))
}
