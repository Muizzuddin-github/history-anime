package controllers

import (
	"context"
	"encoding/json"
	"history_anime/src/db"
	"history_anime/src/logger"
	"history_anime/src/repository"
	"history_anime/src/requestbody"
	"history_anime/src/response"
	"history_anime/src/validation"
	"io"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

var GenreGetAll httprouter.Handle = func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	ctx := context.Background()
	genreCol := repository.GenreRepo(db.DB)
	result, err := genreCol.GetAll(ctx)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		logger.New().WithFields(logrus.Fields{
			"action": "Database Error",
			"status": http.StatusText(http.StatusInternalServerError),
			"path":   r.URL.Path,
			"method": r.Method,
		}).Error(err.Error())
		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	res, err := json.Marshal(response.GenreAll{
		Message: "all data genre",
		Data:    result,
	})
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		logger.New().WithFields(logrus.Fields{
			"action": "Error json.Marshal",
			"status": http.StatusText(http.StatusInternalServerError),
			"path":   r.URL.Path,
			"method": r.Method,
		}).Error(err.Error())
		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	logger.New().WithFields(logrus.Fields{
		"action": "Success",
		"status": http.StatusText(http.StatusOK),
		"path":   r.URL.Path,
		"method": r.Method,
	}).Info("Request Success")
	response.SendJSONResponse(w, http.StatusOK, res)
}

var GenreAdd httprouter.Handle = func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	if r.Header.Get("Content-Type") != "application/json" {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{"content-type must be application/json"},
		})

		logger.New().WithFields(logrus.Fields{
			"action": "Content Type",
			"status": http.StatusText(http.StatusBadRequest),
			"path":   r.URL.Path,
			"method": r.Method,
		}).Warn("Content Type Not Allowed")
		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		logger.New().WithFields(logrus.Fields{
			"action": "Error io.ReadAll",
			"status": http.StatusText(http.StatusBadRequest),
			"path":   r.URL.Path,
			"method": r.Method,
		}).Warn(err.Error())
		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	body := requestbody.Genre{}
	err = json.Unmarshal(bodyByte, &body)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		logger.New().WithFields(logrus.Fields{
			"action": "Error json.Unmarshal",
			"status": http.StatusText(http.StatusBadRequest),
			"path":   r.URL.Path,
			"method": r.Method,
		}).Warn(err.Error())
		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	errResult := validation.ValidateGenre(&body)
	if len(errResult) > 0 {
		res, _ := json.Marshal(response.Errors{
			Errors: errResult,
		})

		logger.New().WithFields(logrus.Fields{
			"action": "Validation error",
			"status": http.StatusText(http.StatusBadRequest),
			"path":   r.URL.Path,
			"method": r.Method,
		}).Warn(strings.Join(errResult, " "))
		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	ctx := context.Background()
	genreCol := repository.GenreRepo(db.DB)

	insertedID, err := genreCol.Add(ctx, &body)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		logger.New().WithFields(logrus.Fields{
			"action": "Database Error",
			"status": http.StatusText(http.StatusInternalServerError),
			"path":   r.URL.Path,
			"method": r.Method,
		}).Error(err.Error())
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

		logger.New().WithFields(logrus.Fields{
			"action": "Error json.Marshal",
			"status": http.StatusText(http.StatusInternalServerError),
			"path":   r.URL.Path,
			"method": r.Method,
		}).Error(err.Error())
		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	logger.New().WithFields(logrus.Fields{
		"action": "Success",
		"status": http.StatusText(http.StatusCreated),
		"path":   r.URL.Path,
		"method": r.Method,
	}).Info(err.Error())
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

		logger.New().WithFields(logrus.Fields{
			"action": "Database Error",
			"status": http.StatusText(http.StatusInternalServerError),
			"path":   r.URL.Path,
			"method": r.Method,
		}).Error(err.Error())
		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	if result.DeletedCount == 0 {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{"genre not found"},
		})

		logger.New().WithFields(logrus.Fields{
			"action": "Genre Not Found",
			"status": http.StatusText(http.StatusNotFound),
			"path":   r.URL.Path,
			"method": r.Method,
		}).Warn(err.Error())
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

		logger.New().WithFields(logrus.Fields{
			"action": "Error json.Marshal",
			"status": http.StatusText(http.StatusInternalServerError),
			"path":   r.URL.Path,
			"method": r.Method,
		}).Error(err.Error())
		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	logger.New().WithFields(logrus.Fields{
		"action": "Success",
		"status": http.StatusText(http.StatusOK),
		"path":   r.URL.Path,
		"method": r.Method,
	}).Info(err.Error())
	response.SendJSONResponse(w, http.StatusOK, res)
}
