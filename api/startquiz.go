package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type quizQuestion struct {
	ID       int    `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Context  string `json:"context"`
}

type quizQuestionList struct {
	Questions []quizQuestion `json:"questions"`
}

type publicQuestion struct {
	ID      int    `json:"id"`
	Text    string `json:"text"`
	Context string `json:"context,omitempty"`
}

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

	data, err := fetchQuizQuizData(r.Context())
	if err != nil {
		http.Error(w, "Failed to load quiz: "+err.Error(), 500)
		return
	}

	rand.Shuffle(len(data.Questions), func(i, j int) {
		data.Questions[i], data.Questions[j] = data.Questions[j], data.Questions[i]
	})

	public := make([]publicQuestion, len(data.Questions))
	for i, q := range data.Questions {
		public[i] = publicQuestion{
			ID:      q.ID,
			Text:    q.Question,
			Context: q.Context,
		}
	}

	json.NewEncoder(w).Encode(public)
}

func fetchQuizQuizData(ctx context.Context) (*quizQuestionList, error) {
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

	var questionList quizQuestionList

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

		questionList.Questions = append(questionList.Questions, quizQuestion{
			ID:       i,
			Question: question,
			Answer:   answer,
			Context:  contextVal,
		})
	}

	return &questionList, nil
}