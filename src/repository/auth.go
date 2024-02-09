package repository

import (
	"context"
	"errors"
	"history_anime/src/entity"
	"history_anime/src/requestbody"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type authRepoInterface interface {
	Register(ctx context.Context, body *requestbody.Register) error
	Login(ctx context.Context, body *requestbody.Login) (*entity.Users, error)
	ResetPassword(ctx context.Context, email string, newHashPassword string) (*mongo.UpdateResult, error)
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

func (auth *authRepo) ResetPassword(ctx context.Context, email string, newHashPassword string) (*mongo.UpdateResult, error) {

	filter := bson.D{
		{
			Key:   "email",
			Value: email,
		},
	}

	updateDoc := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "password",
					Value: newHashPassword,
				},
			},
		},
	}

	up, err := auth.DB.Collection("users").UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return up, nil
}

func AuthRepo(db *mongo.Database) authRepoInterface {
	return &authRepo{
		DB: db,
	}
}
