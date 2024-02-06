package db

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func CreateConnection(ctx context.Context) {

	option := options.Client().ApplyURI(os.Getenv("DATABASE_URI"))
	client, err := mongo.Connect(ctx, option)

	if err != nil {
		client.Disconnect(ctx)
		return
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	db := client.Database("history-anime")
	DB = db
	fmt.Println("database terhubung")
}

func CloseDB(ctx context.Context) {
	DB.Client().Disconnect(ctx)
}
