package handler

// package main

import (
	"context"

	"history_anime/src/db"
	"history_anime/src/routers"

	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
)

var server http.Server

func init() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()
	db.CreateConnection(ctx)
	// defer db.CloseDB(ctx)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	})

	server = http.Server{
		Handler: c.Handler(routers.Router()),
		Addr:    "localhost:" + port,
	}

	// server.ListenAndServe()

}

func Handler(w http.ResponseWriter, r *http.Request) {
	server.Handler.ServeHTTP(w, r)
}
