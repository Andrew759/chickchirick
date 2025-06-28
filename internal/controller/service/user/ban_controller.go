package user

import (
	"chickChirick/internal/controller/abstraction"
	ban "chickChirick/internal/model/user"
	"encoding/json"
	"net/http"
	"os/user"
	"strconv"
)

type BanController struct {
	MainController abstraction.Controller
}

func (uc *BanController) HandleRequest() {
	http.HandleFunc("/bans", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			uc.GetBans(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/ban", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			uc.GetBan(w, r)
		case http.MethodPost:
			uc.Ban(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (uc *BanController) GetBans(w http.ResponseWriter) {
	bans, err := ban.GetBans(uc.MainController.Dependencies.DBDecorator.GDB())
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

func (uc *BanController) GetBan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idVal := r.URL.Query().Get("id")
	if idVal == "" {
		http.Error(w, "Missing ban ID", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(idVal)

	u, err := ban.GetBanById(uc.MainController.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Ban not found: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (uc *BanController) Ban(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var u user.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := user.CreateUser(uc.MainController.Dependencies.DBDecorator.GDB(), &u); err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
