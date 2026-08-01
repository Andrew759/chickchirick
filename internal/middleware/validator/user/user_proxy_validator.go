package user

import (
	"bytes"
	mainService "chickChirick/cmd/service"
	"chickChirick/internal/controller/c_http"
	"chickChirick/internal/middleware"
	"chickChirick/internal/middleware/config"
	"chickChirick/internal/model/user"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type UserProxyValidatorFactory struct{}

func (uvf UserProxyValidatorFactory) NewValidator(dbDecorator mainService.DBDecorator, opts ...middleware.ValidatorOption) middleware.Validator {
	var vOptions middleware.ValidatorOptions
	for _, opt := range opts {
		opt(&vOptions)
	}

	return UserProxyValidator{
		DBDecorator:      dbDecorator,
		ValidatorOptions: vOptions,
	}
}

type UserProxyValidator struct {
	DBDecorator mainService.DBDecorator
	middleware.ValidatorOptions
}

func (upv UserProxyValidator) Validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			c_http.NewResponse().SendError(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		r.Body.Close()

		var u user.User
		if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&u); err != nil {
			c_http.NewResponse().SendError(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if errorList := validateRequestRules(u); len(errorList) > 0 {
			errResponse := c_http.NewResponse()
			errResponse.AddErrorsToErrorContainer(errorList)

			errResponse.Send(w, http.StatusBadRequest)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		ctx := context.WithValue(r.Context(), config.UserUserKey, &u)
		next(w, r.WithContext(ctx))
	}
}
