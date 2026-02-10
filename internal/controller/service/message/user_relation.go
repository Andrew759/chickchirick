package message

import (
	"chickChirick/internal/controller/c_controller"
	"chickChirick/internal/controller/c_http"
	userRelation "chickChirick/internal/model/message"
	"encoding/json"
	"errors"
	"net/http"
)

type UserRelationController struct {
	Controller c_controller.Controller
}

// TODO: тут скорее всего баг
func (urc *UserRelationController) HandleRequest() {
	urc.Controller.ServeMux.HandleFunc("GET /user-relations", func(w http.ResponseWriter, r *http.Request) {
		urc.GetUserRelations(w)
	})

	urc.Controller.ServeMux.HandleFunc("POST /user-relation", func(w http.ResponseWriter, r *http.Request) {
		urc.CreateUserRelation(w, c_http.NewRequest(r))
	})

	urc.Controller.ServeMux.HandleFunc("GET /user-relation/{id}", func(w http.ResponseWriter, r *http.Request) {
		urc.GetUserRelation(w, c_http.NewRequest(r))
	})

	urc.Controller.ServeMux.HandleFunc("PUT /user-relation/{id}", func(w http.ResponseWriter, r *http.Request) {
		urc.UpdateUserRelation(w, c_http.NewRequest(r))
	})

	urc.Controller.ServeMux.HandleFunc("DELETE /user-relation/{id}", func(w http.ResponseWriter, r *http.Request) {
		urc.DeleteUserRelation(w, c_http.NewRequest(r))
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

func (urc *UserRelationController) GetUserRelation(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ur, err := userRelation.GetUserRelationById(urc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "User relation not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, ur, http.StatusOK)
}

func (urc *UserRelationController) CreateUserRelation(w http.ResponseWriter, r *c_http.Request) {
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

func (urc *UserRelationController) UpdateUserRelation(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ur userRelation.UserRelation
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = userRelation.UpdateUserRelationById(urc.Controller.Dependencies.DBDecorator.GDB(), &ur, id)
	if err != nil && errors.Is(err, userRelation.UserRelationNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update user relation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, ur, http.StatusOK)
}

func (urc *UserRelationController) DeleteUserRelation(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = userRelation.DeleteUserRelationById(urc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil && errors.Is(err, userRelation.UserRelationNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete user relation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
