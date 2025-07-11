package message

import (
	"chickChirick/internal/controller/abstraction"
	file "chickChirick/internal/model/message"
	"encoding/json"
	"net/http"
)

type FileController struct {
	Controller abstraction.Controller
}

func (fc *FileController) HandleRequest() {
	http.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fc.GetFiles(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fc.GetFile(w, r)
		case http.MethodPost:
			fc.CreateFile(w, r)
		case http.MethodPut:
			fc.UpdateFile(w, r)
		case http.MethodDelete:
			fc.DeleteFile(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (fc *FileController) GetFiles(w http.ResponseWriter) {
	files, err := file.GetFiles(fc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (fc *FileController) GetFile(w http.ResponseWriter, r *http.Request) {
	id := fc.Controller.GETId(w, r)
	f, err := file.GetFileById(fc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "File not found: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (fc *FileController) CreateFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var f file.File
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := file.CreateFile(fc.Controller.Dependencies.DBDecorator.GDB(), &f); err != nil {
		http.Error(w, "Failed to create file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (fc *FileController) UpdateFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var f file.File
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := file.UpdateFile(fc.Controller.Dependencies.DBDecorator.GDB(), &f)
	if err != nil {
		http.Error(w, "Failed to update file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (fc *FileController) DeleteFile(w http.ResponseWriter, r *http.Request) {
	id := fc.Controller.GETId(w, r)
	err := file.DeleteFileById(fc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Failed to delete file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
