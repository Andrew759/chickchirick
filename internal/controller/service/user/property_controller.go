package user

import (
	"chickChirick/internal/controller/c_controller"
	"chickChirick/internal/controller/c_http"
	"chickChirick/internal/middleware"
	"chickChirick/internal/middleware/config"
	authValidator "chickChirick/internal/middleware/validator"
	property "chickChirick/internal/model/user"
	"errors"
	"net/http"
)

type PropertyController struct {
	Controller c_controller.Controller
	CPV        middleware.Validator
	UPV        middleware.Validator
	authValidator.AuthValidator
}

func (pc *PropertyController) HandleRequest() {
	pc.Controller.ServeMux.HandleFunc("GET /properties", func(w http.ResponseWriter, r *http.Request) {
		pc.GetProperties(w, c_http.NewRequest(r))
	})

	pc.Controller.ServeMux.HandleFunc("POST /property",
		pc.ValidateAuth(pc.CPV.Validate(func(w http.ResponseWriter, r *http.Request) {
			pc.CreateProperty(w, c_http.NewRequest(r))
		})))

	pc.Controller.ServeMux.HandleFunc("GET /user/{id}/property", func(w http.ResponseWriter, r *http.Request) {
		pc.ValidateAuth(func(w http.ResponseWriter, r *http.Request) {
			pc.GetPropertyByUserId(w, c_http.NewRequest(r))
		})
	})

	pc.Controller.ServeMux.HandleFunc("PUT /user/{id}/property",
		pc.ValidateAuth(pc.UPV.Validate(func(w http.ResponseWriter, r *http.Request) {
			pc.UpdateProperty(w, c_http.NewRequest(r))
		})))

	pc.Controller.ServeMux.HandleFunc("DELETE /user/{id}/property", func(w http.ResponseWriter, r *http.Request) {
		pc.ValidateAuth(func(w http.ResponseWriter, r *http.Request) {
			pc.DeleteProperty(w, c_http.NewRequest(r))
		})
	})
}

func (pc *PropertyController) GetProperties(w http.ResponseWriter, r *c_http.Request) {
	ctx := r.Context()
	properties, err := property.GetProperties(ctx, pc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, properties, http.StatusOK)
}

func (pc *PropertyController) GetPropertyByUserId(w http.ResponseWriter, r *c_http.Request) {
	userId, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	p, err := property.GetPropertyByUserId(ctx, pc.Controller.Dependencies.DBDecorator.GDB(), userId)
	if err != nil {
		c_http.NewResponse().SendError(w, "Property not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, p, http.StatusOK)
}

func (pc *PropertyController) CreateProperty(w http.ResponseWriter, r *c_http.Request) {
	ctx := r.Context()
	p := ctx.Value(config.UserPropertyKey).(*property.Property)

	err := property.CreateProperty(ctx, pc.Controller.Dependencies.DBDecorator.GDB(), p)
	if err != nil && errors.Is(err, property.PropertyForUserAlreadyExistsErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusConflict)
		return
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

	ctx := r.Context()
	p := ctx.Value(config.UserPropertyKey).(*property.Property)

	err = property.UpdatePropertyByUserId(ctx, pc.Controller.Dependencies.DBDecorator.GDB(), p, id)
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

	ctx := r.Context()
	err = property.DeletePropertyByUserId(ctx, pc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil && errors.Is(err, property.PropertyNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete property: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
