package handler

type Question struct {
	ID       int    `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Context  string `json:"context"`
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

type QuestionList struct {
	Questions []Question `json:"questions"`
}