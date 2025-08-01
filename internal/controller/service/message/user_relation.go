package message

import (
	"chickChirick/internal/controller/abstraction"
	userRelation "chickChirick/internal/model/message"
	"encoding/json"
	"net/http"
)

type UserRelationController struct {
	Controller abstraction.Controller
}

func (urc *UserRelationController) HandleRequest(mux *http.ServeMux) {
	mux.HandleFunc("/message/user-relations", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			urc.GetUserRelations(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/message/user-relation", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			urc.GetUserRelation(w, r)
		case http.MethodPost:
			urc.CreateUserRelation(w, r)
		case http.MethodPut:
			urc.UpdateUserRelation(w, r)
		case http.MethodDelete:
			urc.DeleteUserRelation(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (urc *UserRelationController) GetUserRelations(w http.ResponseWriter) {
	userRelations, err := userRelation.GetUserRelation(urc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userRelations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (urc *UserRelationController) GetUserRelation(w http.ResponseWriter, r *http.Request) {
	id := urc.Controller.GETId(w, r)
	s, err := userRelation.GetUserRelationById(urc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "User relation not found: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (urc *UserRelationController) CreateUserRelation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var ur userRelation.UserRelation
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := userRelation.CreateUserRelation(urc.Controller.Dependencies.DBDecorator.GDB(), &ur); err != nil {
		http.Error(w, "Failed to create user relation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(ur); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (urc *UserRelationController) UpdateUserRelation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var ur userRelation.UserRelation
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := userRelation.UpdateUserRelation(urc.Controller.Dependencies.DBDecorator.GDB(), &ur)
	if err != nil {
		http.Error(w, "Failed to update user relation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(ur); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (urc *UserRelationController) DeleteUserRelation(w http.ResponseWriter, r *http.Request) {
	id := urc.Controller.GETId(w, r)
	err := userRelation.DeleteUserRelationById(urc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Failed to delete user relation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
