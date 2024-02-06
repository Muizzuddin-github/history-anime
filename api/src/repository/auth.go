package repository

import (
	"context"
	"errors"
	"history_anime/api/src/entity"
	"history_anime/api/src/requestbody"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type authRepoInterface interface {
	Register(ctx context.Context, body *requestbody.Register) error
	Login(ctx context.Context, body *requestbody.Login) (*entity.Users, error)
}

type authRepo struct {
	DB *mongo.Database
}

func (auth *authRepo) Register(ctx context.Context, body *requestbody.Register) error {

	_, err := auth.DB.Collection("users").InsertOne(ctx, body)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (auth *authRepo) Login(ctx context.Context, body *requestbody.Login) (*entity.Users, error) {

	filter := bson.D{{Key: "email", Value: body.Email}}

	result := entity.Users{}
	err := auth.DB.Collection("users").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return &entity.Users{}, errors.New(err.Error())
	}

	return &result, nil
}

func AuthRepo(db *mongo.Database) authRepoInterface {
	return &authRepo{
		DB: db,
	}
}
