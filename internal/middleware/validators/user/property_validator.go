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

type PropertyValidator struct{}

func (pv PropertyValidator) Validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p user.Property

		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http_transaction.NewResponse().SendError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
			return
		}

		if errorList := validateProperty(p); len(errorList) > 0 {
			errResponse := http_transaction.NewResponse()
			errResponse.AddErrorsToErrorContainer(errorList)

			errResponse.Send(w, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), config.UserPropertyKey, &p)
		next(w, r.WithContext(ctx))
	}
}

func validateProperty(p user.Property) []error {
	var errList []error

	if strings.TrimSpace(p.Email) != "" && !service.IsEmail(p.Email) {
		errList = append(errList, errors.New("invalid email"))
	}
	if p.Password != nil || !service.IsHasCorrectLength(*p.Password, 1024) {
		errList = append(errList, errors.New("invalid password"))
	}
	//TODO: валидация таймзон

	return errList
}
