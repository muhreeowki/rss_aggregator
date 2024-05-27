package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not found in env.")
	}

	router := chi.NewRouter()

	// Enable Cors
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"LINK"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Setup server http instance
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	// Start server
	fmt.Printf("Server listening on PORT %v...", port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
