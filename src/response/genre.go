package response

import "history_anime/src/entity"

type GenreAll struct {
	Message string         `json:"message"`
	Data    []entity.Genre `json:"data"`
}

type GenreInsert struct {
	Message    string `json:"message"`
	InsertedID string `json:"insertedID"`
}
