package user

import (
	"chickChirick/internal/controller/http_transaction"
	"chickChirick/internal/middleware/config"
	"chickChirick/internal/middleware/service"
	"chickChirick/internal/model/user"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type UserValidator struct{}

func (uv UserValidator) Validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u user.User

		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http_transaction.NewResponse().SendError(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		if errorList := validateUser(u); len(errorList) > 0 {
			errResponse := http_transaction.NewResponse()
			errResponse.AddErrorsToErrorContainer(errorList)

			errResponse.Send(w, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), config.UserUserKey, &u)
		next(w, r.WithContext(ctx))
	}
}

func validateUser(u user.User) []error {
	var errList []error

	if strings.TrimSpace(u.Name) == "" || !service.IsHasCorrectLength(u.Name, 256) {
		errList = append(errList, errors.New("invalid name"))
	}
	if strings.TrimSpace(u.Surname) == "" || !service.IsHasCorrectLength(u.Surname, 256) {
		errList = append(errList, errors.New("invalid surname"))
	}
	if strings.TrimSpace(u.Login) == "" || !service.IsLogin(u.Login) || !service.IsHasCorrectLength(u.Login, 256) {
		errList = append(errList, errors.New("invalid login"))
	}
	if strings.TrimSpace(u.Phone) == "" || !service.IsPhoneNumber(u.Phone) {
		errList = append(errList, errors.New("invalid phone"))
	}

	return errList
}
