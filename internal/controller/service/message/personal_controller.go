package message

import (
	"chickChirick/internal/controller/abstraction"
	personal "chickChirick/internal/model/message"
	"encoding/json"
	"net/http"
)

type PersonalController struct {
	Controller abstraction.Controller
}

func (pc *PersonalController) HandleRequest() {
	http.HandleFunc("/message/personals", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			pc.GetPersonals(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/message/personal", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			pc.GetPersonal(w, r)
		case http.MethodPost:
			pc.CreatePersonal(w, r)
		case http.MethodPut:
			pc.UpdatePersonal(w, r)
		case http.MethodDelete:
			pc.DeletePersonal(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (pc *PersonalController) GetPersonals(w http.ResponseWriter) {
	personals, err := personal.GetPersonal(pc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(personals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (pc *PersonalController) GetPersonal(w http.ResponseWriter, r *http.Request) {
	id := pc.Controller.GETId(w, r)
	p, err := personal.GetPersonalById(pc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Personal not found: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (pc *PersonalController) CreatePersonal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var p personal.Personal
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := personal.CreatePersonal(pc.Controller.Dependencies.DBDecorator.GDB(), &p); err != nil {
		http.Error(w, "Failed to create personal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (pc *PersonalController) UpdatePersonal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var p personal.Personal
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := personal.UpdatePersonal(pc.Controller.Dependencies.DBDecorator.GDB(), &p)
	if err != nil {
		http.Error(w, "Failed to update personal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (pc *PersonalController) DeletePersonal(w http.ResponseWriter, r *http.Request) {
	id := pc.Controller.GETId(w, r)
	err := personal.DeletePersonalById(pc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Failed to delete personal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
