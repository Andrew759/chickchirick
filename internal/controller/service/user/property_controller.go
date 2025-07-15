package user

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/middleware/config"
	userMiddleware "chickChirick/internal/middleware/validators/user"
	property "chickChirick/internal/model/user"
	"encoding/json"
	"net/http"
)

type PropertyController struct {
	Controller abstraction.Controller
	userMiddleware.PropertyValidator
}

func (pc *PropertyController) HandleRequest() {
	http.HandleFunc("/user/properties", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			pc.GetProperties(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/user/property", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			pc.GetProperty(w, r)
		case http.MethodPost:
			pc.Validate(pc.CreateProperty)(w, r)
		case http.MethodPut:
			pc.Validate(pc.UpdateProperty)(w, r)
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

	p := r.Context().Value(config.UserPropertyKey).(*property.Property)

	if err := property.CreateProperty(pc.Controller.Dependencies.DBDecorator.GDB(), p); err != nil {
		http.Error(w, "Failed to create property: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (pc *PropertyController) UpdateProperty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	p := r.Context().Value(config.UserPropertyKey).(*property.Property)

	err := property.UpdateProperty(pc.Controller.Dependencies.DBDecorator.GDB(), p)
	if err != nil {
		http.Error(w, "Failed to update property: "+err.Error(), http.StatusInternalServerError)
		return
	}

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
