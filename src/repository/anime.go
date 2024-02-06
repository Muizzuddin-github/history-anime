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

type animeRepoInterface interface {
	Add(ctx context.Context, body requestbody.Anime) (string, error)
	Update(ctx context.Context, body requestbody.Anime, id string) (*mongo.UpdateResult, error)
	Del(ctx context.Context, id string) (*mongo.DeleteResult, error)
	GetAll(ctx context.Context) ([]entity.Anime, error)
}

type animeRepo struct {
	DB *mongo.Database
}

func (anime *animeRepo) Add(ctx context.Context, body requestbody.Anime) (string, error) {

	insertDoc := bson.D{
		{Key: "name", Value: body.Name},
		{Key: "image", Value: body.Image},
		{Key: "genre", Value: body.Genre},
		{Key: "description", Value: body.Description},
		{Key: "status", Value: body.Status},
		{Key: "created_at", Value: primitive.NewDateTimeFromTime(time.Now())},
	}
	insert, err := anime.DB.Collection("anime").InsertOne(ctx, insertDoc)
	if err != nil {
		return "", errors.New(err.Error())
	}

	insertID, ok := insert.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New(err.Error())
	}

	return insertID.Hex(), nil
}

func (anime *animeRepo) Update(ctx context.Context, body requestbody.Anime, id string) (*mongo.UpdateResult, error) {

	upDoc := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: body.Name},
			{Key: "image", Value: body.Image},
			{Key: "genre", Value: body.Genre},
			{Key: "description", Value: body.Description},
			{Key: "status", Value: body.Status},
		}},
	}

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	up, err := anime.DB.Collection("anime").UpdateByID(ctx, objId, upDoc)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return up, nil
}

func (anime *animeRepo) Del(ctx context.Context, id string) (*mongo.DeleteResult, error) {

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	filter := bson.D{
		{
			Key:   "_id",
			Value: objId,
		},
	}
	result, err := anime.DB.Collection("anime").DeleteOne(ctx, filter)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return result, nil
}

func (anime *animeRepo) GetAll(ctx context.Context) ([]entity.Anime, error) {

	cur, err := anime.DB.Collection("anime").Find(ctx, bson.D{})
	defer cur.Close(ctx)

	if err != nil {
		return []entity.Anime{}, errors.New(err.Error())
	}

	result := []entity.Anime{}

	for cur.Next(ctx) {
		data := entity.Anime{}
		err := cur.Decode(&data)
		if err != nil {
			return []entity.Anime{}, errors.New(err.Error())
		}

		result = append(result, data)
	}

	return result, nil
}

func AnimeRepo(db *mongo.Database) animeRepoInterface {
	return &animeRepo{
		DB: db,
	}
}
