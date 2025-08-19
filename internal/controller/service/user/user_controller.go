package user

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	"chickChirick/internal/middleware/config"
	userMiddleware "chickChirick/internal/middleware/validators/user"
	"chickChirick/internal/model/user"
	"errors"
	"net/http"
)

type UserController struct {
	Controller abstraction.Controller
	userMiddleware.UserValidator
}

func (uc *UserController) HandleRequest() {
	uc.Controller.ServeMux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			uc.GetUsers(w)
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	uc.Controller.ServeMux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			uc.Validate(func(w http.ResponseWriter, r *http.Request) {
				uc.CreateUser(w, c_http.NewRequest(r))
			})(w, r)
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	uc.Controller.ServeMux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			uc.GetUser(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/user/")))
		case http.MethodPut:
			uc.Validate(func(w http.ResponseWriter, r *http.Request) {
				uc.UpdateUser(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/user/")))
			})(w, r)
		case http.MethodDelete:
			uc.DeleteUser(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/user/")))
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (uc *UserController) GetUsers(w http.ResponseWriter) {
	users, err := user.GetAllUsers(uc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, users, http.StatusOK)
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
	}

	u, err := user.GetUserById(uc.Controller.Dependencies.DBDecorator.GDB(), id)
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
	u := r.Context().Value(config.UserUserKey).(*user.User)

	err := user.CreateUser(uc.Controller.Dependencies.DBDecorator.GDB(), u)

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
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	u := r.Context().Value(config.UserUserKey).(*user.User)

	err = user.UpdateUserById(uc.Controller.Dependencies.DBDecorator.GDB(), u, id)
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
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	//TODO: добавить ошибку, что пользователь не найден
	err = user.DeleteUserById(uc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete user. "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, nil, http.StatusNoContent)
}
