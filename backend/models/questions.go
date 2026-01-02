package models

type Question struct {
	ID     int
	Text   string
	Answer string
}

type PublicQuestion struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}