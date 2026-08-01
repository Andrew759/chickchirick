package user

import (
	"chickChirick/internal/controller/c_controller"
	"chickChirick/internal/controller/c_http"
	"chickChirick/internal/middleware"
	"chickChirick/internal/middleware/config"
	authValidator "chickChirick/internal/middleware/validator"
	photo "chickChirick/internal/model/user"
	"errors"
	"net/http"
)

type PhotoController struct {
	Controller c_controller.Controller
	middleware.Validator
	authValidator.AuthValidator
}

func (pc *PhotoController) HandleRequest() {
	pc.Controller.ServeMux.HandleFunc("GET /photos", func(w http.ResponseWriter, r *http.Request) {
		pc.ValidateAuth(func(w http.ResponseWriter, r *http.Request) {
			pc.GetPhotos(w, c_http.NewRequest(r))
		})
	})

	pc.Controller.ServeMux.HandleFunc("POST /photo",
		pc.ValidateAuth(pc.Validate(func(w http.ResponseWriter, r *http.Request) {
			pc.CreatePhoto(w, c_http.NewRequest(r))
		})))

	pc.Controller.ServeMux.HandleFunc("GET /photo/{id}", func(w http.ResponseWriter, r *http.Request) {
		pc.ValidateAuth(func(w http.ResponseWriter, r *http.Request) {
			pc.GetPhoto(w, c_http.NewRequest(r))
		})
	})

	pc.Controller.ServeMux.HandleFunc("GET /user/{id}/photos", func(w http.ResponseWriter, r *http.Request) {
		pc.ValidateAuth(func(w http.ResponseWriter, r *http.Request) {
			pc.GetPhotosByUserId(w, c_http.NewRequest(r))
		})
	})

	pc.Controller.ServeMux.HandleFunc("PUT /photo/{id}",
		pc.ValidateAuth(pc.Validate(func(w http.ResponseWriter, r *http.Request) {
			pc.UpdatePhoto(w, c_http.NewRequest(r))
		})))

	pc.Controller.ServeMux.HandleFunc("DELETE /photo/{id}", func(w http.ResponseWriter, r *http.Request) {
		pc.ValidateAuth(func(w http.ResponseWriter, r *http.Request) {
			pc.DeletePhoto(w, c_http.NewRequest(r))
		})
	})
}

func (pc *PhotoController) GetPhotos(w http.ResponseWriter, r *c_http.Request) {
	ctx := r.Context()
	photos, err := photo.GetPhotos(ctx, pc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, photos, http.StatusOK)
}

func (pc *PhotoController) GetPhoto(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	p, err := photo.GetPhotoById(ctx, pc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Photo not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, p, http.StatusOK)
}

func (pc *PhotoController) GetPhotosByUserId(w http.ResponseWriter, r *c_http.Request) {
	userId, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	photos, err := photo.GetPhotosByUserId(ctx, pc.Controller.Dependencies.DBDecorator.GDB(), userId)
	if err != nil {
		c_http.NewResponse().SendError(w, "Photos by user id not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, photos, http.StatusOK)
}

func (pc *PhotoController) CreatePhoto(w http.ResponseWriter, r *c_http.Request) {
	ctx := r.Context()
	p := ctx.Value(config.UserPhotoKey).(*photo.Photo)

	//TODO: доработать ошибки
	if err := photo.CreatePhoto(ctx, pc.Controller.Dependencies.DBDecorator.GDB(), p); err != nil {
		c_http.NewResponse().SendError(w, "Failed to create photo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, p, http.StatusCreated)
}

func (pc *PhotoController) UpdatePhoto(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	p := ctx.Value(config.UserPhotoKey).(*photo.Photo)

	err = photo.UpdatePhotoById(ctx, pc.Controller.Dependencies.DBDecorator.GDB(), p, id)
	if err != nil && errors.Is(err, photo.PhotoNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update photo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, p, http.StatusOK)
}

func (pc *PhotoController) DeletePhoto(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = photo.DeletePhotoById(ctx, pc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil && errors.Is(err, photo.PhotoNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete photo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
