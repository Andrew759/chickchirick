package auth

import (
	"chickChirick/internal/controller/abstraction"
	code "chickChirick/internal/model/auth"
	"encoding/json"
	"net/http"
)

const CodeResource = "/auth/code/"

type CodeController struct {
	Controller abstraction.Controller
}

func (cc *CodeController) HandleRequest() {
	cc.Controller.ServeMux.HandleFunc("/auth/codes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			cc.GetCodes(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	cc.Controller.ServeMux.HandleFunc(CodeResource, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			cc.GetCode(w, r)
		case http.MethodPost:
			cc.CreateCode(w, r)
		case http.MethodPut:
			cc.UpdateCode(w, r)
		case http.MethodDelete:
			cc.DeleteCode(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (cc *CodeController) GetCodes(w http.ResponseWriter) {
	codes, err := code.GetCodes(cc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(codes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (cc *CodeController) GetCode(w http.ResponseWriter, r *http.Request) {
	id := cc.Controller.HttpId(w, r, CodeResource)
	c, err := code.GetCodeById(cc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Code not found: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (cc *CodeController) CreateCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var c code.Code
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := code.CreateCode(cc.Controller.Dependencies.DBDecorator.GDB(), &c); err != nil {
		http.Error(w, "Failed to create code: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (cc *CodeController) UpdateCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var c code.Code
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := code.UpdateCode(cc.Controller.Dependencies.DBDecorator.GDB(), &c)
	if err != nil {
		http.Error(w, "Failed to update code: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (cc *CodeController) DeleteCode(w http.ResponseWriter, r *http.Request) {
	id := cc.Controller.HttpId(w, r, CodeResource)
	err := code.DeleteCodeById(cc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Failed to delete code: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
