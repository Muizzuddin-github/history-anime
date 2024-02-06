package controllers

import (
	"context"
	"encoding/json"
	"history_anime/src/db"
	"history_anime/src/repository"
	"history_anime/src/requestbody"
	"history_anime/src/response"
	"history_anime/src/validation"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var AnimeAdd httprouter.Handle = func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {

		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	body := requestbody.Anime{}
	err = json.Unmarshal(bodyByte, &body)
	if err != nil {

		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	errResult := validation.ValidateAnime(&body)
	if len(*errResult) > 0 {
		res, _ := json.Marshal(response.Errors{
			Errors: *errResult,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	ctx := context.Background()
	anime := repository.AnimeRepo(db.DB)

	insertId, err := anime.Add(ctx, body)
	if err != nil {

		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	data := map[string]string{
		"message":    "berhasil menambah anime",
		"insertedID": insertId,
	}

	res, err := json.Marshal(data)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

var AnimeUpdate httprouter.Handle = func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {

		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	body := requestbody.Anime{}
	err = json.Unmarshal(bodyByte, &body)
	if err != nil {

		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	errResult := validation.ValidateAnime(&body)
	if len(*errResult) > 0 {
		res, _ := json.Marshal(response.Errors{
			Errors: *errResult,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	ctx := context.Background()
	anime := repository.AnimeRepo(db.DB)

	result, err := anime.Update(ctx, body, params.ByName("id"))
	if err != nil {

		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	if result.MatchedCount == 0 {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{"anime not found"},
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(res)
		return
	}

	res, err := json.Marshal(response.Msg{Message: "anime update success"})
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

var AnimeDel httprouter.Handle = func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	ctx := context.Background()
	anime := repository.AnimeRepo(db.DB)
	result, err := anime.Del(ctx, params.ByName("id"))
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	if result.DeletedCount == 0 {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{"anime tidak ada"},
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(res)
		return
	}

	res, _ := json.Marshal(response.Msg{
		Message: "anime delete success",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return

}

var AnimeGetAll httprouter.Handle = func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	ctx := context.Background()
	anime := repository.AnimeRepo(db.DB)
	result, err := anime.GetAll(ctx)

	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	res, _ := json.Marshal(response.AnimeAll{
		Message: "data anime",
		Data:    result,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}
