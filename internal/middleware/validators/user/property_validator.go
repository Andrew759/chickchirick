package user

import (
	mainService "chickChirick/cmd/service"
	"chickChirick/internal/controller/c_http"
	"chickChirick/internal/middleware"
	"chickChirick/internal/middleware/config"
	"chickChirick/internal/middleware/service"
	"chickChirick/internal/model/user"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type PropertyValidatorFactory struct{}

func (pvf PropertyValidatorFactory) NewValidator(dbDecorator mainService.DBDecorator, opts ...middleware.ValidatorOption) middleware.Validator {
	var vOptions middleware.ValidatorOptions
	for _, opt := range opts {
		opt(&vOptions)
	}

	return &PropertyValidator{
		DBDecorator:      dbDecorator,
		ValidatorOptions: vOptions,
	}
}

type PropertyValidator struct {
	DBDecorator mainService.DBDecorator
	middleware.ValidatorOptions
}

func (pv PropertyValidator) Validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p user.Property

		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			c_http.NewResponse().SendError(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if errorList := pv.validateRequestRules(p); len(errorList) > 0 {
			errResponse := c_http.NewResponse()
			errResponse.AddErrorsToErrorContainer(errorList)

			errResponse.Send(w, http.StatusBadRequest)
			return
		}

		if pv.IsDBValidationActivated() {
			pv.validateAndSendResponseByDBRules(w, p)
		}

		ctx := context.WithValue(r.Context(), config.UserPropertyKey, &p)
		next(w, r.WithContext(ctx))
	}
}

func (pv PropertyValidator) validateRequestRules(p user.Property) []error {
	var errList []error

	if strings.TrimSpace(p.Email) != "" && !service.IsEmail(p.Email) {
		errList = append(errList, errors.New("invalid email"))
	}
	if p.Password != nil && !service.IsHasCorrectLength(*p.Password, 1024) {
		errList = append(errList, errors.New("invalid password"))
	}
	//TODO: валидация таймзон, после того, как появится ENUM

	return errList
}

func (pv PropertyValidator) validateAndSendResponseByDBRules(w http.ResponseWriter, p user.Property) {
	_, err := user.GetUserById(pv.DBDecorator.GormInterface, p.UserId)
	if err != nil && errors.Is(err, user.UserNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
	}
}
