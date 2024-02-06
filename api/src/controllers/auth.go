package controllers

import (
	"context"
	"encoding/json"
	"history_anime/api/src/db"
	"history_anime/api/src/repository"
	"history_anime/api/src/requestbody"
	"history_anime/api/src/response"
	"history_anime/api/src/utility"
	"history_anime/api/src/validation"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var Register httprouter.Handle = func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	result, err := io.ReadAll(r.Body)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	body := requestbody.Register{}
	err = json.Unmarshal(result, &body)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	res, err := json.Marshal(response.Msg{
		Message: "Registrasi berhasil",
	})

	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	body.Password = string(hash)

	ctx := context.Background()
	con := repository.AuthRepo(db.DB)
	err = con.Register(ctx, &body)
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

var Login httprouter.Handle = func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	result, err := io.ReadAll(r.Body)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	body := requestbody.Login{}
	err = json.Unmarshal(result, &body)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	errResult := validation.ValidateLogin(&body)
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
	con := repository.AuthRepo(db.DB)
	user, err := con.Login(ctx, &body)

	if err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			res, _ := json.Marshal(response.Errors{
				Errors: []string{"check email or password"},
			})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(res)
			return
		}
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{"check email or password"},
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	token, err := utility.CreateToken(os.Getenv("SECRET_KEY"), user.Id.Hex())
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	cookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	res, _ := json.Marshal(response.Msg{
		Message: "login success",
	})

	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

var Logout httprouter.Handle = func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	cookie := http.Cookie{
		Name:     "token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	res, _ := json.Marshal(response.Msg{Message: "logout success"})
	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(res)

}
