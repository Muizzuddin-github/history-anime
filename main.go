package main

import (
	"context"
	"fmt"
	"history_anime/src/db"
	_ "history_anime/src/db"
	"history_anime/src/routers"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Orang struct {
	Nama string
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()
	db.CreateConnection(ctx)
	defer db.CloseDB(ctx)

	server := http.Server{
		Handler: routers.Router(),
		Addr:    "localhost:" + port,
	}

	fmt.Println("server is running on http://localhost:" + port)
	server.ListenAndServe()

}
