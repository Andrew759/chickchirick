package user

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	photo "chickChirick/internal/model/user"
	"encoding/json"
	"net/http"
)

type PhotoController struct {
	Controller abstraction.Controller
}

func (pc *PhotoController) HandleRequest() {
	pc.Controller.ServeMux.HandleFunc("/user/photos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			pc.GetPhotos(w)
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	pc.Controller.ServeMux.HandleFunc("/user/photo", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			pc.CreatePhoto(w, c_http.NewRequest(r))
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	pc.Controller.ServeMux.HandleFunc("/user/photo/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			pc.GetPhoto(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/user/photo/")))
		case http.MethodPut:
			pc.UpdatePhoto(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/user/photo/")))
		case http.MethodDelete:
			pc.DeletePhoto(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/user/photo/")))
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (pc *PhotoController) GetPhotos(w http.ResponseWriter) {
	photos, err := photo.GetPhotos(pc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, photos, http.StatusOK)
}

func (pc *PhotoController) GetPhoto(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := photo.GetPhotoById(pc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Photo not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, p, http.StatusOK)
}

func (pc *PhotoController) CreatePhoto(w http.ResponseWriter, r *c_http.Request) {
	var p photo.Photo
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := photo.CreatePhoto(pc.Controller.Dependencies.DBDecorator.GDB(), &p); err != nil {
		c_http.NewResponse().SendError(w, "Failed to create photo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, p, http.StatusCreated)
}

func (pc *PhotoController) UpdatePhoto(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var p photo.Photo
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = photo.UpdatePhotoById(pc.Controller.Dependencies.DBDecorator.GDB(), &p, id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update photo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, p, http.StatusOK)
}

func (pc *PhotoController) DeletePhoto(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = photo.DeletePhotoById(pc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete photo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
