package service

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/model/user"
	"encoding/json"
	"net/http"
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
		//TODO: implement this
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		//TODO: implement this
	}
}

// GetUser Get user by ID
func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	//implement this
}

// Create user
func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	//implement this
}
