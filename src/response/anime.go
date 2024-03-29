package response

import "history_anime/src/entity"

type AnimeAll struct {
	Message string         `json:"message"`
	Data    []entity.Anime `json:"data"`
}

type AnimeInsert struct {
	Message    string `json:"message"`
	InsertedID string `json:"insertedID"`
}
