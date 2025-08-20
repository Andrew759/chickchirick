package auth

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	code "chickChirick/internal/model/auth"
	"encoding/json"
	"errors"
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

	cc.Controller.ServeMux.HandleFunc("/auth/code", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			cc.CreateCode(w, c_http.NewRequest(r))
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	cc.Controller.ServeMux.HandleFunc("/auth/code/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			cc.GetCode(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/auth/code/")))
		case http.MethodPut:
			cc.UpdateCode(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/auth/code/")))
		case http.MethodDelete:
			cc.DeleteCode(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/auth/code/")))
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

func (cc *CodeController) GetCode(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err := code.GetCodeById(cc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Code not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, c, http.StatusOK)
}

func (cc *CodeController) CreateCode(w http.ResponseWriter, r *c_http.Request) {
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

func (cc *CodeController) UpdateCode(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var c code.Code
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = code.UpdateCodeById(cc.Controller.Dependencies.DBDecorator.GDB(), &c, id)
	if err != nil && errors.Is(err, code.CodeNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update code: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, c, http.StatusOK)
}

func (cc *CodeController) DeleteCode(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = code.DeleteCodeById(cc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil && errors.Is(err, code.CodeNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete code: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
