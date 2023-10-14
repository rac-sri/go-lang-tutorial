package main

import (
	"fmt"
	"net/http"

	"github.com/rachit2501/goserver/internal/auth"
	"github.com/rachit2501/goserver/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		apiKey, err:= auth.GetAPIKey(r.Header)
		if err!=nil {
			responseWithError(w,403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user,err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			responseWithError(w,400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}

		handler(w, r, user)
	}
}