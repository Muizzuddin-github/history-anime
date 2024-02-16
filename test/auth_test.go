package test

import (
	"bytes"
	"encoding/json"
	"history_anime/src/requestbody"
	"history_anime/src/response"
	"history_anime/src/utility"
	"history_anime/test/dbutility"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {

	var loginTesting = requestbody.Register{
		Username: "testing",
		Email:    "testing@gmail.com",
		Password: "123",
	}

	loginTestingByte, _ := json.Marshal(loginTesting)

	resp, err := http.Post(Server.URL+"/api/register", "application/json", bytes.NewReader(loginTestingByte))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	t.Run("success", func(t *testing.T) {

		body, _ := json.Marshal(requestbody.Login{
			Email:    loginTesting.Email,
			Password: loginTesting.Password,
		})

		res, err := http.Post(Server.URL+"/api/login", "application/json", bytes.NewReader(body))
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
			Email:    loginTesting.Email,
			Password: "1234",
		})

		res, err := http.Post(Server.URL+"/api/login", "application/json", bytes.NewReader(body))
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

	t.Run("client error not found email", func(t *testing.T) {
		body, _ := json.Marshal(requestbody.Login{
			Email:    "haskuy@gmail.com",
			Password: loginTesting.Password,
		})

		res, err := http.Post(Server.URL+"/api/login", "application/json", bytes.NewReader(body))
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

		res, err := http.Post(Server.URL+"/api/login", "application/json", bytes.NewReader(body))
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

	err = dbutility.DeleteUser(loginTesting.Email)
	require.Nil(t, err)

}

func TestLogout(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodPost, Server.URL+"/api/logout", nil)
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
			Email: os.Getenv("MY_EMAIL"),
		})
		require.Nil(t, err)

		res, err := http.Post(Server.URL+"/api/forgot-password", "application/json", bytes.NewReader(bodyByte))
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

		res, err := http.Post(Server.URL+"/api/forgot-password", "application/json", bytes.NewReader(bodyByte))
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

		var resetTesting = requestbody.Register{
			Username: "testingreset",
			Email:    "testing@gmail.com",
			Password: "123",
		}

		loginTestingByte, _ := json.Marshal(resetTesting)

		resp, err := http.Post(Server.URL+"/api/register", "application/json", bytes.NewReader(loginTestingByte))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		token, err := utility.CreateTokenForgotPassword(os.Getenv("SECRET_KEY"), resetTesting.Email)
		require.Nil(t, err)

		bodyResetPassword := requestbody.ResetPassword{
			NewPassword: "1234",
			Token:       token,
		}

		bodyResetPasswordByte, err := json.Marshal(bodyResetPassword)
		require.Nil(t, err)

		res, err := http.Post(Server.URL+"/api/reset-password", "application/json", bytes.NewReader(bodyResetPasswordByte))

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
			Email:    resetTesting.Email,
			Password: resetTesting.Password,
		})

		resLogin, err := http.Post(Server.URL+"/api/login", "application/json", bytes.NewReader(bodyLogin))
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

		err = dbutility.DeleteUser(resetTesting.Email)
		require.Nil(t, err)
	})

	t.Run("error token invalid", func(t *testing.T) {

		var resetTesting = requestbody.Register{
			Username: "testingtokeninvalid",
			Email:    "testing@gmail.com",
			Password: "123",
		}

		loginTestingByte, _ := json.Marshal(resetTesting)

		resp, err := http.Post(Server.URL+"/api/register", "application/json", bytes.NewReader(loginTestingByte))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		id, err := dbutility.FindUser(resetTesting.Email)
		require.Nil(t, err)

		token, err := utility.CreateTokenForgotPassword(os.Getenv("SECRET_KEY"), id)
		require.Nil(t, err)

		bodyResetPassword := requestbody.ResetPassword{
			NewPassword: "1234",
			Token:       token,
		}

		bodyResetPasswordByte, err := json.Marshal(bodyResetPassword)
		require.Nil(t, err)

		res, err := http.Post(Server.URL+"/api/reset-password", "application/json", bytes.NewReader(bodyResetPasswordByte))

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

		err = dbutility.DeleteUser(resetTesting.Email)
		require.Nil(t, err)
	})

	t.Run("jwt error type", func(t *testing.T) {
		token, err := utility.CreateToken(os.Getenv("SECRET_KEY"), "email@gmail.com")
		require.Nil(t, err)

		bodyResetPassword := requestbody.ResetPassword{
			NewPassword: "1234",
			Token:       token,
		}

		bodyResetPasswordByte, err := json.Marshal(bodyResetPassword)
		require.Nil(t, err)

		res, err := http.Post(Server.URL+"/api/reset-password", "application/json", bytes.NewReader(bodyResetPasswordByte))

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
		var isLoginTesting = requestbody.Register{
			Username: "testingtokeninvalid",
			Email:    "testing@gmail.com",
			Password: "123",
		}

		loginTestingByte, _ := json.Marshal(isLoginTesting)

		resp, err := http.Post(Server.URL+"/api/register", "application/json", bytes.NewReader(loginTestingByte))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		id, err := dbutility.FindUser(isLoginTesting.Email)
		require.Nil(t, err)

		token, err := utility.CreateToken(os.Getenv("SECRET_KEY"), id)
		require.Nil(t, err)

		request, err := http.NewRequest(http.MethodGet, Server.URL+"/api/islogin", nil)
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

		err = dbutility.DeleteUser(isLoginTesting.Email)
		require.Nil(t, err)
	})

	t.Run("user is not login", func(t *testing.T) {

		res, err := http.Get(Server.URL + "/api/islogin")
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

		token, err := utility.CreateToken(os.Getenv("SECRET_KEY"), "email@gmail.com")
		require.Nil(t, err)

		request, err := http.NewRequest(http.MethodGet, Server.URL+"/api/islogin", nil)
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
