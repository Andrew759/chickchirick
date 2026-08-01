package user

import (
	"chickChirick/internal/controller/c_controller"
	"chickChirick/internal/controller/c_http"
	"chickChirick/internal/middleware"
	"chickChirick/internal/middleware/config"
	"chickChirick/internal/model/user"
	"context"
	"errors"
	"net/http"
)

type CreateUserProxy struct {
	Controller c_controller.Controller
	middleware.Validator
}

func (cup *CreateUserProxy) HandleRequest() {
	cup.Controller.ServeMux.HandleFunc("POST /frontend/user", func(w http.ResponseWriter, r *http.Request) {
		cup.Validate(func(w http.ResponseWriter, r *http.Request) {
			cup.CreateUser(w, c_http.NewRequest(r))
		})
	})
}

func (cup *CreateUserProxy) CreateUser(w http.ResponseWriter, r *c_http.Request) {
	ctx := r.Context()
	u := ctx.Value(config.UserProxyKey).(*user.User)

	err := user.CreateUser(ctx, cup.Controller.Dependencies.DBDecorator.GDB(), u)

	var userAlreadyExistError *user.UserAlreadyExistErr
	if err != nil && errors.As(err, &userAlreadyExistError) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusConflict)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to create user. "+err.Error(), http.StatusInternalServerError)
		return
	}

	cup.createMeta(ctx, w, u)

}

func (cup *CreateUserProxy) createMeta(ctx context.Context, w http.ResponseWriter, u *user.User) user.Meta {
	m := user.Meta{
		UserId: u.Id,
	}

	if err := user.CreateMeta(ctx, cup.Controller.Dependencies.DBDecorator.GDB(), &m); err != nil {
		c_http.NewResponse().SendError(w, "Failed to create meta: "+err.Error(), http.StatusInternalServerError)
		return m
	}

	return m
}

func (cup *CreateUserProxy) createProperty(ctx context.Context, w http.ResponseWriter, u *user.User) {

}
