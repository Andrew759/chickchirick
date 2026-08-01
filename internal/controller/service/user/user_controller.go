package user

import (
	"chickChirick/internal/controller/c_controller"
	"chickChirick/internal/controller/c_http"
	"chickChirick/internal/middleware"
	"chickChirick/internal/middleware/config"
	authValidator "chickChirick/internal/middleware/validator"
	"chickChirick/internal/model/user"
	"errors"
	"net/http"
)

type UserController struct {
	Controller c_controller.Controller
	middleware.Validator
	authValidator.AuthValidator
}

func (uc *UserController) HandleRequest() {
	uc.Controller.ServeMux.HandleFunc("GET /users",
		uc.ValidateAuth(func(w http.ResponseWriter, r *http.Request) {
			uc.GetUsers(w, c_http.NewRequest(r))
		}))

	uc.Controller.ServeMux.HandleFunc("POST /user",
		uc.Validate(func(w http.ResponseWriter, r *http.Request) {
			uc.CreateUser(w, c_http.NewRequest(r))
		}))

	uc.Controller.ServeMux.HandleFunc("GET /user/{id}",
		uc.ValidateAuth(func(w http.ResponseWriter, r *http.Request) {
			uc.GetUser(w, c_http.NewRequest(r))
		}))

	uc.Controller.ServeMux.HandleFunc("PUT /user/{id}",
		uc.ValidateAuth(uc.Validate(func(w http.ResponseWriter, r *http.Request) {
			uc.UpdateUser(w, c_http.NewRequest(r))
		})))

	uc.Controller.ServeMux.HandleFunc("DELETE /user/{id}",
		uc.ValidateAuth(func(w http.ResponseWriter, r *http.Request) {
			uc.DeleteUser(w, c_http.NewRequest(r))
		}))
}

func (uc *UserController) GetUsers(w http.ResponseWriter, r *c_http.Request) {
	ctx := r.Context()
	users, err := user.GetAllUsers(ctx, uc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, users, http.StatusOK)
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	u, err := user.GetUserById(ctx, uc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil && errors.Is(err, user.UserNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, u, http.StatusOK)
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *c_http.Request) {
	ctx := r.Context()
	u := ctx.Value(config.UserUserKey).(*user.User)

	err := user.CreateUser(ctx, uc.Controller.Dependencies.DBDecorator.GDB(), u)

	var userAlreadyExistError *user.UserAlreadyExistErr
	if err != nil && errors.As(err, &userAlreadyExistError) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusConflict)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to create user. "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, u, http.StatusCreated)
}

func (uc *UserController) UpdateUser(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	u := ctx.Value(config.UserUserKey).(*user.User)

	err = user.UpdateUserById(ctx, uc.Controller.Dependencies.DBDecorator.GDB(), u, id)
	if err != nil && errors.Is(err, user.UserNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update user. "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, u, http.StatusOK)
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = user.DeleteUserById(ctx, uc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil && errors.Is(err, user.UserNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete user. "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, nil, http.StatusNoContent)
}
