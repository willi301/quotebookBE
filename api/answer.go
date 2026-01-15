package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var req AnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", 400)
		return
	}

	data, err := fetchQuizData(r.Context())
	if err != nil {
		http.Error(w, "Failed to load quiz data: "+err.Error(), 500)
		return
	}

	for _, q := range data.Questions {
		if q.ID == req.QuestionID {
			json.NewEncoder(w).Encode(map[string]bool{
				"correct": strings.EqualFold(q.Answer, req.Answer),
			})
			return
		}
	}

	http.Error(w, "Question not found", 404)
}