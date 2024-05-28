package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/muhreeowki/rss_aggregator/internal/database"
)

func (apiCfg *apiConfig) handlerCreateUser(writer http.ResponseWriter, r *http.Request) {
	type params struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	p := params{}
	err := decoder.Decode(&p)
	if err != nil {
		respondWithError(writer, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.CreatUser(r.Context(), database.CreatUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      p.Name,
	})
	if err != nil {
		respondWithError(writer, 400, fmt.Sprintf("Couldnt create user: %v", err))
		return
	}

	respondWithJSON(writer, 201, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUser(writer http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(writer, 200, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetPostsForUser(writer http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondWithError(writer, 400, fmt.Sprintf("Couldnt get posts for user: %v", err))
		return
	}
	respondWithJSON(writer, 200, databasePostsToPosts(posts))
}
