package response

import "history_anime/src/entity"

type AnimeAll struct {
	Message string         `json:"message"`
	Data    []entity.Anime `json:"data"`
}
