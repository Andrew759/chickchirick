package user

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	"chickChirick/internal/middleware/config"
	userMiddleware "chickChirick/internal/middleware/validators/user"
	property "chickChirick/internal/model/user"
	"net/http"
)

type PropertyController struct {
	Controller abstraction.Controller
	userMiddleware.PropertyValidator
}

func (pc *PropertyController) HandleRequest() {
	pc.Controller.ServeMux.HandleFunc("/user/properties", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			pc.GetProperties(w)
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	pc.Controller.ServeMux.HandleFunc("/user/property", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			pc.Validate(func(w http.ResponseWriter, r *http.Request) {
				pc.CreateProperty(w, c_http.NewRequest(r))
			})(w, r)
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	pc.Controller.ServeMux.HandleFunc("/user/property/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			pc.GetProperty(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/property/")))
		case http.MethodPut:
			pc.Validate(func(w http.ResponseWriter, r *http.Request) {
				pc.UpdateProperty(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/property/")))
			})(w, r)
		case http.MethodDelete:
			pc.DeleteProperty(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/property/")))
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (pc *PropertyController) GetProperties(w http.ResponseWriter) {
	properties, err := property.GetProperties(pc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, properties, http.StatusOK)
}

func (pc *PropertyController) GetProperty(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := property.GetPropertyById(pc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Property not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, p, http.StatusOK)
}

func (pc *PropertyController) CreateProperty(w http.ResponseWriter, r *c_http.Request) {
	p := r.Context().Value(config.UserPropertyKey).(*property.Property)

	if err := property.CreateProperty(pc.Controller.Dependencies.DBDecorator.GDB(), p); err != nil {
		c_http.NewResponse().SendError(w, "Failed to create property: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, p, http.StatusCreated)
}

func (pc *PropertyController) UpdateProperty(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	p := r.Context().Value(config.UserPropertyKey).(*property.Property)

	err = property.UpdatePropertyById(pc.Controller.Dependencies.DBDecorator.GDB(), p, id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update property: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, p, http.StatusOK)
}

func (pc *PropertyController) DeleteProperty(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = property.DeletePropertyById(pc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete property: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
