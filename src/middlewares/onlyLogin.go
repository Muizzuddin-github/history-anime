package middlewares

import (
	"context"
	"encoding/json"
	"history_anime/src/db"
	"history_anime/src/entity"
	"history_anime/src/logger"
	"history_anime/src/response"
	"history_anime/src/utility"

	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func OnlyLogin(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

		token, err := r.Cookie("token")
		if err == http.ErrNoCookie {

			res, _ := json.Marshal(response.Errors{
				Errors: []string{"Unauthorized"},
			})

			logger.New().WithFields(logrus.Fields{
				"action": "No Token",
				"status": http.StatusText(http.StatusUnauthorized),
			}).Warn(err.Error())
			response.SendJSONResponse(w, http.StatusUnauthorized, res)
			return
		}

		id, err := utility.VerifyToken(os.Getenv("SECRET_KEY"), token.Value)
		if err != nil {
			res, _ := json.Marshal(response.Errors{
				Errors: []string{"Unauthorized"},
			})

			logger.New().WithFields(logrus.Fields{
				"action": "Token Invalid",
				"status": http.StatusText(http.StatusUnauthorized),
			}).Warn(err.Error())
			response.SendJSONResponse(w, http.StatusUnauthorized, res)
			return
		}

		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			res, _ := json.Marshal(response.Errors{
				Errors: []string{"Unauthorized"},
			})

			logger.New().WithFields(logrus.Fields{
				"action": "ObjectID Mongodb Invalid",
				"status": http.StatusText(http.StatusUnauthorized),
			}).Warn(err.Error())
			response.SendJSONResponse(w, http.StatusUnauthorized, res)
			return
		}

		ctx := context.Background()
		filter := bson.D{{Key: "_id", Value: objId}}
		result := entity.Users{}
		err = db.DB.Collection("users").FindOne(ctx, filter).Decode(&result)
		if err == mongo.ErrNoDocuments {
			res, _ := json.Marshal(response.Errors{
				Errors: []string{"Unauthorized"},
			})

			logger.New().WithFields(logrus.Fields{
				"action": "User Not Found",
				"status": http.StatusText(http.StatusUnauthorized),
			}).Warn(err.Error())
			response.SendJSONResponse(w, http.StatusUnauthorized, res)
			return
		} else if err != nil {
			res, _ := json.Marshal(response.Errors{
				Errors: []string{err.Error()},
			})
			logger.New().WithFields(logrus.Fields{
				"action": "database error",
				"status": http.StatusText(http.StatusUnauthorized),
			}).Error(err.Error())
			response.SendJSONResponse(w, http.StatusInternalServerError, res)
			return
		}

		next(w, r, params)
	}
}
