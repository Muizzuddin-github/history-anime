package test

import (
	"bytes"
	"context"
	"encoding/json"
	"history_anime/src/db"
	"history_anime/src/requestbody"
	"history_anime/src/response"
	"history_anime/src/routers"
	"history_anime/src/utility"
	"history_anime/test/dbutility"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var server *httptest.Server

var bodyRegister = requestbody.Register{
	Username: "hasan",
	Email:    "hasan@gmail.com",
	Password: "123",
}

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()
	db.CreateConnection(ctx)
	server = httptest.NewServer(routers.Router())
	m.Run()
	defer db.CloseDB(ctx)
	defer server.Close()
	err = dbutility.DeleteUser(bodyRegister.Email)
	if err != nil {
		panic(err)
	}
}

func TestLogin(t *testing.T) {

	bodyRegisterByte, _ := json.Marshal(bodyRegister)

	resp, err := http.Post(server.URL+"/api/register", "application/json", bytes.NewReader(bodyRegisterByte))
	require.Nil(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	t.Run("success", func(t *testing.T) {

		body, _ := json.Marshal(requestbody.Login{
			Email:    bodyRegister.Email,
			Password: bodyRegister.Password,
		})

		res, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewReader(body))
		require.Nil(t, err)

		bodyResult := res.Body

		bodyByte, err := io.ReadAll(bodyResult)
		require.Nil(t, err)
		defer bodyResult.Close()

		resBody := response.Login{}
		err = json.Unmarshal(bodyByte, &resBody)
		require.Nil(t, err)

		cookie := res.Cookies()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "login success", resBody.Message)
		assert.NotNil(t, resBody.Token)
		assert.Equal(t, len(cookie), 1)

	})

	t.Run("client error password", func(t *testing.T) {
		body, _ := json.Marshal(requestbody.Login{
			Email:    "haskuy12@gmail.com",
			Password: bodyRegister.Password,
		})

		res, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewReader(body))
		require.Nil(t, err)
		bodyResult := res.Body

		bodyByte, err := io.ReadAll(bodyResult)
		require.Nil(t, err)
		defer bodyResult.Close()

		resBody := response.Errors{}
		err = json.Unmarshal(bodyByte, &resBody)
		require.Nil(t, err)

		cookie := res.Cookies()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "check email or password", resBody.Errors[0])
		assert.Equal(t, len(cookie), 0)
	})

	t.Run("client error email", func(t *testing.T) {
		body, _ := json.Marshal(requestbody.Login{
			Email:    "haskuy@gmail.com",
			Password: "haskuy",
		})

		res, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewReader(body))
		require.Nil(t, err)
		bodyResult := res.Body

		bodyByte, err := io.ReadAll(bodyResult)
		require.Nil(t, err)
		defer bodyResult.Close()

		resBody := response.Errors{}
		err = json.Unmarshal(bodyByte, &resBody)
		require.Nil(t, err)

		cookie := res.Cookies()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "check email or password", resBody.Errors[0])
		assert.Equal(t, len(cookie), 0)
	})

	t.Run("error validation", func(t *testing.T) {
		body, _ := json.Marshal(requestbody.Login{
			Email:    "haskuy",
			Password: "haskuy",
		})

		res, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewReader(body))
		require.Nil(t, err)

		bodyResult := res.Body

		bodyByte, err := io.ReadAll(bodyResult)
		require.Nil(t, err)
		defer bodyResult.Close()

		resBody := response.Errors{}
		err = json.Unmarshal(bodyByte, &resBody)
		require.Nil(t, err)

		cookie := res.Cookies()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "Error:Field validation for 'Email' failed on the 'required' tag", resBody.Errors[0])
		assert.Equal(t, len(cookie), 0)
	})

}

func TestLogout(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodPost, server.URL+"/api/logout", nil)
		require.Nil(t, err)

		req.AddCookie(&http.Cookie{
			Name:     "token",
			Value:    "halo",
			Expires:  time.Now().Add(time.Hour * 24),
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		})

		client := &http.Client{}
		res, err := client.Do(req)
		require.Nil(t, err)

		resultBody := res.Body
		defer resultBody.Close()

		bodyByte, err := io.ReadAll(resultBody)
		require.Nil(t, err)

		resBody := response.Msg{}
		err = json.Unmarshal(bodyByte, &resBody)
		require.Nil(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "logout success", resBody.Message)
		assert.Equal(t, -1, res.Cookies()[0].MaxAge)
	})
}

func TestForgotPassword(t *testing.T) {
	t.Skip()
	t.Run("success", func(t *testing.T) {

		bodyByte, err := json.Marshal(requestbody.ForgotPassword{
			Email: "muizzuddin334@gmail.com",
		})

		require.Nil(t, err)

		res, err := http.Post(server.URL+"/api/forgot-password", "application/json", bytes.NewReader(bodyByte))
		require.Nil(t, err)

		result := res.Body
		defer result.Close()

		resBodyByte, err := io.ReadAll(result)
		require.Nil(t, err)

		resBody := response.Msg{}

		err = json.Unmarshal(resBodyByte, &resBody)
		require.Nil(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "email has been sent and check your email", resBody.Message)

	})

	t.Run("error validatoin", func(t *testing.T) {
		bodyByte, err := json.Marshal(requestbody.ForgotPassword{
			Email: "hasan",
		})

		require.Nil(t, err)

		res, err := http.Post(server.URL+"/api/forgot-password", "application/json", bytes.NewReader(bodyByte))
		require.Nil(t, err)

		result := res.Body
		defer result.Close()

		resBodyByte, err := io.ReadAll(result)
		require.Nil(t, err)

		resBody := response.Errors{}

		err = json.Unmarshal(resBodyByte, &resBody)
		require.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "Error:Field validation for 'Email' failed on the 'required' tag", resBody.Errors[0])

	})
}

func TestResetPassword(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		token, err := utility.CreateTokenForgotPassword(os.Getenv("SECRET_KEY"), bodyRegister.Email)
		require.Nil(t, err)

		bodyResetPassword := requestbody.ResetPassword{
			NewPassword: "1234",
			Token:       token,
		}

		bodyResetPasswordByte, err := json.Marshal(bodyResetPassword)
		require.Nil(t, err)

		res, err := http.Post(server.URL+"/api/reset-password", "application/json", bytes.NewReader(bodyResetPasswordByte))

		require.Nil(t, err)

		resBody := res.Body
		defer resBody.Close()

		resBodyByte, err := io.ReadAll(resBody)
		require.Nil(t, err)

		body := response.Msg{}
		err = json.Unmarshal(resBodyByte, &body)
		require.Nil(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "reset password success", body.Message)

		bodyLogin, _ := json.Marshal(requestbody.Login{
			Email:    "haskuy12@gmail.com",
			Password: bodyRegister.Password,
		})

		resLogin, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewReader(bodyLogin))
		require.Nil(t, err)
		bodyResult := resLogin.Body
		defer bodyResult.Close()

		bodyByte, err := io.ReadAll(bodyResult)
		require.Nil(t, err)

		resBodyLogin := response.Errors{}
		err = json.Unmarshal(bodyByte, &resBodyLogin)
		require.Nil(t, err)

		cookie := res.Cookies()
		assert.Equal(t, http.StatusBadRequest, resLogin.StatusCode)
		assert.Equal(t, "check email or password", resBodyLogin.Errors[0])
		assert.Equal(t, len(cookie), 0)
	})

	t.Run("error token invalid", func(t *testing.T) {

		id, err := dbutility.FindUser(bodyRegister.Email)
		require.Nil(t, err)

		token, err := utility.CreateTokenForgotPassword(os.Getenv("SECRET_KEY"), id)
		require.Nil(t, err)

		bodyResetPassword := requestbody.ResetPassword{
			NewPassword: "1234",
			Token:       token,
		}

		bodyResetPasswordByte, err := json.Marshal(bodyResetPassword)
		require.Nil(t, err)

		res, err := http.Post(server.URL+"/api/reset-password", "application/json", bytes.NewReader(bodyResetPasswordByte))

		require.Nil(t, err)

		resBody := res.Body
		defer resBody.Close()

		resBodyByte, err := io.ReadAll(resBody)
		require.Nil(t, err)

		body := response.Errors{}
		err = json.Unmarshal(resBodyByte, &body)
		require.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "token invalid", body.Errors[0])
	})

	t.Run("jwt error type", func(t *testing.T) {
		token, err := utility.CreateToken(os.Getenv("SECRET_KEY"), bodyRegister.Email)
		require.Nil(t, err)

		bodyResetPassword := requestbody.ResetPassword{
			NewPassword: "1234",
			Token:       token,
		}

		bodyResetPasswordByte, err := json.Marshal(bodyResetPassword)
		require.Nil(t, err)

		res, err := http.Post(server.URL+"/api/reset-password", "application/json", bytes.NewReader(bodyResetPasswordByte))

		require.Nil(t, err)

		resBody := res.Body
		defer resBody.Close()

		resBodyByte, err := io.ReadAll(resBody)
		require.Nil(t, err)

		body := response.Errors{}
		err = json.Unmarshal(resBodyByte, &body)
		require.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "type token error", body.Errors[0])
	})

}

func TestIsLogin(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		id, err := dbutility.FindUser(bodyRegister.Email)
		require.Nil(t, err)

		token, err := utility.CreateToken(os.Getenv("SECRET_KEY"), id)
		require.Nil(t, err)

		request, err := http.NewRequest(http.MethodGet, server.URL+"/api/islogin", nil)
		require.Nil(t, err)

		request.AddCookie(&http.Cookie{
			Name:     "token",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 24),
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		})

		client := &http.Client{}
		res, err := client.Do(request)
		require.Nil(t, err)

		resBody := res.Body
		defer resBody.Close()

		resBodyByte, err := io.ReadAll(resBody)
		require.Nil(t, err)

		body := response.Msg{}
		err = json.Unmarshal(resBodyByte, &body)
		require.Nil(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "user has been login", body.Message)
	})

	t.Run("user is not login", func(t *testing.T) {

		res, err := http.Get(server.URL + "/api/islogin")
		require.Nil(t, err)

		resBody := res.Body
		defer resBody.Close()

		resBodyByte, err := io.ReadAll(resBody)
		require.Nil(t, err)

		body := response.Errors{}
		err = json.Unmarshal(resBodyByte, &body)
		require.Nil(t, err)

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
		assert.Equal(t, "user is not logged in", body.Errors[0])
	})

	t.Run("jwt error type", func(t *testing.T) {
		token, err := utility.CreateToken(os.Getenv("SECRET_KEY"), bodyRegister.Email)
		require.Nil(t, err)

		request, err := http.NewRequest(http.MethodGet, server.URL+"/api/islogin", nil)
		require.Nil(t, err)

		request.AddCookie(&http.Cookie{
			Name:     "token",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 24),
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		})

		client := &http.Client{}
		res, err := client.Do(request)
		require.Nil(t, err)

		resBody := res.Body
		defer resBody.Close()

		resBodyByte, err := io.ReadAll(resBody)
		require.Nil(t, err)

		body := response.Errors{}
		err = json.Unmarshal(resBodyByte, &body)
		require.Nil(t, err)

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
		assert.Equal(t, "user is not logged in", body.Errors[0])

	})
}
