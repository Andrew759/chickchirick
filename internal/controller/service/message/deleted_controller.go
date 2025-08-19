package message

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	deleted "chickChirick/internal/model/message"
	"encoding/json"
	"net/http"
)

type DeletedController struct {
	Controller abstraction.Controller
}

func (dc *DeletedController) HandleRequest() {
	dc.Controller.ServeMux.HandleFunc("/message/deleted-list", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			dc.GetDeletedList(w)
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	dc.Controller.ServeMux.HandleFunc("/message/deleted", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			dc.CreateDeleted(w, c_http.NewRequest(r))
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	dc.Controller.ServeMux.HandleFunc("/message/deleted/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			dc.GetDeleted(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/message/deleted/")))
		case http.MethodDelete:
			dc.DeleteDeleted(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/message/deleted/")))
		case http.MethodPut:
			dc.UpdateDeleted(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/message/deleted/")))
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (dc *DeletedController) GetDeletedList(w http.ResponseWriter) {
	deletedList, err := deleted.GetDeleted(dc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, deletedList, http.StatusOK)
}

func (dc *DeletedController) GetDeleted(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	d, err := deleted.GetDeletedById(dc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Deleted not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, d, http.StatusOK)
}

func (dc *DeletedController) CreateDeleted(w http.ResponseWriter, r *c_http.Request) {
	var d deleted.Deleted
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := deleted.CreateDeleted(dc.Controller.Dependencies.DBDecorator.GDB(), &d); err != nil {
		c_http.NewResponse().SendError(w, "Failed to create deleted: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, d, http.StatusCreated)
}

func (dc *DeletedController) UpdateDeleted(w http.ResponseWriter, r *c_http.Request) {
	var d deleted.Deleted
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := deleted.UpdateDeleted(dc.Controller.Dependencies.DBDecorator.GDB(), &d)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update deleted: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, d, http.StatusOK)
}

func (dc *DeletedController) DeleteDeleted(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = deleted.DeleteDeletedById(dc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete deleted: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
