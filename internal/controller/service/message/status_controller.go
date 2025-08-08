package message

import (
	"chickChirick/internal/controller/abstraction"
	status "chickChirick/internal/model/message"
	"encoding/json"
	"net/http"
)

type StatusController struct {
	Controller abstraction.Controller
}

func (sc *StatusController) HandleRequest() {
	sc.Controller.ServeMux.HandleFunc("/message/statuses", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			sc.GetStatuses(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	sc.Controller.ServeMux.HandleFunc("/message/status", func(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (sc *StatusController) GetStatuses(w http.ResponseWriter) {
	statuses, err := status.GetStatus(sc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(statuses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (sc *StatusController) GetStatus(w http.ResponseWriter, r *http.Request) {
	id := sc.Controller.GETId(w, r)
	s, err := status.GetStatusById(sc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Status not found: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (sc *StatusController) CreateStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var s status.Status
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := status.CreateStatus(sc.Controller.Dependencies.DBDecorator.GDB(), &s); err != nil {
		http.Error(w, "Failed to create status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (sc *StatusController) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var s status.Status
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := status.UpdateStatus(sc.Controller.Dependencies.DBDecorator.GDB(), &s)
	if err != nil {
		http.Error(w, "Failed to update status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (sc *StatusController) DeleteStatus(w http.ResponseWriter, r *http.Request) {
	id := sc.Controller.GETId(w, r)
	err := status.DeleteStatusById(sc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Failed to delete status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
