package store

import (
	"encoding/json"
	"os"

	"backend/models"
)

type QuizData struct {
	Questions []models.Question `json:"questions"`
}

func LoadQuestions(path string) ([]models.Question, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data QuizData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}

	return data.Questions, nil
}
