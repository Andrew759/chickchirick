package user

import (
	"chickChirick/internal/controller/abstraction"
	ban "chickChirick/internal/model/user"
	"encoding/json"
	"net/http"
	"strconv"
)

type BanController struct {
	MainController abstraction.Controller
}

func (bc *BanController) HandleRequest() {
	http.HandleFunc("/bans", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bc.GetBans(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/ban", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bc.GetBan(w, r)
		case http.MethodPost:
			bc.Ban(w, r)
		case http.MethodDelete:
			bc.DeleteBan(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (bc *BanController) GetBans(w http.ResponseWriter) {
	bans, err := ban.GetBans(bc.MainController.Dependencies.DBDecorator.GDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(bans)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (bc *BanController) GetBan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idVal := r.URL.Query().Get("id")
	if idVal == "" {
		http.Error(w, "Missing ban ID", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(idVal)

	u, err := ban.GetBanById(bc.MainController.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Ban not found: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (bc *BanController) Ban(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var b ban.Ban
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := ban.CreateBan(bc.MainController.Dependencies.DBDecorator.GDB(), &b); err != nil {
		http.Error(w, "Failed to create ban: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (bc *BanController) DeleteBan(w http.ResponseWriter, r *http.Request) {
	id := bc.MainController.GETId(w, r)

	err := ban.DeleteBanById(bc.MainController.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Failed to delete ban: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
