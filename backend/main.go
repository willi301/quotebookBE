package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"backend/handler"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	//create router
	mux := http.NewServeMux()

	//Register routes
	mux.HandleFunc("/api/quiz/start", handler.StartQuiz)
	// mux.HandleFunc("/api/quiz/answer", handlers.SubmitAnswer)
	// mux.HandleFunc("/api/quiz/finish", handlers.FinishQuiz)

	// // (optional health check)
	// mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("ok"))
	// })

	// 4️⃣ Start server
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
