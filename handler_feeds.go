package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/muhreeowki/rss_aggregator/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFeed(writer http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	p := params{}
	err := decoder.Decode(&p)
	if err != nil {
		respondWithError(writer, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      p.Name,
		Url:       p.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(writer, 400, fmt.Sprintf("Couldnt create user: %v", err))
		return
	}

	respondWithJSON(writer, 201, databaseFeedToFeed(feed))
}

// func (apiCfg *apiConfig) handlerGetFeed(writer http.ResponseWriter, r *http.Request, user database.User) {
// 	respondWithJSON(writer, 200, databaseUserToUser(user))
// }
