package main

import (
	"fmt"
	"net/http"

	"github.com/muhreeowki/rss_aggregator/internal/auth"
	"github.com/muhreeowki/rss_aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		// Check caller has passed APIKEY in header
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(writer, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}
		// Get the user
		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(writer, 400, fmt.Sprintf("Couldnt find user: %v", err))
			return
		}

		// Run handler
		handler(writer, r, user)
	}
}
