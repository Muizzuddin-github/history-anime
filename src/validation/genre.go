package validation

import (
	"fmt"
	"history_anime/src/requestbody"

	"github.com/go-playground/validator/v10"
)

func ValidateGenre(body *requestbody.Genre) []string {

	validate := validator.New(validator.WithRequiredStructEnabled())

	errResult := []string{}
	err := validate.Struct(body)

	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			errResult = append(errResult, fmt.Sprintf("Error:Field validation for '%s' failed on the 'required' tag", err.Field()))
		}
	}

	return errResult

}
