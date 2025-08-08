package user

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/http_transaction"
	ban "chickChirick/internal/model/user"
	"encoding/json"
	"net/http"
)

const BanResource = "/user/ban/"

type BanController struct {
	Controller abstraction.Controller
}

func (bc *BanController) HandleRequest() {
	bc.Controller.ServeMux.HandleFunc("/user/bans", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bc.GetBans(w)
		default:
			http_transaction.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	bc.Controller.ServeMux.HandleFunc(BanResource, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bc.GetBan(w, r)
		case http.MethodPost:
			bc.Ban(w, r)
		case http.MethodPut:
			bc.UpdateBan(w, r)
		case http.MethodDelete:
			bc.DeleteBan(w, r)
		default:
			http_transaction.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (bc *BanController) GetBans(w http.ResponseWriter) {
	bans, err := ban.GetBans(bc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		http_transaction.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http_transaction.NewResponse().SendSuccess(w, bans, http.StatusOK)
}

func (bc *BanController) GetBan(w http.ResponseWriter, r *http.Request) {
	id := bc.Controller.HttpId(w, r, BanResource)
	b, err := ban.GetBanById(bc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http_transaction.NewResponse().SendError(w, "Ban not found: "+err.Error(), http.StatusNotFound)
		return
	}

	http_transaction.NewResponse().SendSuccess(w, b, http.StatusOK)
}

func (bc *BanController) Ban(w http.ResponseWriter, r *http.Request) {
	var b ban.Ban
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http_transaction.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := ban.CreateBan(bc.Controller.Dependencies.DBDecorator.GDB(), &b); err != nil {
		http_transaction.NewResponse().SendError(w, "Failed to create ban: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http_transaction.NewResponse().SendSuccess(w, b, http.StatusOK)
}

func (bc *BanController) UpdateBan(w http.ResponseWriter, r *http.Request) {
	var b ban.Ban
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http_transaction.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := ban.UpdateBan(bc.Controller.Dependencies.DBDecorator.GDB(), &b)
	if err != nil {
		http_transaction.NewResponse().SendError(w, "Failed to update ban: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http_transaction.NewResponse().SendSuccess(w, b, http.StatusOK)
}

func (bc *BanController) DeleteBan(w http.ResponseWriter, r *http.Request) {
	id := bc.Controller.HttpId(w, r, BanResource)
	err := ban.DeleteBanById(bc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http_transaction.NewResponse().SendError(w, "Failed to delete ban: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
