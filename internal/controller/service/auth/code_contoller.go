package auth

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	code "chickChirick/internal/model/auth"
	"encoding/json"
	"net/http"
)

type CodeController struct {
	Controller abstraction.Controller
}

func (cc *CodeController) HandleRequest() {
	cc.Controller.ServeMux.HandleFunc("/auth/codes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			cc.GetCodes(w)
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
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
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (cc *CodeController) GetCodes(w http.ResponseWriter) {
	codes, err := code.GetCodes(cc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, codes, http.StatusOK)
}

func (cc *CodeController) GetCode(w http.ResponseWriter, r *http.Request) {
	id, err := cc.Controller.HttpId(w, r, CodeResource)
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
	}

	c, err := code.GetCodeById(cc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Code not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, c, http.StatusOK)
}

func (cc *CodeController) CreateCode(w http.ResponseWriter, r *http.Request) {
	var c code.Code
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := code.CreateCode(cc.Controller.Dependencies.DBDecorator.GDB(), &c); err != nil {
		c_http.NewResponse().SendError(w, "Failed to create code: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, c, http.StatusCreated)
}

func (cc *CodeController) UpdateCode(w http.ResponseWriter, r *http.Request) {
	var c code.Code
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := code.UpdateCode(cc.Controller.Dependencies.DBDecorator.GDB(), &c)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update code: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, c, http.StatusOK)
}

func (cc *CodeController) DeleteCode(w http.ResponseWriter, r *http.Request) {
	id, err := cc.Controller.HttpId(w, r, CodeResource)
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
	}

	err = code.DeleteCodeById(cc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete code: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
