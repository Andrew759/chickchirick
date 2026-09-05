package user

import (
	"bytes"
	"chickChirick/internal/controller/c_controller"
	"chickChirick/internal/controller/c_http"
	"chickChirick/internal/middleware"
	"chickChirick/internal/middleware/config"
	"chickChirick/internal/model/user"
	"chickChirick/pkg/chirik_config"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// TODO: отрефакторить, поправить баг с введением емейл
type CreateAuthUserRequest struct {
	Password string    `json:"password"`
	UserUuid uuid.UUID `json:"user_uuid"`
}

type CreateMessageUserRequest struct {
	UserUuid uuid.UUID `json:"user_uuid"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CreateUserProxy struct {
	Controller c_controller.Controller
	UPV        middleware.Validator
	PV         middleware.Validator
}

func (cup *CreateUserProxy) HandleRequest() {
	cup.Controller.ServeMux.HandleFunc("POST /frontend/user", func(w http.ResponseWriter, r *http.Request) {
		cup.UPV.Validate(func(w http.ResponseWriter, r *http.Request) {
			cup.PV.Validate(func(w http.ResponseWriter, r *http.Request) {
				cup.CreateUser(w, c_http.NewRequest(r))
			})(w, r)
		})(w, r)
	})

	cup.Controller.ServeMux.HandleFunc("GET /frontend/user/me", func(w http.ResponseWriter, r *http.Request) {
		cup.GetMe(w, c_http.NewRequest(r))
	})
}

func (cup *CreateUserProxy) CreateUser(w http.ResponseWriter, r *c_http.Request) {
	ctx := r.Context()

	u, ok := ctx.Value(config.UserUserKey).(*user.User)
	if !ok {
		c_http.NewResponse().SendError(w, "User identity missing in context", http.StatusInternalServerError)
		return
	}

	p, ok := ctx.Value(config.UserPropertyKey).(*user.Property)
	if !ok {
		c_http.NewResponse().SendError(w, "Property identity missing in context", http.StatusInternalServerError)
		return
	}

	var meta user.Meta

	err := cup.Controller.Dependencies.DBDecorator.GDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := user.CreateUser(ctx, tx, u); err != nil {
			return err
		}

		var err error
		meta, err = cup.createMeta(ctx, tx, u)
		if err != nil {
			return err
		}

		if err := cup.createProperty(ctx, tx, p, u); err != nil {
			return err
		}

		return nil
	})

	var userAlreadyExistError *user.UserAlreadyExistErr
	if err != nil && errors.As(err, &userAlreadyExistError) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusConflict)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to create user sequence. "+err.Error(), http.StatusInternalServerError)
		return
	}

	//TODO: оба запроса во внешние сервисы должны выполняться безусловно, а транзакция выше должна откатываться
	// тут нужно применить реббит или кафку
	tokens, err := cup.createAuthUser(ctx, p, meta)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to auth. "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = cup.createMessageUser(ctx, meta)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to create message user. "+err.Error(), http.StatusInternalServerError)
		return
	}

	cookieAccess := &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, //TODO: для https: true
		SameSite: http.SameSiteLaxMode,
	}

	cookieRefresh := &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	w.Header().Add("Set-Cookie", cookieAccess.String())
	w.Header().Add("Set-Cookie", cookieRefresh.String())

	c_http.NewResponse().SendSuccess(w, u, http.StatusCreated)
}

func (cup *CreateUserProxy) createMeta(ctx context.Context, tx *gorm.DB, u *user.User) (user.Meta, error) {
	m := user.Meta{
		UserId: u.Id,
	}

	if err := user.CreateMeta(ctx, tx, &m); err != nil {
		return m, err
	}

	return m, nil
}

func (cup *CreateUserProxy) createProperty(ctx context.Context, tx *gorm.DB, p *user.Property, u *user.User) error {
	if p == nil {
		return nil
	}

	if p.UserId == 0 {
		p.UserId = u.Id
	}

	return user.CreateProperty(ctx, tx, p)
}

func (cup *CreateUserProxy) createAuthUser(
	ctx context.Context,
	p *user.Property,
	m user.Meta,
) (Tokens, error) {
	var tokens Tokens

	reqData := CreateAuthUserRequest{
		UserUuid: m.UserUuid,
	}
	if p != nil && p.Password != nil {
		reqData.Password = *p.Password
	}

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		slog.Error("error marshaling auth request data: ", err.Error())
		return tokens, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", viper.GetString(chirik_config.AuthAppUrl)+"/user/create", bytes.NewBuffer(reqBody))
	if err != nil {
		return tokens, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := cup.Controller.Dependencies.Client.Do(req)
	if err != nil {
		slog.Error("error sending request to auth service: ", err.Error())
		return tokens, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			slog.Error("error closing auth response body: ", err.Error())
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return tokens, errors.New("auth service returned non-ok status")
	}

	for _, cookie := range resp.Cookies() {
		switch cookie.Name {
		case "access_token":
			tokens.AccessToken = cookie.Value
		case "refresh_token":
			tokens.RefreshToken = cookie.Value
		}
	}

	return tokens, nil
}

func (cup *CreateUserProxy) createMessageUser(ctx context.Context, m user.Meta) error {
	reqData := CreateMessageUserRequest{
		UserUuid: m.UserUuid,
	}

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		slog.Error("error marshaling message request data: ", err.Error())
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", viper.GetString(chirik_config.MessageAppUrl)+"/user-relation", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := cup.Controller.Dependencies.Client.Do(req)
	if err != nil {
		slog.Error("error sending request to message service: ", err.Error())
		return err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			slog.Error("error closing message response body: ", err.Error())
		}
	}(resp.Body)
	return nil
}

func (cup *CreateUserProxy) GetMe(w http.ResponseWriter, r *c_http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie("access_token")
	if err != nil {
		c_http.NewResponse().SendError(w, "Unauthorized: missing token", http.StatusUnauthorized)
		return
	}

	req, err := http.NewRequestWithContext(ctx, "GET", viper.GetString(chirik_config.AuthAppUrl)+"/auth/validate", nil)
	if err != nil {
		c_http.NewResponse().SendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	req.AddCookie(cookie)

	resp, err := cup.Controller.Dependencies.Client.Do(req)
	if err != nil {
		slog.Error("error sending request to auth service: ", err.Error())
		c_http.NewResponse().SendError(w, "Auth service unavailable", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c_http.NewResponse().SendError(w, "Unauthorized: invalid session", http.StatusUnauthorized)
		return
	}

	var authResult struct {
		Payload struct {
			Valid    bool   `json:"valid"`
			UserUuid string `json:"user_uuid"`
		} `json:"payload"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&authResult); err != nil {
		slog.Error("error decoding token validation response: ", err.Error())
		c_http.NewResponse().SendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, authResult, http.StatusOK)
}
