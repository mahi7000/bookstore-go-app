package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mahi7000/bookstore-go-app/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (apiCfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name	 string `json:"name"`
		Email	 string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	decoder.Decode(&params)

	hashed_password, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		RespondWithError(w, 403, "Unable to hash password")
		return
	}

	user, err := apiCfg.DB.CreateNewUser(r.Context(), database.CreateNewUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: params.Name,
		Email: params.Email,
		Password: string(hashed_password),
	})
	if err != nil {
		RespondWithError(w, 403, "Unables to create new user")
		return
	}

	token, err := GenerateJWT(user.ID)
	if err != nil {
		RespondWithError(w, 500, "Couldn't generate the token")
	}

	fmt.Println(token)
	
	RespondWithJson(w, 201, user)
}

func (apiCfg *ApiConfig)HandlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	RespondWithJson(w, 200, user)
}

func (apiCfg *ApiConfig) HandlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 	 string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 403, "Invalid parameters")
		return
	}

	user, err := apiCfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, 500, "Couldn't get user")
		return
	}

	if !CheckPasswordHash(params.Password, user.Password) {
		RespondWithError(w, http.StatusUnauthorized, "Invalid Credentials")
		return
	}

	if tokenString, err := GetTokenFromHeader(r); err == nil {
		if userID, err := ValidateJWT(tokenString); err == nil {
			if userID == user.ID {
				RespondWithError(w, 403, "User already logged in")
				return
			}
		}
	}

	token, err := GenerateJWT(user.ID)
	if err != nil {
		RespondWithError(w, 500, "Unable to generate token")
		return
	}

	fmt.Println(token)

	RespondWithJson(w, 201, user)
}

func (apiCfg *ApiConfig)HandlerLogoutUser(w http.ResponseWriter, r *http.Request, user database.User) {
	token, err := GetTokenFromHeader(r)
	if err != nil {
		RespondWithError(w, 403, "User not logged in"+err.Error())
	}
	
	err = RevokeJWT(token)
	if err != nil {
		RespondWithError(w, 403, "Failed to revoke token")
	}

	RespondWithJson(w, 200, "Successfully logged out")
}