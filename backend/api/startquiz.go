package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"backend/models"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", 405)
		return
	}

	data, err := fetcher.FetchQuizData(r.Context())
	if err != nil {
		http.Error(w, "Failed to load quiz", 500)
		return
	}

	// shuffle
	rand.Shuffle(len(data.Questions), func(i, j int) {
		data.Questions[i], data.Questions[j] = data.Questions[j], data.Questions[i]
	})

	// strip answers
	public := make([]models.PublicQuestion, len(data.Questions))
	for i, q := range data.Questions {
		public[i] = models.PublicQuestion{
			ID:   q.ID,
			Text: q.Question,
		}
	}

	json.NewEncoder(w).Encode(public)
}