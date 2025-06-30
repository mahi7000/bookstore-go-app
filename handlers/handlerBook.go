package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	
	"github.com/mahi7000/bookstore-go-app/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *ApiConfig)HandlerAddBook(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Author      string `json:"author"`
		Cover       string `json:"cover"`
		URL         string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 400, "Couldn't decode params")
		return
	}

	book, err := apiCfg.DB.AddNewBook(r.Context(), database.AddNewBookParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Name:        params.Name,
		Description: params.Description,
		Author:      params.Author,
		Url:         params.URL,
		BookCover: 	 params.Cover,
	})
	if err != nil {
		RespondWithError(w, 400, "Couldn't add new book")
		return
	}

	RespondWithJson(w, 201, book)
}