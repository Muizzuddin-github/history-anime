package repository

import (
	"context"
	"errors"
	"history_anime/src/entity"
	"history_anime/src/requestbody"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type genreInterface interface {
	GetAll(ctx context.Context) ([]entity.Genre, error)
	Add(ctx context.Context, body *requestbody.Genre) (string, error)
	Del(ctx context.Context, id string) (*mongo.DeleteResult, error)
}

type genreRepo struct {
	DB *mongo.Database
}

func (genre *genreRepo) GetAll(ctx context.Context) ([]entity.Genre, error) {

	sortDoc := bson.D{
		{
			Key: "$sort",
			Value: bson.D{
				{
					Key:   "created_at",
					Value: -1,
				},
			},
		},
	}

	cur, err := genre.DB.Collection("genre").Aggregate(ctx, mongo.Pipeline{sortDoc})
	if err != nil {
		return []entity.Genre{}, errors.New(err.Error())
	}
	defer cur.Close(ctx)

	genres := []entity.Genre{}
	for cur.Next(ctx) {
		data := entity.Genre{}

		err := cur.Decode(&data)
		if err != nil {
			return []entity.Genre{}, errors.New(err.Error())
		}

		genres = append(genres, data)
	}

	return genres, nil

}

func (genre *genreRepo) Add(ctx context.Context, body *requestbody.Genre) (string, error) {

	insertDoc := bson.D{
		{
			Key:   "name",
			Value: body.Name,
		},
		{
			Key:   "created_at",
			Value: primitive.NewDateTimeFromTime(time.Now()),
		},
	}
	insert, err := genre.DB.Collection("genre").InsertOne(ctx, insertDoc)
	if err != nil {
		return "", errors.New(err.Error())
	}

	insertedID, ok := insert.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New(err.Error())
	}

	return insertedID.Hex(), nil
}

func (genre *genreRepo) Del(ctx context.Context, id string) (*mongo.DeleteResult, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	filter := bson.D{
		{
			Key:   "_id",
			Value: objID,
		},
	}
	result, err := genre.DB.Collection("genre").DeleteOne(ctx, filter)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return result, nil
}

func GenreRepo(db *mongo.Database) genreInterface {
	return &genreRepo{
		DB: db,
	}
}
