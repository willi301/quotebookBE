package fetcher

import (
	"backend/dto"
	"context"
	"fmt"
	"os"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func FetchQuizData(ctx context.Context) (*dto.QuestionList, error) {
	srv, err := sheets.NewService(
		ctx,
		option.WithCredentialsJSON([]byte(os.Getenv("GOOGLE_CREDENTIALS"))),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create sheets client: %w", err)
	}

	spreadsheetID := os.Getenv("SHEET_ID")
	readRange := "Sheet1!A:C"

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sheet data: %w", err)
	}

	var questionList dto.QuestionList

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

		questionList.Questions = append(questionList.Questions, dto.Question{
			ID:       i,
			Question: question,
			Answer:   answer,
			Context:  contextVal,
		})
	}

	return &questionList, nil
}

