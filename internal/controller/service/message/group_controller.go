package message

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	group "chickChirick/internal/model/message"
	"encoding/json"
	"errors"
	"net/http"
)

type GroupController struct {
	Controller abstraction.Controller
}

func (gc *GroupController) HandleRequest() {
	gc.Controller.ServeMux.HandleFunc("GET /groups", func(w http.ResponseWriter, r *http.Request) {
		gc.GetGroups(w)
	})

	gc.Controller.ServeMux.HandleFunc("POST /group", func(w http.ResponseWriter, r *http.Request) {
		gc.CreateGroup(w, c_http.NewRequest(r))
	})

	gc.Controller.ServeMux.HandleFunc("GET /group/{id}", func(w http.ResponseWriter, r *http.Request) {
		gc.GetGroup(w, c_http.NewRequest(r))
	})

	gc.Controller.ServeMux.HandleFunc("PUT /group/{id}", func(w http.ResponseWriter, r *http.Request) {
		gc.UpdateGroup(w, c_http.NewRequest(r))
	})

	gc.Controller.ServeMux.HandleFunc("DELETE /group/{id}", func(w http.ResponseWriter, r *http.Request) {
		gc.DeleteGroup(w, c_http.NewRequest(r))
	})
}

func (gc *GroupController) GetGroups(w http.ResponseWriter) {
	groups, err := group.GetGroups(gc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, groups, http.StatusOK)
}

func (gc *GroupController) GetGroup(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	g, err := group.GetGroupById(gc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Group not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, g, http.StatusOK)
}

func (gc *GroupController) CreateGroup(w http.ResponseWriter, r *c_http.Request) {
	var g group.Group
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := group.CreateGroup(gc.Controller.Dependencies.DBDecorator.GDB(), &g); err != nil {
		c_http.NewResponse().SendError(w, "Failed to create group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, g, http.StatusCreated)
}

func (gc *GroupController) UpdateGroup(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var g group.Group
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = group.UpdateGroupById(gc.Controller.Dependencies.DBDecorator.GDB(), &g, id)
	if err != nil && errors.Is(err, group.GroupNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, g, http.StatusOK)
}

func (gc *GroupController) DeleteGroup(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = group.DeleteGroupById(gc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil && errors.Is(err, group.GroupNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
