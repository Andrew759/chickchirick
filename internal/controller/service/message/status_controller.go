package message

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	status "chickChirick/internal/model/message"
	"encoding/json"
	"errors"
	"net/http"
)

type StatusController struct {
	Controller abstraction.Controller
}

func (sc *StatusController) HandleRequest() {
	sc.Controller.ServeMux.HandleFunc("GET /statuses", func(w http.ResponseWriter, r *http.Request) {
		sc.GetStatuses(w)
	})

	sc.Controller.ServeMux.HandleFunc("POST /status", func(w http.ResponseWriter, r *http.Request) {
		sc.CreateStatus(w, c_http.NewRequest(r))
	})

	sc.Controller.ServeMux.HandleFunc("GET /status/{id}", func(w http.ResponseWriter, r *http.Request) {
		sc.GetStatus(w, c_http.NewRequest(r))
	})

	sc.Controller.ServeMux.HandleFunc("PUT /status/{id}", func(w http.ResponseWriter, r *http.Request) {
		sc.UpdateStatus(w, c_http.NewRequest(r))
	})

	sc.Controller.ServeMux.HandleFunc("DELETE /status/{id}", func(w http.ResponseWriter, r *http.Request) {
		sc.DeleteStatus(w, c_http.NewRequest(r))
	})
}

func (sc *StatusController) GetStatuses(w http.ResponseWriter) {
	statuses, err := status.GetStatus(sc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, statuses, http.StatusOK)
}

func (sc *StatusController) GetStatus(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	s, err := status.GetStatusById(sc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Status not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, s, http.StatusOK)
}

func (sc *StatusController) CreateStatus(w http.ResponseWriter, r *c_http.Request) {
	var s status.Status
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := status.CreateStatus(sc.Controller.Dependencies.DBDecorator.GDB(), &s); err != nil {
		c_http.NewResponse().SendError(w, "Failed to create status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, s, http.StatusCreated)
}

func (sc *StatusController) UpdateStatus(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var s status.Status
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = status.UpdateStatusById(sc.Controller.Dependencies.DBDecorator.GDB(), &s, id)
	if err != nil && errors.Is(err, status.StatusNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, s, http.StatusOK)
}

func (sc *StatusController) DeleteStatus(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = status.DeleteStatusById(sc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil && errors.Is(err, status.StatusNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
