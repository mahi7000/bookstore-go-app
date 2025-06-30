package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println("Couldn't marshal json")
		w.WriteHeader(400)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Severe error most likely server error: ", code)
	}
	
	type errorMsg struct {
		Error string `json:"error"`
	}
	RespondWithJson(w, code, errorMsg{
		Error: msg,
	})
}