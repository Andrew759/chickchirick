package service

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/model/user"
	"encoding/json"
	"net/http"
	"strconv"
)

type UserController struct {
	AbstractController abstraction.Controller
}

func (uc *UserController) HandleRequest() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			uc.GetUsers(w)
		case http.MethodPost:
			uc.CreateUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			uc.GetUser(w, r)
		case http.MethodPost:
			uc.CreateUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (uc *UserController) GetUsers(w http.ResponseWriter) {
	users, err := user.GetAllUsers(uc.AbstractController.Dependencies.DBDecorator.GDB())
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
	w.Header().Set("Content-Type", "application/json")
	idVal := r.URL.Query().Get("id")
	if idVal == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(idVal)
	u, err := user.GetUserById(uc.AbstractController.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
		return
	}

	// Возвращаем пользователя
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var u user.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := user.CreateUser(uc.AbstractController.Dependencies.DBDecorator.GDB(), &u); err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
