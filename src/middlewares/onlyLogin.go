package middlewares

import (
	"context"
	"encoding/json"
	"history_anime/src/db"
	"history_anime/src/entity"
	"history_anime/src/response"
	"history_anime/src/utility"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
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

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(res)
			return
		}

		id, err := utility.VerifyToken(os.Getenv("SECRET_KEY"), token.Value)
		if err != nil {
			res, _ := json.Marshal(response.Errors{
				Errors: []string{"Unauthorized"},
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(res)
			return
		}

		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			res, _ := json.Marshal(response.Errors{
				Errors: []string{"Unauthorized"},
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(res)
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

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(res)
			return
		}

		next(w, r, params)
	}
}
