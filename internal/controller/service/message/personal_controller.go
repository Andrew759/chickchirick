package message

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	personal "chickChirick/internal/model/message"
	"encoding/json"
	"net/http"
)

type PersonalController struct {
	Controller abstraction.Controller
}

func (pc *PersonalController) HandleRequest() {
	pc.Controller.ServeMux.HandleFunc("/message/personals", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			pc.GetPersonals(w)
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	pc.Controller.ServeMux.HandleFunc(PersonalResource, func(w http.ResponseWriter, r *http.Request) {
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
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (pc *PersonalController) GetPersonals(w http.ResponseWriter) {
	personals, err := personal.GetPersonal(pc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, personals, http.StatusOK)
}

func (pc *PersonalController) GetPersonal(w http.ResponseWriter, r *http.Request) {
	id, err := pc.Controller.HttpId(w, r, PersonalResource)
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
	}

	p, err := personal.GetPersonalById(pc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Personal not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, p, http.StatusOK)
}

func (pc *PersonalController) CreatePersonal(w http.ResponseWriter, r *http.Request) {
	var p personal.Personal
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := personal.CreatePersonal(pc.Controller.Dependencies.DBDecorator.GDB(), &p); err != nil {
		c_http.NewResponse().SendError(w, "Failed to create personal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, p, http.StatusCreated)
}

func (pc *PersonalController) UpdatePersonal(w http.ResponseWriter, r *http.Request) {
	var p personal.Personal
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := personal.UpdatePersonal(pc.Controller.Dependencies.DBDecorator.GDB(), &p)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update personal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, p, http.StatusOK)
}

func (pc *PersonalController) DeletePersonal(w http.ResponseWriter, r *http.Request) {
	id, err := pc.Controller.HttpId(w, r, PersonalResource)
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
	}

	err = personal.DeletePersonalById(pc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete personal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
