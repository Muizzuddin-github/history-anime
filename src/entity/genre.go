package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Genre struct {
	Id         primitive.ObjectID `bson:"_id" json:"_id"`
	Name       string             `bson:"name" json:"name"`
	Created_at time.Time          `bson:"created_at" json:"created_at"`
}
