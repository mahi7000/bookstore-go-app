package handlers

import (
	"net/http"

	"github.com/mahi7000/bookstore-go-app/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *ApiConfig)AuthMiddleware(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := GetTokenFromHeader(r)
		if err != nil {
			RespondWithError(w, 403, "Unauthorized")
			return
		}

		userId, err := ValidateJWT(tokenString)
		if err != nil {
			RespondWithError(w, 403, "Token is invalid")
			return
		}

		user, err := apiCfg.DB.GetUserByID(r.Context(), userId)
		if err != nil {
			RespondWithError(w, 400, "Couldn't get user from id")
			return
		}

		handler(w, r, user)
	}
}