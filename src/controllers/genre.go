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

var GenreGetAll httprouter.Handle = func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	ctx := context.Background()
	genreCol := repository.GenreRepo(db.DB)
	result, err := genreCol.GetAll(ctx)
	if err != nil {

		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	res, err := json.Marshal(response.GenreAll{
		Message: "data genre",
		Data:    result,
	})
	if err != nil {
		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	response.SendJSONResponse(w, http.StatusOK, res)
}

var GenreAdd httprouter.Handle = func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

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

	body := requestbody.Genre{}
	err = json.Unmarshal(bodyByte, &body)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})
		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	invalid := validation.ValidateGenre(&body)
	if len(invalid) > 0 {
		res, _ := json.Marshal(response.Errors{
			Errors: invalid,
		})
		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	ctx := context.Background()
	genreCol := repository.GenreRepo(db.DB)

	insertedID, err := genreCol.Add(ctx, &body)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: invalid,
		})
		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	res, err := json.Marshal(response.GenreInsert{
		Message:    "insert genre success",
		InsertedID: insertedID,
	})

	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})
		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	response.SendJSONResponse(w, http.StatusCreated, res)
}

var GenreDelete httprouter.Handle = func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	id := params.ByName("id")

	ctx := context.Background()
	genreCol := repository.GenreRepo(db.DB)
	result, err := genreCol.Del(ctx, id)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	if result.DeletedCount == 0 {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{"genre not found"},
		})

		response.SendJSONResponse(w, http.StatusNotFound, res)
		return
	}

	res, err := json.Marshal(response.Msg{
		Message: "delete genre success",
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
