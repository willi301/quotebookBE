package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)


type Question struct {
	ID       int    `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type Questions struct {
	Questions []Question `json:"questions"`
}


func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"message":"Hello from Go!"}`)
	})

	mux.HandleFunc("/api/json", func(w http.ResponseWriter, r *http.Request) {
		// Open the questions.json file
		jsonFile, err := os.Open("questions.json")
		if err != nil {
			http.Error(w, "failed to open questions.json", http.StatusInternalServerError)
			log.Println("error opening questions.json:", err)
			return
		}
		defer jsonFile.Close()

		var qs Questions
		// Creates decoder
		dec := json.NewDecoder(jsonFile)
		// Decode JSON data
		if err := dec.Decode(&qs); err != nil {
			http.Error(w, "failed to decode questions.json", http.StatusInternalServerError)
			log.Println("error decoding questions.json:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(qs); err != nil {
			log.Println("error encoding response:", err)
		}
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server running on http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
