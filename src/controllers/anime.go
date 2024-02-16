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

	if r.Header.Get("Content-Type") != "application/json" {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{"content-type must be application/json"},
		})
		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return

	}

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {

		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	body := requestbody.Anime{}
	err = json.Unmarshal(bodyByte, &body)
	if err != nil {

		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	errResult := validation.ValidateAnime(&body)
	if len(errResult) > 0 {
		res, _ := json.Marshal(response.Errors{
			Errors: errResult,
		})

		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	ctx := context.Background()
	anime := repository.AnimeRepo(db.DB)
	insertID, err := anime.Add(ctx, &body)
	if err != nil {

		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	res, err := json.Marshal(response.AnimeInsert{
		Message:    "insert anime success",
		InsertedID: insertID,
	})

	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	response.SendJSONResponse(w, http.StatusOK, res)
}

var AnimeUpdate httprouter.Handle = func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	if r.Header.Get("Content-Type") != "application/json" {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{"content-type must be application/json"},
		})
		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return

	}

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {

		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	body := requestbody.Anime{}
	err = json.Unmarshal(bodyByte, &body)
	if err != nil {

		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	errResult := validation.ValidateAnime(&body)
	if len(errResult) > 0 {
		res, _ := json.Marshal(response.Errors{
			Errors: errResult,
		})

		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	ctx := context.Background()
	anime := repository.AnimeRepo(db.DB)

	result, err := anime.Update(ctx, &body, params.ByName("id"))
	if err != nil {

		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})
		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	if result.MatchedCount == 0 {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{"anime not found"},
		})

		response.SendJSONResponse(w, http.StatusNotFound, res)
		return
	}

	res, err := json.Marshal(response.Msg{Message: "update anime success"})
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	response.SendJSONResponse(w, http.StatusOK, res)
}

var AnimeDel httprouter.Handle = func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	ctx := context.Background()
	anime := repository.AnimeRepo(db.DB)
	result, err := anime.Del(ctx, params.ByName("id"))
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	if result.DeletedCount == 0 {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{"anime not found"},
		})

		response.SendJSONResponse(w, http.StatusNotFound, res)
		return
	}

	res, _ := json.Marshal(response.Msg{
		Message: "delete anime success",
	})

	response.SendJSONResponse(w, http.StatusOK, res)
}

var AnimeGetAll httprouter.Handle = func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	ctx := context.Background()
	anime := repository.AnimeRepo(db.DB)
	result, err := anime.GetAll(ctx)

	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	res, _ := json.Marshal(response.AnimeAll{
		Message: "all data anime",
		Data:    result,
	})

	response.SendJSONResponse(w, http.StatusOK, res)
}
