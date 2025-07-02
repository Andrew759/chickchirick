package user

import (
	"chickChirick/internal/controller/abstraction"
	property "chickChirick/internal/model/user"
	"encoding/json"
	"net/http"
)

type PropertyController struct {
	Controller abstraction.Controller
}

func (pc *PropertyController) HandleRequest() {
	http.HandleFunc("/properties", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			pc.GetProperties(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/property", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			pc.GetProperty(w, r)
		case http.MethodPost:
			pc.CreateProperty(w, r)
		case http.MethodDelete:
			pc.DeleteProperty(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (pc *PropertyController) GetProperties(w http.ResponseWriter) {
	properties, err := property.GetProperties(pc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(properties)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (pc *PropertyController) GetProperty(w http.ResponseWriter, r *http.Request) {
	id := pc.Controller.GETId(w, r)
	p, err := property.GetPropertyById(pc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Property not found: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (pc *PropertyController) CreateProperty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var p property.Property
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := property.CreateProperty(pc.Controller.Dependencies.DBDecorator.GDB(), &p); err != nil {
		http.Error(w, "Failed to create property: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (pc *PropertyController) DeleteProperty(w http.ResponseWriter, r *http.Request) {
	id := pc.Controller.GETId(w, r)
	err := property.DeletePropertyById(pc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Failed to delete property: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
