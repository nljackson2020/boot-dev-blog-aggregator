package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nljackson2020/boot-dev-blog-aggregator/internal/database"
)

type User struct {
	ID         string `json:"id"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
	Name       string `json:"name"`
}

func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode request body")
		return
	}

	user, err := cfg.DB.CreatUser(r.Context(), database.CreatUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) handlerUserGet(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	apiKey := splitAuth[1]

	users, err := cfg.DB.GetUsers(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get users")
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}
