package user

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/middleware/config"
	userMiddleware "chickChirick/internal/middleware/validators/user"
	"chickChirick/internal/model/user"
	"encoding/json"
	"net/http"
)

type UserController struct {
	Controller abstraction.Controller
	userMiddleware.UserValidator
}

func (uc *UserController) HandleRequest() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			uc.GetUsers(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (uc *UserController) GetUsers(w http.ResponseWriter) {
	users, err := user.GetAllUsers(uc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	id := uc.Controller.GETId(w, r)
	u, err := user.GetUserById(uc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := r.Context().Value(config.UserUserKey).(*user.User)

	if err := user.CreateUser(uc.Controller.Dependencies.DBDecorator.GDB(), u); err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := r.Context().Value(config.UserUserKey).(*user.User)

	err := user.UpdateUser(uc.Controller.Dependencies.DBDecorator.GDB(), u)
	if err != nil {
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := uc.Controller.GETId(w, r)
	err := user.DeleteUserById(uc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
