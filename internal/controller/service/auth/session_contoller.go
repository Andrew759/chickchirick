package auth

import (
	"chickChirick/internal/controller/abstraction"
	session "chickChirick/internal/model/auth"
	"encoding/json"
	"net/http"
)

const SessionResource = "/auth/session/"

type SessionController struct {
	Controller abstraction.Controller
}

func (sc *SessionController) HandleRequest() {
	sc.Controller.ServeMux.HandleFunc("/auth/sessions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			sc.GetSessions(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	sc.Controller.ServeMux.HandleFunc(SessionResource, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			sc.GetSession(w, r)
		case http.MethodPost:
			sc.CreateSession(w, r)
		case http.MethodPut:
			sc.UpdateSession(w, r)
		case http.MethodDelete:
			sc.DeleteSession(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (sc *SessionController) GetSessions(w http.ResponseWriter) {
	codes, err := session.GetSessions(sc.Controller.Dependencies.DBDecorator.GDB())
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

func (sc *SessionController) GetSession(w http.ResponseWriter, r *http.Request) {
	id := sc.Controller.HttpId(w, r, SessionResource)
	c, err := session.GetSessionById(sc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Session not found: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (sc *SessionController) CreateSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var s session.Session
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := session.CreateSession(sc.Controller.Dependencies.DBDecorator.GDB(), &s); err != nil {
		http.Error(w, "Failed to create session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (sc *SessionController) UpdateSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var s session.Session
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := session.UpdateSession(sc.Controller.Dependencies.DBDecorator.GDB(), &s)
	if err != nil {
		http.Error(w, "Failed to update session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (sc *SessionController) DeleteSession(w http.ResponseWriter, r *http.Request) {
	id := sc.Controller.HttpId(w, r, SessionResource)
	err := session.DeleteSessionById(sc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Failed to delete session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
