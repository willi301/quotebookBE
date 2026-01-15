package handler

import (
	"context"
	"fmt"
	"os"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Shared types
type Question struct {
	ID       int    `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Context  string `json:"context"`
}

type QuestionList struct {
	Questions []Question `json:"questions"`
}

type PublicQuestion struct {
	ID      int    `json:"id"`
	Text    string `json:"text"`
	Context string `json:"context,omitempty"`
}

type AnswerRequest struct {
	QuestionID int    `json:"question_id"`
	Answer     string `json:"answer"`
}

// Shared fetcher function
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