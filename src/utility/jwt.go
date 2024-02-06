package utility

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(key string, id string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", errors.New(err.Error())
	}

	return tokenString, nil
}

func VerifyToken(key string, token string) (string, error) {
	tokenv, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return "", errors.New(err.Error())
	}

	if !tokenv.Valid {
		return "", errors.New(err.Error())
	}

	id := tokenv.Claims.(jwt.MapClaims)["id"].(string)

	return id, nil
}
