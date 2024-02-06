package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func CreateConnection(ctx context.Context) {

	option := options.Client().ApplyURI("mongodb://localhost:27017")
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
