package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rachit2501/goserver/internal/database"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){

	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)

	if err!= nil {
		responseWithError(w,400, fmt.Sprintf("Error parsing JSON %v",err))
		return
	}

	user, err:= apiCfg.DB.CreateUser(r.Context(),database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
	})

	if err!= nil {
		responseWithError(w,400, fmt.Sprintf("Couldn't create user %v", err))
	}

	
	respondWithJSON(w,200,databaseUserToUser(user))
}