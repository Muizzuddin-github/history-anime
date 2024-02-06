package handler

import (
	"context"
	"history_anime/api/src/db"
	_ "history_anime/api/src/db"
	"history_anime/api/src/routers"

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
	defer db.CloseDB(ctx)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})

	server = http.Server{
		Handler: c.Handler(routers.Router()),
		Addr:    "localhost:" + port,
	}

}

func Handler(w http.ResponseWriter, r *http.Request) {
	server.Handler.ServeHTTP(w, r)
}
