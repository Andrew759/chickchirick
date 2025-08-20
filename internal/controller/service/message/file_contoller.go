package message

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	file "chickChirick/internal/model/message"
	"encoding/json"
	"net/http"
)

type FileController struct {
	Controller abstraction.Controller
}

func (fc *FileController) HandleRequest() {
	fc.Controller.ServeMux.HandleFunc("/message/files", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fc.GetFiles(w)
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fc.Controller.ServeMux.HandleFunc("/message/file", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			fc.CreateFile(w, c_http.NewRequest(r))
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fc.Controller.ServeMux.HandleFunc("/message/file/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fc.GetFile(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/message/file/")))
		case http.MethodPut:
			fc.UpdateFile(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/message/file/")))
		case http.MethodDelete:
			fc.DeleteFile(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/message/file/")))
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
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
	id, err := r.HttpId()
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
	var f file.File
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := file.UpdateFile(fc.Controller.Dependencies.DBDecorator.GDB(), &f)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, f, http.StatusOK)
}

func (fc *FileController) DeleteFile(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = file.DeleteFileById(fc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
