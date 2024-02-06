package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Anime struct {
	Id          primitive.ObjectID `bson:"_id" json:"_id"`
	Name        string             `bson:"name" json:"name"`
	Genre       []string           `bson:"genre" json:"genre"`
	Description string             `bson:"description" json:"description"`
	Image       string             `bson:"image" json:"image"`
	Status      string             `bson:"status" json:"status"`
	Created_at  time.Time          `bsoon:"created_at" json:"created_at"`
}
