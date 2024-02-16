package dbutility

import (
	"context"
	"errors"
	"history_anime/src/db"
	"history_anime/src/requestbody"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GenreDeleteById(id string) error {

	ctx := context.Background()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filters := bson.D{
		{
			Key:   "_id",
			Value: objID,
		},
	}

	result, err := db.DB.Collection("genre").DeleteOne(ctx, filters)
	if err != nil {
		return errors.New(err.Error())
	}

	if result.DeletedCount == 0 {
		return errors.New("not found")
	}

	return nil

}

func GenreDeleteByName(name string) error {

	ctx := context.Background()

	filters := bson.D{
		{
			Key:   "name",
			Value: name,
		},
	}

	result, err := db.DB.Collection("genre").DeleteOne(ctx, filters)
	if err != nil {
		return errors.New(err.Error())
	}

	if result.DeletedCount == 0 {
		return errors.New("not found")
	}

	return nil

}

func GenreAdd(data *requestbody.Genre) (string, error) {
	ctx := context.Background()

	insertDoc := bson.D{
		{
			Key:   "name",
			Value: data.Name,
		},
	}
	result, err := db.DB.Collection("genre").InsertOne(ctx, insertDoc)
	if err != nil {
		return "", errors.New(err.Error())
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("type error")
	}

	return id.Hex(), nil
}
