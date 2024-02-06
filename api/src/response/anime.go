package response

import "history_anime/api/src/entity"

type AnimeAll struct {
	Message string         `json:"message"`
	Data    []entity.Anime `json:"data"`
}
