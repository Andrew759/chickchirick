package message

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	personal "chickChirick/internal/model/message"
	"encoding/json"
	"errors"
	"net/http"
)

type PersonalController struct {
	Controller abstraction.Controller
}

func (pc *PersonalController) HandleRequest() {
	pc.Controller.ServeMux.HandleFunc("GET /personals", func(w http.ResponseWriter, r *http.Request) {
		pc.GetPersonals(w)
	})

	pc.Controller.ServeMux.HandleFunc("POST /personal", func(w http.ResponseWriter, r *http.Request) {
		pc.CreatePersonal(w, c_http.NewRequest(r))
	})

	pc.Controller.ServeMux.HandleFunc("GET /personal/{id}", func(w http.ResponseWriter, r *http.Request) {
		pc.GetPersonal(w, c_http.NewRequest(r))
	})

	pc.Controller.ServeMux.HandleFunc("PUT /personal/{id}", func(w http.ResponseWriter, r *http.Request) {
		pc.UpdatePersonal(w, c_http.NewRequest(r))
	})

	pc.Controller.ServeMux.HandleFunc("DELETE /personal/{id}", func(w http.ResponseWriter, r *http.Request) {
		pc.DeletePersonal(w, c_http.NewRequest(r))
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

func (pc *PersonalController) GetPersonal(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := personal.GetPersonalById(pc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Personal not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, p, http.StatusOK)
}

func (pc *PersonalController) CreatePersonal(w http.ResponseWriter, r *c_http.Request) {
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

func (pc *PersonalController) UpdatePersonal(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var p personal.Personal
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = personal.UpdatePersonalById(pc.Controller.Dependencies.DBDecorator.GDB(), &p, id)
	if err != nil && errors.Is(err, personal.PersonalNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update personal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, p, http.StatusOK)
}

func (pc *PersonalController) DeletePersonal(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = personal.DeletePersonalById(pc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil && errors.Is(err, personal.PersonalNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete personal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
