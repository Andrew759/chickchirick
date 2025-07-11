package auth

import (
	"chickChirick/internal/controller/abstraction"
	token "chickChirick/internal/model/auth"
	"encoding/json"
	"net/http"
)

type TokenController struct {
	Controller abstraction.Controller
}

func (tc *TokenController) HandleRequest() {
	http.HandleFunc("/tokens", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			tc.GetTokens(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			tc.GetToken(w, r)
		case http.MethodPost:
			tc.CreateToken(w, r)
		case http.MethodPut:
			tc.UpdateToken(w, r)
		case http.MethodDelete:
			tc.DeleteToken(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (tc *TokenController) GetTokens(w http.ResponseWriter) {
	tokens, err := token.GetTokens(tc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tokens)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (tc *TokenController) GetToken(w http.ResponseWriter, r *http.Request) {
	id := tc.Controller.GETId(w, r)
	t, err := token.GetTokenById(tc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Token not found: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (tc *TokenController) CreateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var t token.Token
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := token.CreateToken(tc.Controller.Dependencies.DBDecorator.GDB(), &t); err != nil {
		http.Error(w, "Failed to create token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (tc *TokenController) UpdateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var t token.Token
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := token.UpdateToken(tc.Controller.Dependencies.DBDecorator.GDB(), &t)
	if err != nil {
		http.Error(w, "Failed to update token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (tc *TokenController) DeleteToken(w http.ResponseWriter, r *http.Request) {
	id := tc.Controller.GETId(w, r)
	err := token.DeleteTokenById(tc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Failed to delete token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
