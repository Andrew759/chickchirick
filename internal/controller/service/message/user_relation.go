package message

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	userRelation "chickChirick/internal/model/message"
	"encoding/json"
	"net/http"
)

const UserRelationResource = "/message/user-relation/"

type UserRelationController struct {
	Controller abstraction.Controller
}

func (urc *UserRelationController) HandleRequest() {
	urc.Controller.ServeMux.HandleFunc("/message/user-relations", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			urc.GetUserRelations(w)
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	urc.Controller.ServeMux.HandleFunc(UserRelationResource, func(w http.ResponseWriter, r *http.Request) {
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
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (urc *UserRelationController) GetUserRelations(w http.ResponseWriter) {
	userRelations, err := userRelation.GetUserRelation(urc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, userRelations, http.StatusOK)
}

func (urc *UserRelationController) GetUserRelation(w http.ResponseWriter, r *http.Request) {
	id := urc.Controller.HttpId(w, r, UserRelationResource)
	ur, err := userRelation.GetUserRelationById(urc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "User relation not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, ur, http.StatusOK)
}

func (urc *UserRelationController) CreateUserRelation(w http.ResponseWriter, r *http.Request) {
	var ur userRelation.UserRelation
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := userRelation.CreateUserRelation(urc.Controller.Dependencies.DBDecorator.GDB(), &ur); err != nil {
		c_http.NewResponse().SendError(w, "Failed to create user relation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, ur, http.StatusCreated)
}

func (urc *UserRelationController) UpdateUserRelation(w http.ResponseWriter, r *http.Request) {
	var ur userRelation.UserRelation
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := userRelation.UpdateUserRelation(urc.Controller.Dependencies.DBDecorator.GDB(), &ur)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update user relation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, ur, http.StatusOK)
}

func (urc *UserRelationController) DeleteUserRelation(w http.ResponseWriter, r *http.Request) {
	id := urc.Controller.HttpId(w, r, UserRelationResource)
	err := userRelation.DeleteUserRelationById(urc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete user relation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
