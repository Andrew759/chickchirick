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

	"gorm.io/gorm"
)

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

	err := cup.Controller.Dependencies.DBDecorator.GDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := user.CreateUser(ctx, tx, u); err != nil {
			return err
		}

		if _, err := cup.createMeta(ctx, tx, u); err != nil {
			return err
		}

		if err := cup.createProperty(ctx, tx, p); err != nil {
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

func (cup *CreateUserProxy) createProperty(ctx context.Context, tx *gorm.DB, p *user.Property) error {
	if p == nil {
		return nil
	}
	return user.CreateProperty(ctx, tx, p)
}
