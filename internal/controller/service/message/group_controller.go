package message

import (
	"chickChirick/internal/controller/abstraction"
	group "chickChirick/internal/model/message"
	"encoding/json"
	"net/http"
)

type GroupController struct {
	Controller abstraction.Controller
}

func (gc *GroupController) HandleRequest(mux *http.ServeMux) {
	mux.HandleFunc("/message/groups", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			gc.GetGroups(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/message/group", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			gc.GetGroup(w, r)
		case http.MethodPost:
			gc.CreateGroup(w, r)
		case http.MethodPut:
			gc.UpdateGroup(w, r)
		case http.MethodDelete:
			gc.DeleteGroup(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (gc *GroupController) GetGroups(w http.ResponseWriter) {
	groups, err := group.GetGroups(gc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(groups)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (gc *GroupController) GetGroup(w http.ResponseWriter, r *http.Request) {
	id := gc.Controller.GETId(w, r)
	g, err := group.GetGroupById(gc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Group not found: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (gc *GroupController) CreateGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var g group.Group
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := group.CreateGroup(gc.Controller.Dependencies.DBDecorator.GDB(), &g); err != nil {
		http.Error(w, "Failed to create group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (gc *GroupController) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var g group.Group
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := group.UpdateGroup(gc.Controller.Dependencies.DBDecorator.GDB(), &g)
	if err != nil {
		http.Error(w, "Failed to update group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (gc *GroupController) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	id := gc.Controller.GETId(w, r)
	err := group.DeleteGroupById(gc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Failed to delete group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
