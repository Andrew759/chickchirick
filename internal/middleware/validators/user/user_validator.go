package user

import (
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
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := validateUser(u); err != nil {
			http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), config.UserUserKey, &u)
		next(w, r.WithContext(ctx))
	}
}

func validateUser(u user.User) error {
	if strings.TrimSpace(u.Name) == "" || !service.IsHasCorrectLength(u.Name, 256) {
		return errors.New("name is required")
	}
	if strings.TrimSpace(u.Surname) == "" || service.IsHasCorrectLength(u.Surname, 256) {
		return errors.New("surname is required")
	}
	if strings.TrimSpace(u.Login) == "" || !service.IsLatinSymbolOnly(u.Login) || service.IsHasCorrectLength(u.Login, 256) {
		return errors.New("invalid login")
	}
	if strings.TrimSpace(u.Phone) == "" || !service.IsPhoneNumber(u.Phone) {
		return errors.New("invalid phone")
	}

	return nil
}
