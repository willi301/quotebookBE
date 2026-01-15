package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Question struct {
	ID       int    `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Context  string `json:"context"`
}

type QuestionList struct {
	Questions []Question `json:"questions"`
}

type AnswerRequest struct {
	QuestionID int    `json:"question_id"`
	Answer     string `json:"answer"`
}

func fetchQuizData(ctx context.Context) (*QuestionList, error) {
	credsJSON := os.Getenv("GOOGLE_CREDENTIALS")
	if credsJSON == "" {
		return nil, fmt.Errorf("GOOGLE_CREDENTIALS not set")
	}

	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON([]byte(credsJSON)))
	if err != nil {
		return nil, fmt.Errorf("failed to create sheets client: %w", err)
	}

	spreadsheetID := os.Getenv("SHEET_ID")
	if spreadsheetID == "" {
		return nil, fmt.Errorf("SHEET_ID not set")
	}

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, "Sheet1!A:C").Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sheet data: %w", err)
	}

	var questionList QuestionList

	for i, row := range resp.Values {
		if i == 0 || len(row) < 2 {
			continue
		}

		question := strings.TrimSpace(fmt.Sprintf("%v", row[0]))
		answer := strings.TrimSpace(fmt.Sprintf("%v", row[1]))

		contextVal := ""
		if len(row) >= 3 {
			contextVal = strings.TrimSpace(fmt.Sprintf("%v", row[2]))
		}

		if question == "" || answer == "" {
			continue
		}

		questionList.Questions = append(questionList.Questions, Question{
			ID:       i,
			Question: question,
			Answer:   answer,
			Context:  contextVal,
		})
	}

	return &questionList, nil
}

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