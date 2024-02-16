package test

import (
	"bytes"
	"context"
	"encoding/json"
	"history_anime/src/db"
	"history_anime/src/requestbody"
	"history_anime/src/routers"
	"history_anime/test/dbutility"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
)

var Server *httptest.Server
var TokenUser string

func setupUser() {
	var bodyRegister = requestbody.Register{
		Username: "hasan",
		Email:    "hasan@gmail.com",
		Password: "123",
	}
	bodyRegisterByte, _ := json.Marshal(bodyRegister)

	resp, err := http.Post(Server.URL+"/api/register", "application/json", bytes.NewReader(bodyRegisterByte))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := json.Marshal(requestbody.Login{
		Email:    bodyRegister.Email,
		Password: bodyRegister.Password,
	})

	res, err := http.Post(Server.URL+"/api/login", "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}

	cookie := res.Cookies()[0]
	TokenUser = cookie.Value
}

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()
	db.CreateConnection(ctx)
	Server = httptest.NewServer(routers.Router())
	setupUser()
	m.Run()
	defer db.CloseDB(ctx)
	defer Server.Close()

	err = dbutility.DeleteUser("hasan@gmail.com")
	if err != nil {
		panic(err)
	}
}
