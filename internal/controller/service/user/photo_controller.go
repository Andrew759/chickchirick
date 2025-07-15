package user

import (
	"chickChirick/internal/controller/abstraction"
	photo "chickChirick/internal/model/user"
	"encoding/json"
	"net/http"
)

type PhotoController struct {
	Controller abstraction.Controller
}

func (pc *PhotoController) HandleRequest() {
	http.HandleFunc("/user/photos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			pc.GetPhotos(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/user/photo", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			pc.GetPhoto(w, r)
		case http.MethodPost:
			pc.CreatePhoto(w, r)
		case http.MethodPut:
			pc.UpdatePhoto(w, r)
		case http.MethodDelete:
			pc.DeletePhoto(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (pc *PhotoController) GetPhotos(w http.ResponseWriter) {
	photos, err := photo.GetPhotos(pc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(photos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (pc *PhotoController) GetPhoto(w http.ResponseWriter, r *http.Request) {
	id := pc.Controller.GETId(w, r)
	p, err := photo.GetPhotoById(pc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Photo not found: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (pc *PhotoController) CreatePhoto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var p photo.Photo
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := photo.CreatePhoto(pc.Controller.Dependencies.DBDecorator.GDB(), &p); err != nil {
		http.Error(w, "Failed to create photo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (pc *PhotoController) UpdatePhoto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var p photo.Photo
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := photo.UpdatePhoto(pc.Controller.Dependencies.DBDecorator.GDB(), &p)
	if err != nil {
		http.Error(w, "Failed to update photo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (pc *PhotoController) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	id := pc.Controller.GETId(w, r)
	err := photo.DeletePhotoById(pc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Failed to delete photo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
