package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         string `json:"id"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
	Name       string `json:"name"`
}

func handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode request body")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:         uuid.NewString(),
		Created_at: time.Now().Format(time.RFC3339),
		Updated_at: time.Now().Format(time.RFC3339),
		Name:       params.Name,
	})
}
