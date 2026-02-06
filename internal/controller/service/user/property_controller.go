package user

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	"chickChirick/internal/middleware"
	"chickChirick/internal/middleware/config"
	property "chickChirick/internal/model/user"
	"errors"
	"net/http"
)

type PropertyController struct {
	Controller abstraction.Controller
	middleware.Validator
}

func (pc *PropertyController) HandleRequest() {
	pc.Controller.ServeMux.HandleFunc("GET /properties", func(w http.ResponseWriter, r *http.Request) {
		pc.GetProperties(w)
	})

	pc.Controller.ServeMux.HandleFunc("POST /property",
		pc.Validate(func(w http.ResponseWriter, r *http.Request) {
			pc.CreateProperty(w, c_http.NewRequest(r))
		}))

	pc.Controller.ServeMux.HandleFunc("GET /property/{id}", func(w http.ResponseWriter, r *http.Request) {
		pc.GetProperty(w, c_http.NewRequest(r))

	})

	pc.Controller.ServeMux.HandleFunc("PUT /property/{id}",
		pc.Validate(func(w http.ResponseWriter, r *http.Request) {
			pc.UpdateProperty(w, c_http.NewRequest(r))
		}))

	pc.Controller.ServeMux.HandleFunc("DELETE /property/{id}", func(w http.ResponseWriter, r *http.Request) {
		pc.DeleteProperty(w, c_http.NewRequest(r))
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
	id, err := r.HTTPId()
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

	err := property.CreateProperty(pc.Controller.Dependencies.DBDecorator.GDB(), p)
	if err != nil && errors.Is(err, property.PropertyForUserAlreadyExistsErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusConflict)
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to create property: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, p, http.StatusCreated)
}

func (pc *PropertyController) UpdateProperty(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	p := r.Context().Value(config.UserPropertyKey).(*property.Property)

	err = property.UpdatePropertyById(pc.Controller.Dependencies.DBDecorator.GDB(), p, id)
	if err != nil && errors.Is(err, property.PropertyNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update property: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, p, http.StatusOK)
}

func (pc *PropertyController) DeleteProperty(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = property.DeletePropertyById(pc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil && errors.Is(err, property.PropertyNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete property: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
