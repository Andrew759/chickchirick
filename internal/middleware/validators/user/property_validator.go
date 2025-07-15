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

type PropertyValidator struct{}

func (pc PropertyValidator) Validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p user.Property

		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := validateProperty(p); err != nil {
			http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), config.UserPropertyKey, &p)
		next(w, r.WithContext(ctx))
	}
}

func validateProperty(p user.Property) error {
	if strings.TrimSpace(p.Email) != "" && !service.IsEmail(p.Email) {
		return errors.New("invalid email")
	}
	if p.Password != nil || !service.IsHasCorrectLength(*p.Password, 1024) {
		return errors.New("invalid password")
	}
	//TODO: валидация таймзон

	return nil
}
