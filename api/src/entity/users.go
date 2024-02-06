package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Users struct {
	Id       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}
