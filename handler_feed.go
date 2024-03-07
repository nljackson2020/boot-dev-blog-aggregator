package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nljackson2020/boot-dev-blog-aggregator/internal/database"
)

func (cfg *apiConfig) handlerFeedCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	type ReturnObject struct {
		Feed        Feed          `json:"feed"`
		Feed_Follow []FeedsFollow `json:"feed_follow"`
	}

	var feed ReturnObject

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode request body")
		return
	}

	databaseFeed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create feed")
		return
	}

	databaseFeedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    databaseFeed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to associate feed to user")
		return
	}
	feed.Feed = databaseFeedToFeed(databaseFeed)
	feed.Feed_Follow = make([]FeedsFollow, len(databaseFeedFollow))
	for i, follow := range databaseFeedFollow {
		feed.Feed_Follow[i] = databaseFeedsFollowToFeedsFollow(follow)
	}

	respondWithJSON(w, http.StatusCreated, feed)
}

func (cfg *apiConfig) handlerFeedGetAll(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeedsAll(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get feeds")
		return
	}
	respondWithJSON(w, http.StatusOK, feeds)
}

func (cfg *apiConfig) handlerFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Feed_ID string `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode request body")
		return
	}
	feedID, err := uuid.Parse(params.Feed_ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	feed_follow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feedID,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to associate feed to user")
		return
	}

	respondWithJSON(w, http.StatusOK, feed_follow)
}

func (cfg *apiConfig) handlerFeedFollowsDelete(w http.ResponseWriter, r *http.Request, user database.User) {
	feedIDString := chi.URLParam(r, "feedFollowID")
	feedID, err := uuid.Parse(feedIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID format")
		return
	}
	err = cfg.DB.DeleteFeedFollow(r.Context(), feedID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete feed follow")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (cfg *apiConfig) handlerGetFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := cfg.DB.GetUserFeedAll(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "No user found with that ID")
		return
	}
	respondWithJSON(w, http.StatusOK, feeds)
}
