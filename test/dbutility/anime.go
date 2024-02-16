package dbutility

import (
	"context"
	"errors"
	"history_anime/src/db"
	"history_anime/src/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AnimeFindOne(name string) (string, error) {

	ctx := context.Background()

	filter := bson.D{
		{
			Key:   "name",
			Value: name,
		},
	}

	data := entity.Anime{}
	err := db.DB.Collection("anime").FindOne(ctx, filter).Decode(&data)
	if err == mongo.ErrNoDocuments {
		return "", errors.New("not found")
	} else if err != nil {
		return "", errors.New(err.Error())
	}

	return data.Id.Hex(), nil
}

func AnimeDeleteOne(name string) error {
	ctx := context.Background()

	filter := bson.D{
		{
			Key:   "name",
			Value: name,
		},
	}

	result, err := db.DB.Collection("anime").DeleteOne(ctx, filter)
	if err != nil {
		return errors.New(err.Error())
	}

	if result.DeletedCount == 0 {
		return errors.New("not found")
	}

	return nil
}

func AnimeDeleteOneById(id string) error {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(err.Error())
	}

	ctx := context.Background()
	filter := bson.D{
		{
			Key:   "_id",
			Value: objID,
		},
	}

	result, err := db.DB.Collection("anime").DeleteOne(ctx, filter)
	if err != nil {
		return errors.New(err.Error())
	}

	if result.DeletedCount == 0 {
		return errors.New("not found")
	}

	return nil
}

func AnimeAdd(name string, image string, genre []string, description string, status string) (string, error) {

	insertDoc := bson.D{
		{Key: "name", Value: name},
		{Key: "image", Value: image},
		{Key: "genre", Value: genre},
		{Key: "description", Value: description},
		{Key: "status", Value: status},
		{Key: "created_at", Value: primitive.NewDateTimeFromTime(time.Now())},
	}

	ctx := context.Background()

	insert, err := db.DB.Collection("anime").InsertOne(ctx, insertDoc)
	if err != nil {
		return "", errors.New(err.Error())
	}

	insertID, ok := insert.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New(err.Error())
	}

	return insertID.Hex(), nil

}
