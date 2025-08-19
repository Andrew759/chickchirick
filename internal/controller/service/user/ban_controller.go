package user

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	ban "chickChirick/internal/model/user"
	"encoding/json"
	"net/http"
)

type BanController struct {
	Controller abstraction.Controller
}

func (bc *BanController) HandleRequest() {
	bc.Controller.ServeMux.HandleFunc("/user/bans", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bc.GetBans(w)
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	bc.Controller.ServeMux.HandleFunc("/user/ban", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			bc.CreateBan(w, c_http.NewRequest(r))
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	bc.Controller.ServeMux.HandleFunc("/user/ban/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bc.GetBan(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/user/ban/")))
		case http.MethodPut:
			bc.UpdateBan(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/user/ban/")))
		case http.MethodDelete:
			bc.DeleteBan(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/user/ban/")))
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (bc *BanController) GetBans(w http.ResponseWriter) {
	bans, err := ban.GetBans(bc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, bans, http.StatusOK)
}

func (bc *BanController) GetBan(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, err := ban.GetBanById(bc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Ban not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, b, http.StatusOK)
}

func (bc *BanController) CreateBan(w http.ResponseWriter, r *c_http.Request) {
	var b ban.Ban
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := ban.CreateBan(bc.Controller.Dependencies.DBDecorator.GDB(), &b); err != nil {
		c_http.NewResponse().SendError(w, "Failed to create ban: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, b, http.StatusCreated)
}

func (bc *BanController) UpdateBan(w http.ResponseWriter, r *c_http.Request) {
	var b ban.Ban
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := ban.UpdateBan(bc.Controller.Dependencies.DBDecorator.GDB(), &b)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update ban: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, b, http.StatusOK)
}

func (bc *BanController) DeleteBan(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ban.DeleteBanById(bc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete ban: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
