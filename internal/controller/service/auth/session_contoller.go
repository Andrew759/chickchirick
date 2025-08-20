package auth

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	session "chickChirick/internal/model/auth"
	"encoding/json"
	"errors"
	"net/http"
)

type SessionController struct {
	Controller abstraction.Controller
}

func (sc *SessionController) HandleRequest() {
	sc.Controller.ServeMux.HandleFunc("/auth/sessions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			sc.GetSessions(w)
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	sc.Controller.ServeMux.HandleFunc("/auth/session", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			sc.CreateSession(w, c_http.NewRequest(r))
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	sc.Controller.ServeMux.HandleFunc("/auth/session/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			sc.GetSession(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/auth/session/")))
		case http.MethodPut:
			sc.UpdateSession(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/auth/session/")))
		case http.MethodDelete:
			sc.DeleteSession(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/auth/session/")))
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (sc *SessionController) GetSessions(w http.ResponseWriter) {
	codes, err := session.GetSessions(sc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, codes, http.StatusOK)
}

func (sc *SessionController) GetSession(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err := session.GetSessionById(sc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Session not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, c, http.StatusOK)
}

func (sc *SessionController) CreateSession(w http.ResponseWriter, r *c_http.Request) {
	var s session.Session
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := session.CreateSession(sc.Controller.Dependencies.DBDecorator.GDB(), &s); err != nil {
		c_http.NewResponse().SendError(w, "Failed to create session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, s, http.StatusCreated)
}

func (sc *SessionController) UpdateSession(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var s session.Session
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = session.UpdateSessionById(sc.Controller.Dependencies.DBDecorator.GDB(), &s, id)
	if err != nil && errors.Is(err, session.SessionNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, s, http.StatusOK)
}

func (sc *SessionController) DeleteSession(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = session.DeleteSessionById(sc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil && errors.Is(err, session.SessionNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
