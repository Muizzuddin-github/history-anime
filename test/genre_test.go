package test

import (
	"bytes"
	"encoding/json"
	"history_anime/src/requestbody"
	"history_anime/src/response"
	"history_anime/test/dbutility"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenreAdd(t *testing.T) {

	t.Run("success", func(t *testing.T) {

		bodyInsert := requestbody.Genre{
			Name: "contoh",
		}

		bodyInsertByte, err := json.Marshal(bodyInsert)
		require.Nil(t, err)

		request, err := http.NewRequest(http.MethodPost, Server.URL+"/api/genre", bytes.NewReader(bodyInsertByte))
		require.Nil(t, err)

		request.Header.Set("Content-Type", "application/json")
		request.AddCookie(&http.Cookie{
			Name:     "token",
			Value:    TokenUser,
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

		resBodyJson := response.GenreInsert{}
		err = json.Unmarshal(resBodyByte, &resBodyJson)
		require.Nil(t, err)

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "insert genre success", resBodyJson.Message)
		assert.NotEmpty(t, resBodyJson.InsertedID)

		err = dbutility.GenreDeleteById(resBodyJson.InsertedID)
		require.Nil(t, err)
	})

	t.Run("content type error", func(t *testing.T) {

		bodyInsert := requestbody.Genre{
			Name: "contoh",
		}

		bodyInsertByte, err := json.Marshal(bodyInsert)
		require.Nil(t, err)

		request, err := http.NewRequest(http.MethodPost, Server.URL+"/api/genre", bytes.NewReader(bodyInsertByte))
		require.Nil(t, err)

		request.AddCookie(&http.Cookie{
			Name:     "token",
			Value:    TokenUser,
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

		resBodyJson := response.Errors{}
		err = json.Unmarshal(resBodyByte, &resBodyJson)
		require.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "content-type must be application/json", resBodyJson.Errors[0])
	})
}

func TestGenreDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		id, err := dbutility.GenreAdd(&requestbody.Genre{
			Name: "contohkuy",
		})
		require.Nil(t, err)

		request, err := http.NewRequest(http.MethodDelete, Server.URL+"/api/genre/"+id, nil)
		require.Nil(t, err)

		request.AddCookie(&http.Cookie{
			Name:     "token",
			Value:    TokenUser,
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

		resBodyJson := response.Msg{}
		err = json.Unmarshal(resBodyByte, &resBodyJson)
		require.Nil(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "delete genre success", resBodyJson.Message)

	})

	t.Run("not found", func(t *testing.T) {

		request, err := http.NewRequest(http.MethodDelete, Server.URL+"/api/genre/65c5c8634a978c8f77b310b2", nil)
		require.Nil(t, err)

		request.AddCookie(&http.Cookie{
			Name:     "token",
			Value:    TokenUser,
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

		resBodyJson := response.Errors{}
		err = json.Unmarshal(resBodyByte, &resBodyJson)
		require.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		assert.Equal(t, "genre not found", resBodyJson.Errors[0])
	})
}
