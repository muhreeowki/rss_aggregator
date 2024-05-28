package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/muhreeowki/rss_aggregator/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(writer http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	p := params{}
	err := decoder.Decode(&p)
	if err != nil {
		respondWithError(writer, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    p.FeedID,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(writer, 400, fmt.Sprintf("Couldnt create feed follow: %v", err))
		return
	}

	respondWithJSON(writer, 201, databaseFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *apiConfig) handlerGetFeedFollows(writer http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(writer, 400, fmt.Sprintf("Couldnt get feed follows: %v", err))
		return
	}

	respondWithJSON(writer, 200, databaseFeedFollowsToFeedFollows(feedFollows))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollow(writer http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(writer, 400, fmt.Sprintf("Couldnt parse feed follow id: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollows(r.Context(), database.DeleteFeedFollowsParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(writer, 400, fmt.Sprintf("Couldnt delete feed follow: %v", err))
		return
	}

	respondWithJSON(writer, 200, struct{}{})
}
