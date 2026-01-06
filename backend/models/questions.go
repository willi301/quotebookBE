package models

type Question struct {
	ID     int    `json:"id"`
	Text   string `json:"question"` // ✅ FIX
	Answer string `json:"answer"`
}


type PublicQuestion struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}