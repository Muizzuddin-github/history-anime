package handler

import (
	"context"
	"fmt"
	"history_anime/api/src/db"
	_ "history_anime/api/src/db"
	"history_anime/api/src/routers"

	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
)

func main() {

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

	server := http.Server{
		Handler: c.Handler(routers.Router()),
		Addr:    "localhost:" + port,
	}

	fmt.Println("server is running on http://localhost:" + port)
	server.ListenAndServe()

}
