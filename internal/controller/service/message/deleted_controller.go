package message

import (
	"chickChirick/internal/controller/abstraction"
	deleted "chickChirick/internal/model/message"
	"encoding/json"
	"net/http"
)

const DeletedResource = "/message/deleted/"

type DeletedController struct {
	Controller abstraction.Controller
}

func (dc *DeletedController) HandleRequest() {
	dc.Controller.ServeMux.HandleFunc("/message/deleted-list", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			dc.GetDeletedList(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	dc.Controller.ServeMux.HandleFunc(DeletedResource, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			dc.GetDeleted(w, r)
		case http.MethodPost:
			dc.CreateDeleted(w, r)
		case http.MethodDelete:
			dc.DeleteDeleted(w, r)
		case http.MethodPut:
			dc.UpdateDeleted(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (dc *DeletedController) GetDeletedList(w http.ResponseWriter) {
	deletedList, err := deleted.GetDeleted(dc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(deletedList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (dc *DeletedController) GetDeleted(w http.ResponseWriter, r *http.Request) {
	id := dc.Controller.HttpId(w, r, DeletedResource)
	d, err := deleted.GetDeletedById(dc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Deleted not found: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (dc *DeletedController) CreateDeleted(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var d deleted.Deleted
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := deleted.CreateDeleted(dc.Controller.Dependencies.DBDecorator.GDB(), &d); err != nil {
		http.Error(w, "Failed to create deleted: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (dc *DeletedController) UpdateDeleted(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var d deleted.Deleted
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := deleted.UpdateDeleted(dc.Controller.Dependencies.DBDecorator.GDB(), &d)
	if err != nil {
		http.Error(w, "Failed to update deleted: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (dc *DeletedController) DeleteDeleted(w http.ResponseWriter, r *http.Request) {
	id := dc.Controller.HttpId(w, r, DeletedResource)
	err := deleted.DeleteDeletedById(dc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Failed to delete deleted: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
