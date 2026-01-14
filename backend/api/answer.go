package Handler

import (
	"encoding/json"
	"net/http"

	"backend/models"
)

type AnswerRequest struct {
	QuestionID int    `json:"question_id"`
	Answer     string `json:"answer"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var req AnswerRequest
	json.NewDecoder(r.Body).Decode(&req)

	data, _ := fetcher.FetchQuizData(r.Context())

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