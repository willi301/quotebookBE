package fetcher

import (
	"backend/backend/dto"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func UpdateQuizData() {
	spreadsheetID := "1WJlV_qezd_aiWnTeCpTlDDxFaZx2ipkdI9rM3KItbzM"
	readRange := "Sheet1!A:C"

	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile("credentials.json"))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data: %v", err)
	}

	var questionList dto.QuestionList

	for i, row := range resp.Values {
		if i == 0 {
			continue
		}

		if len(row) < 2 {
			continue
		}

		question := strings.TrimSpace(fmt.Sprintf("%v", row[0]))
		answer := strings.TrimSpace(fmt.Sprintf("%v", row[1]))

		context := ""
		if len(row) >= 3 {
			context = strings.TrimSpace(fmt.Sprintf("%v", row[2]))
		}

		if question == "" || answer == "" {
			continue
		}

		questionList.Questions = append(questionList.Questions, dto.Question{
			ID:       i,
			Question: question,
			Answer:   answer,
			Context:  context,
		})
	}

	file, _ := json.MarshalIndent(questionList, "", "  ")

	err = os.WriteFile("quiz_data.json", file, 0644)

	if err != nil {
		log.Fatalf("Error writing file: %v", err)
	}

	fmt.Println("✅ Data updated inside fetcher package")
}
