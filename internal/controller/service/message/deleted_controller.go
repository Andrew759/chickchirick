package message

import (
	"chickChirick/internal/controller/c_controller"
	"chickChirick/internal/controller/c_http"
	deleted "chickChirick/internal/model/message"
	"encoding/json"
	"errors"
	"net/http"
)

type DeletedController struct {
	Controller c_controller.Controller
}

func (dc *DeletedController) HandleRequest() {
	dc.Controller.ServeMux.HandleFunc("GET /deleted-list", func(w http.ResponseWriter, r *http.Request) {
		dc.GetDeletedList(w)
	})

	dc.Controller.ServeMux.HandleFunc("POST /deleted", func(w http.ResponseWriter, r *http.Request) {
		dc.CreateDeleted(w, c_http.NewRequest(r))
	})

	dc.Controller.ServeMux.HandleFunc("GET /deleted/{id}", func(w http.ResponseWriter, r *http.Request) {
		dc.GetDeleted(w, c_http.NewRequest(r))
	})

	dc.Controller.ServeMux.HandleFunc("PUT /deleted/{id}", func(w http.ResponseWriter, r *http.Request) {
		dc.UpdateDeleted(w, c_http.NewRequest(r))
	})

	dc.Controller.ServeMux.HandleFunc("DELETE /deleted/{id}", func(w http.ResponseWriter, r *http.Request) {
		dc.DeleteDeleted(w, c_http.NewRequest(r))
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
	id, err := r.HTTPId()
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
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var d deleted.Deleted
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = deleted.UpdateDeletedById(dc.Controller.Dependencies.DBDecorator.GDB(), &d, id)
	if err != nil && errors.Is(err, deleted.DeletedNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update deleted: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, d, http.StatusOK)
}

func (dc *DeletedController) DeleteDeleted(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = deleted.DeleteDeletedById(dc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil && errors.Is(err, deleted.DeletedNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete deleted: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
