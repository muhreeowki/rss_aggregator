package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/muhreeowki/rss_aggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not found in env.")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not found in env.")
	}

	// Cconnecting to DB
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Couldnt connect to DB")
	}

	queries := database.New(conn)

	apiCfg := apiConfig{
		DB: queries,
	}

	// Routers
	router := chi.NewRouter()
	v1Router := chi.NewRouter()

	// Enable Cors
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"LINK"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Routes and Handlers
	v1Router.Get("/ready", handlerReadiness)
	v1Router.Get("/error", handlerError)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	router.Mount("/v1", v1Router)

	// Setup server http instance
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	// Start server
	fmt.Printf("Server listening on PORT %v...", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
