package dbutility

import (
	"context"
	"errors"
	"history_anime/src/db"
	"history_anime/src/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteUser(email string) error {
	ctx := context.Background()
	filter := bson.D{
		{
			Key:   "email",
			Value: email,
		},
	}
	result, err := db.DB.Collection("users").DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("not found")
	}
	return nil
}

func FindUser(email string) (string, error) {

	ctx := context.Background()
	filter := bson.D{
		{
			Key:   "email",
			Value: email,
		},
	}
	data := entity.Users{}
	err := db.DB.Collection("users").FindOne(ctx, filter).Decode(&data)
	if err == mongo.ErrNoDocuments {
		return "", errors.New("not found")
	}

	return data.Id.Hex(), nil

}
