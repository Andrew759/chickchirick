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

func (controller *UserController) HandleRequest() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetUsers(w, r)
		case http.MethodPost:
			CreateUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			GetUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

var users = []user.User{
	user.User{
		Id:       1,
		Phone:    79634823344,
		Name:     "Andrey",
		Surname:  "Velkov",
		Password: nil,
	},
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Get user by ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	//implement this
}

// Create user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	//implement this
}
