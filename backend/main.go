package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"backend/handler"
	"backend/store"
)

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

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}