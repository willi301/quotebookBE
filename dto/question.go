package dto

type Question struct {
	ID       int    `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Context  string `json:"context"`
}

type QuestionList struct {
	Questions []Question `json:"questions"`
}
