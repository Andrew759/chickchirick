package auth

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	token "chickChirick/internal/model/auth"
	"encoding/json"
	"net/http"
)

type TokenController struct {
	Controller abstraction.Controller
}

func (tc *TokenController) HandleRequest() {
	tc.Controller.ServeMux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			tc.GetTokens(w)
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	tc.Controller.ServeMux.HandleFunc(TokenResource, func(w http.ResponseWriter, r *http.Request) {
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
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (tc *TokenController) GetTokens(w http.ResponseWriter) {
	tokens, err := token.GetTokens(tc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, tokens, http.StatusOK)
}

func (tc *TokenController) GetToken(w http.ResponseWriter, r *http.Request) {
	id, err := tc.Controller.HttpId(w, r, TokenResource)
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
	}

	t, err := token.GetTokenById(tc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Token not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, t, http.StatusOK)
}

func (tc *TokenController) CreateToken(w http.ResponseWriter, r *http.Request) {
	var t token.Token
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := token.CreateToken(tc.Controller.Dependencies.DBDecorator.GDB(), &t); err != nil {
		c_http.NewResponse().SendError(w, "Failed to create token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, t, http.StatusCreated)
}

func (tc *TokenController) UpdateToken(w http.ResponseWriter, r *http.Request) {
	var t token.Token
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := token.UpdateToken(tc.Controller.Dependencies.DBDecorator.GDB(), &t)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, t, http.StatusOK)
}

func (tc *TokenController) DeleteToken(w http.ResponseWriter, r *http.Request) {
	id, err := tc.Controller.HttpId(w, r, TokenResource)
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
	}

	err = token.DeleteTokenById(tc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
