package message

import (
	"chickChirick/internal/controller/c_controller"
	"chickChirick/internal/controller/c_http"
	file "chickChirick/internal/model/message"
	"encoding/json"
	"errors"
	"net/http"
)

type FileController struct {
	Controller c_controller.Controller
}

func (fc *FileController) HandleRequest() {
	fc.Controller.ServeMux.HandleFunc("GET /files", func(w http.ResponseWriter, r *http.Request) {
		fc.GetFiles(w)
	})

	fc.Controller.ServeMux.HandleFunc("POST /file", func(w http.ResponseWriter, r *http.Request) {
		fc.CreateFile(w, c_http.NewRequest(r))
	})

	fc.Controller.ServeMux.HandleFunc("GET /file/{id}", func(w http.ResponseWriter, r *http.Request) {
		fc.GetFile(w, c_http.NewRequest(r))
	})

	fc.Controller.ServeMux.HandleFunc("PUT /file/{id}", func(w http.ResponseWriter, r *http.Request) {
		fc.UpdateFile(w, c_http.NewRequest(r))
	})

	fc.Controller.ServeMux.HandleFunc("DELETE /file/{id}", func(w http.ResponseWriter, r *http.Request) {
		fc.DeleteFile(w, c_http.NewRequest(r))
	})
}

func (fc *FileController) GetFiles(w http.ResponseWriter) {
	files, err := file.GetFiles(fc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, files, http.StatusOK)
}

func (fc *FileController) GetFile(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	f, err := file.GetFileById(fc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "File not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, f, http.StatusOK)
}

func (fc *FileController) CreateFile(w http.ResponseWriter, r *c_http.Request) {
	var f file.File
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := file.CreateFile(fc.Controller.Dependencies.DBDecorator.GDB(), &f); err != nil {
		c_http.NewResponse().SendError(w, "Failed to create file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, f, http.StatusCreated)
}

func (fc *FileController) UpdateFile(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var f file.File
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = file.UpdateFileById(fc.Controller.Dependencies.DBDecorator.GDB(), &f, id)
	if err != nil && errors.Is(err, file.FileNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, f, http.StatusOK)
}

func (fc *FileController) DeleteFile(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = file.DeleteFileById(fc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil && errors.Is(err, file.FileNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
