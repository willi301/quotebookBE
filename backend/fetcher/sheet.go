package fetcher

import (
	"backend/dto"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func UpdateQuizData(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile("credentials.json"))
	if err != nil {
		log.Printf("Unable to retrieve Sheets client: %v", err)
		http.Error(w, "Internal Server Error: Could not connect to Google Sheets", http.StatusInternalServerError)
		return
	}

	spreadsheetID := "1WJlV_qezd_aiWnTeCpTlDDxFaZx2ipkdI9rM3KItbzM"
	readRange := "Sheet1!A:C"

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Printf("Unable to retrieve data: %v", err)
		http.Error(w, "Failed to fetch data from Google Sheets", http.StatusInternalServerError)
		return
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

	file, err := json.MarshalIndent(questionList, "", "  ")
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		http.Error(w, "Failed to process data", http.StatusInternalServerError)
		return
	}

	err = os.WriteFile("quiz_data.json", file, 0644)
	if err != nil {
		log.Printf("Error writing file: %v", err)
		http.Error(w, "Failed to save data to disk", http.StatusInternalServerError)
		return
	}

	log.Println("✅ Data updated successfully via API request")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Quiz data refreshed successfully", "count": ` + fmt.Sprint(len(questionList.Questions)) + `}`))
}
