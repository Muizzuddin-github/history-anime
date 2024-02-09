package controllers

import (
	"context"
	"encoding/json"
	"fmt"

	"history_anime/src/db"
	"history_anime/src/repository"
	"history_anime/src/requestbody"
	"history_anime/src/response"
	"history_anime/src/utility"
	"history_anime/src/validation"
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

		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	body := requestbody.Register{}
	err = json.Unmarshal(result, &body)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})
		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	res, err := json.Marshal(response.Msg{
		Message: "Registrasi berhasil",
	})

	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusInternalServerError, res)
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

		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	response.SendJSONResponse(w, http.StatusCreated, res)
}

var Login httprouter.Handle = func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	result, err := io.ReadAll(r.Body)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})
		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	body := requestbody.Login{}
	err = json.Unmarshal(result, &body)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})
		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	errResult := validation.ValidateLogin(&body)
	if len(errResult) > 0 {
		res, _ := json.Marshal(response.Errors{
			Errors: errResult,
		})
		response.SendJSONResponse(w, http.StatusBadRequest, res)
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
			response.SendJSONResponse(w, http.StatusBadRequest, res)
			return
		}
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})
		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{"check email or password"},
		})
		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	token, err := utility.CreateToken(os.Getenv("SECRET_KEY"), user.Id.Hex())
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})
		response.SendJSONResponse(w, http.StatusInternalServerError, res)
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

	res, _ := json.Marshal(response.Login{
		Message: "login success",
		Token:   token,
	})

	http.SetCookie(w, &cookie)
	response.SendJSONResponse(w, http.StatusOK, res)
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
	response.SendJSONResponse(w, http.StatusOK, res)
}

var ForgotPassword httprouter.Handle = func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	body := requestbody.ForgotPassword{}
	err = json.Unmarshal(bodyByte, &body)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	errResult := validation.ValidateForgotPassword(&body)
	if len(errResult) > 0 {
		res, _ := json.Marshal(response.Errors{
			Errors: errResult,
		})

		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	token, err := utility.CreateTokenForgotPassword(os.Getenv("SECRET_KEY"), body.Email)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: errResult,
		})

		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	url := fmt.Sprintf("%s/reset-password/%s", os.Getenv("CLIENT_URL_HOST"), token)
	err = utility.SendEmail(utility.Email{
		From:    "History Anime",
		To:      body.Email,
		Subject: "reset password",
		Html: fmt.Sprintf(`
			<p style="font-weight: bold;">silahkan reset password anda <a href="%s" target="_blank"> reset password </a> </p>
		  <p style="font-weight: bold;"> link berlaku 10 menit </p>
		`, url),
	})

	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	res, err := json.Marshal(response.Msg{
		Message: "email has been sent and check your email",
	})

	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	response.SendJSONResponse(w, http.StatusInternalServerError, res)
}

var ResetPassword httprouter.Handle = func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	body := requestbody.ResetPassword{}
	err = json.Unmarshal(bodyByte, &body)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	errResult := validation.ValidateResetPassword(&body)
	if len(errResult) > 0 {
		res, _ := json.Marshal(response.Errors{
			Errors: errResult,
		})

		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	email, err := utility.VerifyTokenForgotPassword(os.Getenv("SECRET_KEY"), body.Token)

	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	newHashPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 10)
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	ctx := context.Background()
	auth := repository.AuthRepo(db.DB)
	result, err := auth.ResetPassword(ctx, email, string(newHashPassword))
	if err != nil {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{err.Error()},
		})

		response.SendJSONResponse(w, http.StatusInternalServerError, res)
		return
	}

	if result.MatchedCount == 0 {
		res, _ := json.Marshal(response.Errors{
			Errors: []string{"token invalid"},
		})

		response.SendJSONResponse(w, http.StatusBadRequest, res)
		return
	}

	res, err := json.Marshal(response.Msg{
		Message: "reset password success",
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
