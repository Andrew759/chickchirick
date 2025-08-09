package message

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	status "chickChirick/internal/model/message"
	"encoding/json"
	"net/http"
)

const StatusResource = "/message/status/"

type StatusController struct {
	Controller abstraction.Controller
}

func (sc *StatusController) HandleRequest() {
	sc.Controller.ServeMux.HandleFunc("/message/statuses", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			sc.GetStatuses(w)
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	sc.Controller.ServeMux.HandleFunc(StatusResource, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			sc.GetStatus(w, r)
		case http.MethodPost:
			sc.CreateStatus(w, r)
		case http.MethodPut:
			sc.UpdateStatus(w, r)
		case http.MethodDelete:
			sc.DeleteStatus(w, r)
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
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

func (sc *StatusController) GetStatus(w http.ResponseWriter, r *http.Request) {
	id := sc.Controller.HttpId(w, r, StatusResource)
	s, err := status.GetStatusById(sc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Status not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, s, http.StatusOK)
}

func (sc *StatusController) CreateStatus(w http.ResponseWriter, r *http.Request) {
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

func (sc *StatusController) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	var s status.Status
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := status.UpdateStatus(sc.Controller.Dependencies.DBDecorator.GDB(), &s)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, s, http.StatusOK)
}

func (sc *StatusController) DeleteStatus(w http.ResponseWriter, r *http.Request) {
	id := sc.Controller.HttpId(w, r, StatusResource)
	err := status.DeleteStatusById(sc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
