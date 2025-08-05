package user

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/http_transaction"
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
			http_transaction.NewResponse().Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	uc.Controller.ServeMux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			uc.GetUser(w, r)
		case http.MethodPost:
			uc.Validate(uc.CreateUser)(w, r)
		case http.MethodPut:
			uc.Validate(uc.UpdateUser)(w, r)
		case http.MethodDelete:
			uc.DeleteUser(w, r)
		default:
			http_transaction.NewResponse().Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})
}

func (uc *UserController) GetUsers(w http.ResponseWriter) {
	users, err := user.GetAllUsers(uc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		http_transaction.NewResponse().Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	http_transaction.NewResponse().Success(w, http.StatusOK, users)
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	id := uc.Controller.GETId(w, r)
	u, err := user.GetUserById(uc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil && errors.Is(err, user.UserNotFoundErr) {
		http_transaction.NewResponse().Error(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		http_transaction.NewResponse().Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	http_transaction.NewResponse().Success(w, http.StatusOK, u)
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := r.Context().Value(config.UserUserKey).(*user.User)

	err := user.CreateUser(uc.Controller.Dependencies.DBDecorator.GDB(), u)

	var userAlreadyExistError *user.UserAlreadyExistErr
	if err != nil && errors.As(err, &userAlreadyExistError) {
		http_transaction.NewResponse().Error(w, http.StatusConflict, err.Error())
		return
	} else if err != nil {
		http_transaction.NewResponse().Error(w, http.StatusInternalServerError, "Failed to create user. "+err.Error())
		return
	}

	http_transaction.NewResponse().Success(w, http.StatusCreated, u)
}

func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := r.Context().Value(config.UserUserKey).(*user.User)

	err := user.UpdateUser(uc.Controller.Dependencies.DBDecorator.GDB(), u)
	if err != nil {
		http_transaction.NewResponse().Error(w, http.StatusInternalServerError, "Failed to update user. "+err.Error())
		return
	}

	http_transaction.NewResponse().Success(w, http.StatusOK, u)
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := uc.Controller.GETId(w, r)
	err := user.DeleteUserById(uc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http_transaction.NewResponse().Error(w, http.StatusInternalServerError, "Failed to delete user. "+err.Error())
		return
	}

	http_transaction.NewResponse().Success(w, http.StatusNoContent, nil)
}
