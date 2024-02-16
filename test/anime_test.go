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

func TestAnimeAdd(t *testing.T) {

	t.Run("succcess", func(t *testing.T) {
		dataInsert := requestbody.Anime{
			Name:        "testing",
			Description: "lorem",
			Genre:       []string{"testing"},
			Image:       "https://example.com",
			Status:      "watching",
		}

		bodyByte, err := json.Marshal(dataInsert)
		require.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, Server.URL+"/api/anime", bytes.NewReader(bodyByte))
		require.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:     "token",
			Value:    TokenUser,
			Expires:  time.Now().Add(time.Hour * 24),
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		})

		client := &http.Client{}
		res, err := client.Do(req)
		require.Nil(t, err)

		body := res.Body
		defer body.Close()

		resBodyByte, err := io.ReadAll(body)
		require.Nil(t, err)

		resBody := response.Msg{}

		err = json.Unmarshal(resBodyByte, &resBody)
		require.Nil(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "insert anime success", resBody.Message)

		err = dbutility.AnimeDeleteOne(dataInsert.Name)
		require.Nil(t, err)
	})

	t.Run("validation error required", func(t *testing.T) {
		data := requestbody.Anime{
			Name:        "",
			Description: "",
			Genre:       []string{},
			Image:       "",
			Status:      "",
		}

		bodyByte, err := json.Marshal(data)
		require.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, Server.URL+"/api/anime", bytes.NewReader(bodyByte))
		require.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:     "token",
			Value:    TokenUser,
			Expires:  time.Now().Add(time.Hour * 24),
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		})

		client := &http.Client{}
		res, err := client.Do(req)
		require.Nil(t, err)

		body := res.Body
		defer body.Close()

		resBodyByte, err := io.ReadAll(body)
		require.Nil(t, err)

		resBody := response.Errors{}
		err = json.Unmarshal(resBodyByte, &resBody)
		require.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, 5, len(resBody.Errors))

		err = dbutility.AnimeDeleteOne(data.Name)
		require.NotNil(t, err)
	})
}

func TestAnimeUpdate(t *testing.T) {
	data := requestbody.Anime{
		Name:        "testing",
		Genre:       []string{"testing"},
		Description: "testingubah",
		Image:       "https://history.com",
		Status:      "finish",
	}

	id, err := dbutility.AnimeAdd(data.Name, data.Image, data.Genre, data.Description, data.Status)
	require.Nil(t, err)

	t.Run("success", func(t *testing.T) {

		data := requestbody.Anime{
			Name:        "testingubah",
			Genre:       []string{"testingubah"},
			Description: "testingubah",
			Image:       "https://history.com",
			Status:      "finish",
		}

		bodyByte, err := json.Marshal(data)
		require.Nil(t, err)

		request, err := http.NewRequest(http.MethodPut, Server.URL+"/api/anime/"+id, bytes.NewReader(bodyByte))
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

		resBodyRead, err := io.ReadAll(resBody)
		require.Nil(t, err)

		resBodyJson := response.Msg{}
		err = json.Unmarshal(resBodyRead, &resBodyJson)
		require.Nil(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "update anime success", resBodyJson.Message)
	})

	t.Run("validation error required", func(t *testing.T) {

		data := requestbody.Anime{
			Name:        "",
			Genre:       []string{},
			Description: "",
			Image:       "",
			Status:      "",
		}

		bodyByte, err := json.Marshal(data)
		require.Nil(t, err)

		request, err := http.NewRequest(http.MethodPut, Server.URL+"/api/anime/"+id, bytes.NewReader(bodyByte))
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

		resBodyRead, err := io.ReadAll(resBody)
		require.Nil(t, err)

		resBodyJson := response.Errors{}
		err = json.Unmarshal(resBodyRead, &resBodyJson)
		require.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, 5, len(resBodyJson.Errors))
	})

	err = dbutility.AnimeDeleteOneById(id)
	require.Nil(t, err)
}

func TestAnimeDelete(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		dataInsert := requestbody.Anime{
			Name:        "testing",
			Description: "lorem",
			Genre:       []string{"testing"},
			Image:       "https://example.com",
			Status:      "watching",
		}
		id, err := dbutility.AnimeAdd(dataInsert.Name, dataInsert.Image, dataInsert.Genre, dataInsert.Description, dataInsert.Status)
		require.Nil(t, err)

		request, err := http.NewRequest(http.MethodDelete, Server.URL+"/api/anime/"+id, nil)
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

		body := res.Body
		defer body.Close()

		bodyByte, err := io.ReadAll(body)
		require.Nil(t, err)

		bodyJson := response.Msg{}

		err = json.Unmarshal(bodyByte, &bodyJson)
		require.Nil(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "delete anime success", bodyJson.Message)
	})

	t.Run("error not found", func(t *testing.T) {

		request, err := http.NewRequest(http.MethodDelete, Server.URL+"/api/anime/65c217fc6556430b3dc4ce61", nil)
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

		body := res.Body
		defer body.Close()

		bodyByte, err := io.ReadAll(body)
		require.Nil(t, err)

		bodyJson := response.Errors{}

		err = json.Unmarshal(bodyByte, &bodyJson)
		require.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		assert.Equal(t, "anime not found", bodyJson.Errors[0])
	})

}

func TestAnimeGetAllSuccess(t *testing.T) {

	dataInsert := requestbody.Anime{
		Name:        "testing1",
		Description: "lorem",
		Genre:       []string{"testing1"},
		Image:       "https://example.com",
		Status:      "watching",
	}

	id, err := dbutility.AnimeAdd(dataInsert.Name, dataInsert.Image, dataInsert.Genre, dataInsert.Description, dataInsert.Status)
	require.Nil(t, err)

	request, err := http.NewRequest(http.MethodGet, Server.URL+"/api/anime", nil)
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

	body := res.Body
	defer body.Close()

	bodyByte, err := io.ReadAll(body)
	require.Nil(t, err)

	bodyJson := response.AnimeAll{}
	err = json.Unmarshal(bodyByte, &bodyJson)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "all data anime", bodyJson.Message)

	err = dbutility.AnimeDeleteOneById(id)
	require.Nil(t, err)

}
