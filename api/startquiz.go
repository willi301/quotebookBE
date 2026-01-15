package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", 405)
		return
	}

	data, err := fetchQuizData(r.Context())
	if err != nil {
		http.Error(w, "Failed to load quiz: "+err.Error(), 500)
		return
	}

	rand.Shuffle(len(data.Questions), func(i, j int) {
		data.Questions[i], data.Questions[j] = data.Questions[j], data.Questions[i]
	})

	public := make([]PublicQuestion, len(data.Questions))
	for i, q := range data.Questions {
		public[i] = PublicQuestion{
			ID:      q.ID,
			Text:    q.Question,
			Context: q.Context,
		}
	}

	json.NewEncoder(w).Encode(public)
}
