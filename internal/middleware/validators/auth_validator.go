package validators

import (
	"chickChirick/internal/controller/c_http"
	"chickChirick/internal/middleware"
	"chickChirick/internal/middleware/config"
	"context"
	"encoding/json"
	"net/http"
)

type AuthValidatorFactory struct{}

func (avf AuthValidatorFactory) NewValidator(authURL string, client *http.Client, opts ...middleware.ValidatorOption) AuthValidator {
	var vOptions middleware.ValidatorOptions
	for _, opt := range opts {
		opt(&vOptions)
	}

	return AuthValidator{
		AuthServiceURL:   authURL,
		Client:           client,
		ValidatorOptions: vOptions,
	}
}

type ValidateResponsePayload struct {
	Valid    bool   `json:"valid"`
	UserUuid string `json:"user_uuid"`
}

type ValidateResponse struct {
	ValidateResponsePayload `json:"payload"`
	Error                   string `json:"error"`
}

type AuthValidator struct {
	AuthServiceURL string
	Client         *http.Client
	middleware.ValidatorOptions
}

func (av AuthValidator) ValidateAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil || cookie.Value == "" {
			c_http.NewResponse().SendError(w, "invalid access token", http.StatusUnauthorized)
			return
		}

		req, err := http.NewRequestWithContext(r.Context(), "GET", av.AuthServiceURL+"/auth/validate", nil)
		if err != nil {
			c_http.NewResponse().SendError(w, ""+
				"internal error at constructing request to the authorization service",
				http.StatusInternalServerError,
			)
			return
		}
		req.AddCookie(cookie)

		resp, err := av.Client.Do(req)
		if err != nil {
			c_http.NewResponse().SendError(w, "auth service unavailable", http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		var result ValidateResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			c_http.NewResponse().SendError(w, "failed to parse auth response", http.StatusInternalServerError)
			return
		}

		if !result.Valid || resp.StatusCode != http.StatusOK {
			c_http.NewResponse().SendError(w, "forbidden: "+result.Error, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), config.AuthKey, result.UserUuid)

		next(w, r.WithContext(ctx))
	}
}
